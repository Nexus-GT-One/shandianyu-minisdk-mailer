package service

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-version"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"shandianyu-minisdk-mailer/constant"
	CheckBSideResult "shandianyu-minisdk-mailer/constant"
	"shandianyu-minisdk-mailer/entity"
	loggerFactory "shandianyu-minisdk-mailer/provider/logger"
	"shandianyu-minisdk-mailer/thirdparty/feishu"
	"shandianyu-minisdk-mailer/util/arrayutil"
	"shandianyu-minisdk-mailer/util/httputil"
	"shandianyu-minisdk-mailer/util/maputil"
	"shandianyu-minisdk-mailer/util/randomutil"
	"shandianyu-minisdk-mailer/util/secretutil"
	"shandianyu-minisdk-mailer/util/stringutil"
	"strings"
	"time"
)

type applicationService struct{}

var logger = loggerFactory.GetLogger()
var ApplicationService = newApplicationService()
var storeMap = map[string]string{"iOS": "AppStore", "GooglePlay": "GooglePlay", "Samsung": "galaxyStore"}
var appStoreGeneratorMap = map[string]func(game *entity.Game) string{
	"iOS":        generateAppStoreUrl,
	"GooglePlay": generateGooglePlayStoreUrl,
	"Samsung":    generateGalaxyStoreUrl,
}

func newApplicationService() *applicationService {
	return &applicationService{}
}

func (a *applicationService) CheckApplicationNewVersion(game *entity.Game, gameMail *entity.GameMail) {
	// 下架的游戏、已经的发布过的游戏不检查
	if !game.Enable || !game.Audit[GameService.GetLastSubmitVersion(game)] {
		return
	}

	// 通过邮件获取应用版本
	publishVersion := gameMail.AppVersion

	// 如果没有新的版本上架，则结束
	storePublishVersion, _ := version.NewSemver(publishVersion)
	gamePublishVersion, _ := version.NewSemver(game.PublishVersion)
	logger.Info("[%s] publishVersion version：%s", game.Symbol, publishVersion)
	if storePublishVersion == nil || gamePublishVersion.GreaterThanOrEqual(storePublishVersion) {
		return
	}

	// 延迟解除审核模式
	publishTime := time.Now().UnixMilli() - 100000
	autoReleaseDays := GameService.GetLastSubmitAutoReleaseDays(game.BundleId, GameService.GetAuditingVersion(game))
	game = GameService.UpdateGame(game.Id.Hex(), bson.D{{"publishVersion", publishVersion}, {"publishTime", publishTime}})

	// 如果有配置自动解除审核天数，优先返回
	if autoReleaseDays > 0 {
		publishTime = time.Now().UnixMilli() + int64(autoReleaseDays*86400000)

		// 发消息通知一下延后开启审核模式
		producterMessage := fmt.Sprintf(`【过审消息】
名称：%s
代号：%s
包名：%s
商店：%s
版本：%s
备注：邮件确认过审，%s`, game.Name, game.Symbol, game.BundleId, storeMap[game.Channel], game.PublishVersion,
			fmt.Sprintf("将继续保持审核模式，预计在北京时间 %s 之后关闭审核模式", time.UnixMilli(publishTime).Format(time.DateTime)))
		game = GameService.UpdateGame(game.Id.Hex(), bson.D{{"publishVersion", publishVersion}, {"publishTime", publishTime}})
		feishu.DemandRobot().SendRobotMessage(producterMessage, game.Producter...)
		feishu.AdminRobot().SendRobotMessage(producterMessage)
	}

	// 判断是否到时间发布
	checkPublishedOrNot(game)

	// 清空服务器缓存
	SystemService.SendClearCache()
}

