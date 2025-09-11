package mail_parser

import (
	"fmt"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"shandianyu-minisdk-mailer/thirdparty/feishu"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

var appStoreGeneratorMap = map[string]func(game *entity.Game) string{
	"GooglePlay": generateGooglePlayStoreUrl,
}

// 应用下架
type suspendedMailParser struct{}

func init() {
	registerImplement(&suspendedMailParser{})
}

func (o *suspendedMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "Google Play Support")
}

func (o *suspendedMailParser) checkTitle(title string) bool {
	return strings.Contains(title, "Your app is not compliant")
}

func (o *suspendedMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "Suspended")
}

func (o *suspendedMailParser) parse(from, to, bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := service.GameService.GetByBundleId(extractBundleId(bodyText))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: oneGame.PublishVersion,
		Status:     "应用下架",
		Content:    bodyText,
	}
}

func (o *suspendedMailParser) after(game *entity.Game, gameMail *entity.GameMail) {
	game = service.GameService.GetOneGameById(game.Id.Hex())
	if game == nil || !game.Enable {
		return
	}

	// 1. 修改游戏状态
	service.GameService.UpdateGame(game.Id.Hex(), bson.D{{"enable", false}})

	// 2. 生成操作记录
	service.GameService.RecordGameOperateHistory(game.Id.Hex(), "mailer", "gameConfig", "邮件确认下架")

	// 3. 发消息通知产品、技术
	message := fmt.Sprintf(`【下架消息】
名称：%s
代号：%s
包名：%s
商店：%s
版本：%s
地址：%s
备注：邮件确认下架`, game.Name, game.Symbol, game.BundleId, storeMap[game.Channel], game.PublishVersion, generateGooglePlayStoreUrl(game))
	feishu.DemandRobot().SendRobotMessage(message, game.Producter...)
	feishu.AdminRobot().SendRobotMessage(message)
	service.SystemService.SendClearCache()
}

func generateGooglePlayStoreUrl(game *entity.Game) string {
	return fmt.Sprintf("https://play.google.com/store/apps/details?id=%s", game.BundleId)
}