func checkPublishedOrNot(game *entity.Game) {
	// 判断是不是第一个包
	isFirstPackage := reflect.DeepEqual("1.0.0", game.PublishVersion) || len(game.Audit) == 1

	// GP：
	// 第一个包：在 n 天后去除审核模式，不开热更
	// 第二个包：在 n 天后去除审核模式，开启热更
	// 第三个包之后，马上去除审核模式，开启热更
	// iOS；
	// 第一个包：马上去除审核模式，不开热更
	// 第二个包：马上去除审核模式，开启热更
	autoReleaseDays := GameService.GetLastSubmitAutoReleaseDays(game.BundleId, GameService.GetAuditingVersion(game))
	game.Audit[game.PublishVersion] = false
	game.Hotfix[game.PublishVersion] = !isFirstPackage
	hotfix := arrayutil.First(GameHotfixService.List(game.Id.Hex()))

	// 游戏状态处理
	GameService.UpdateGame(game.Id.Hex(), bson.D{
		{"publishTime", 0},      // 把发布时间设置成0
		{"audit", game.Audit},   // 关闭审核模式
		{"hotfix", game.Hotfix}, // 开启热更 (旧)
		{"monitorEnable", reflect.DeepEqual("iOS", game.Channel) && !isFirstPackage}, // 可投包开启数据监控
	})

	// 审核轨迹标记“已过审”
	GameService.RecordApproved(game.BundleId, game.PublishVersion)

	// 记录游戏可投包版本
	GameService.RecordGameDispatchedVersion(game, game.PublishVersion)

	// 把数据监控开关从 关闭 修改为 开启
	if !isFirstPackage && !game.MonitorEnable && reflect.DeepEqual("iOS", game.Channel) {
		GameService.RecordGameOperateHistory(game.Id.Hex(), "mailer", "gameConfig", "把数据监控开关从 关闭 修改为 开启")
	}

	if isFirstPackage {
		GameService.RecordGameOperateHistory(game.Id.Hex(), "mailer", "gameConfig", fmt.Sprintf("把版本 v_%s 的审核状态修改为 关闭", game.PublishVersion))
	} else {
		if hotfix != nil {
			GameHotfixService.AddVersion(hotfix.GameId, hotfix.Version, game.PublishVersion)
			GameService.RecordGameOperateHistory(game.Id.Hex(), "mailer", "gameConfig", fmt.Sprintf("热更版本 %s 已经关联app版本 %s", hotfix.Version, game.PublishVersion))
			GameHotfixService.EnableVersion(hotfix.GameId, hotfix.Version)
		}
	}

	// 延迟解除审核
	hotfixStatus := map[bool]string{true: "未开启热更", false: "已开启热更，请注意监听B面数据是否正常"}
	producterMessage := fmt.Sprintf("已过审满 %d 天，已解除审核模式，%s", autoReleaseDays, hotfixStatus[isFirstPackage])
	if autoReleaseDays <= 0 {
		// 不延迟解除审核
		producterMessage = fmt.Sprintf("已过审，%s", hotfixStatus[isFirstPackage])
	}

	// 清空服务器缓存
	SystemService.SendClearCache()
	time.Sleep(5 * time.Second)

	// 白包只需要发消息给产品就好
	if isFirstPackage {
		websiteUrl := strings.ReplaceAll(game.PolicyUrl[0:strings.LastIndex(game.PolicyUrl, "/")], "www.", "")
		producterMessage = fmt.Sprintf(`【过审消息】
名称：%s
代号：%s
包名：%s
商店：%s
版本：%s
地址：%s
官网：%s
备注：邮件确认过审，%s`, game.Name, game.Symbol, game.BundleId, storeMap[game.Channel], game.PublishVersion, appStoreGeneratorMap[game.Channel](game), websiteUrl, producterMessage)
		feishu.DemandRobot().SendRobotMessage(producterMessage, game.Producter...)
		feishu.AdminRobot().SendRobotMessage(producterMessage)
		return
	}

	// 检查B面
	checkBSideResultMap := make(map[int]constant.CheckBSideResult)
	for i := 0; i < 5; i++ {
		checkBSideResultMap[i+1] = checkCanEnterBSide(game)
	}

	flag := true
	result := arrayutil.Uniq(maputil.Values(checkBSideResultMap))
	if len(result) != 1 || (!reflect.DeepEqual(CheckBSideResult.ENTER, arrayutil.First(result)) || !reflect.DeepEqual(CheckBSideResult.ENTER, arrayutil.Last(result))) {
		// flag = false
	}

	logger.Info("B面检测结果：%v", checkBSideResultMap)
	checkResult := map[bool]string{true: "已经能", false: "不能"}
	developerMessage := "请注意查看下正式服，请检查热更配置否正常无误"
	testerMessage := fmt.Sprintf("前置机测试%s进入B面，请再手动检查一遍%s的最新包是否能进入B面", checkResult[flag], storeMap[game.Channel])
	remark := fmt.Sprintf(`邮件确认过审
产品端：{{.producter}}
%s

开发端：{{.developer}}
%s

测试端：{{.tester}}
%s`, producterMessage, developerMessage, testerMessage)
	message := fmt.Sprintf(`【过审消息】
名称：%s
代号：%s
包名：%s
商店：%s
版本：%s
地址：%s
备注：%s`, game.Name, game.Symbol, game.BundleId, storeMap[game.Channel], game.PublishVersion, appStoreGeneratorMap[game.Channel](game), remark)
	producterList := strings.Join(arrayutil.Map(game.Producter, func(producter string) string {
		return fmt.Sprintf(`<at user_id="%s"></at>`, feishu.GetUserOpenId(producter).OpenId)
	}), "")
	developerList := strings.Join(arrayutil.Map(game.Developer, func(developer string) string {
		return fmt.Sprintf(`<at user_id="%s"></at>`, feishu.GetUserOpenId(developer).OpenId)
	}), "")
	testerList := strings.Join(arrayutil.Map(game.Tester, func(tester string) string {
		return fmt.Sprintf(`<at user_id="%s"></at>`, feishu.GetUserOpenId(tester).OpenId)
	}), "")
	message = stringutil.TemplateParse(message, map[string]string{
		"producter": producterList,
		"developer": developerList,
		"tester":    testerList,
	})
	feishu.DemandRobot().SendRobotMessage(message)
	feishu.AdminRobot().SendRobotMessage(message)
}

func checkCanEnterBSide(game *entity.Game) constant.CheckBSideResult {
	initUrl := game.ProductionDomain + game.Api["/api/v1/sdk/sdk/init"]
	headers := make(map[string]string)
	for key, header := range game.Header {
		if reflect.DeepEqual(key, "appVersion") {
			headers[header] = game.PublishVersion
		} else if reflect.DeepEqual(key, "device") {
			headers[header] = "shandianyu-minisdk-monitor"
		} else {
			headers[header] = randomutil.String(32)
		}
	}

	headers["bundleId"] = game.BundleId
	initResponse := make(map[string]any)
	json.Unmarshal(httputil.Get(initUrl, headers), &initResponse)
	if len(initResponse) <= 0 {
		return CheckBSideResult.CLOSE
	}

	// 检查热更配置是否跟游戏配置的一样
	desConfig := strings.TrimSpace(fmt.Sprintf("%v", initResponse["data"].(map[string]any)[game.InitRequired["des"]]))
	desDecrypt, _ := secretutil.AesDecrypt([]byte(desConfig), []byte(game.AesKey), []byte(game.AesIV))
	if len(desDecrypt) <= 0 {
		return CheckBSideResult.CLOSE
	}

	return CheckBSideResult.ENTER
}

func generateAppStoreUrl(game *entity.Game) string {
	return fmt.Sprintf("https://apps.apple.com/app/id%s", game.AppId)
}

func generateGooglePlayStoreUrl(game *entity.Game) string {
	return fmt.Sprintf("https://play.google.com/store/apps/details?id=%s", game.BundleId)
}

func generateGalaxyStoreUrl(game *entity.Game) string {
	return fmt.Sprintf("https://galaxystore.samsung.com/detail/%s", game.BundleId)
}
