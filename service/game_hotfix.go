package service

import (
	"context"
	"encoding/json"
	"fmt"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/provider/config"
	"shandianyu-minisdk-mailer/provider/mongodb"
	"shandianyu-minisdk-mailer/thirdparty/feishu"
	"shandianyu-minisdk-mailer/util/httputil"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type gameHotfixService struct{}

var defaultValue any

var GameHotfixService = newGameHotfixService()

func newGameHotfixService() *gameHotfixService {
	return &gameHotfixService{}
}

func (a *gameHotfixService) CompareHotfixConfig(gameId, version string) {
	game := GameService.GetOneGameById(gameId)
	if game == nil || len(version) <= 0 {
		return
	}
	testUrl := fmt.Sprintf("%s/api/v1/game/gameHotfix/detail?gameId=%s&version=%s", config.GetString("application.test-dashboard-url"), gameId, version)
	prodUrl := fmt.Sprintf("%s/api/v1/game/gameHotfix/detail?gameId=%s&version=%s", config.GetString("application.prod-dashboard-url"), gameId, version)
	testHotfixConfig := getHotfixConfig(testUrl)
	prodHotfixConfig := getHotfixConfig(prodUrl)
	if len(cmp.Diff(prodHotfixConfig, testHotfixConfig)) <= 0 {
		return
	}

	message := fmt.Sprintf("【监控消息】\n游戏 %s v_%s 已过审，热更版本 %s 正式环境和测试环境的热更配置有差异，请确认！", game.Symbol, game.PublishVersion, version)
	feishu.DevRobot().SendRobotMessage(message, game.Developer...)
}

func getHotfixConfig(url string) map[string]any {
	bytes := httputil.Get(url, make(map[string]string))
	gameHotfix := make(map[string]any)
	json.Unmarshal(bytes, &gameHotfix)
	if gameHotfix["data"] == nil {
		return map[string]any{}
	}
	gameHotfix = gameHotfix["data"].(map[string]any)
	hotfixConfig := gameHotfix["value"].(map[string]any)["hotfixConfig"].(map[string]any)
	webConfig := gameHotfix["value"].(map[string]any)["webConfig"].(map[string]any)
	taskConfig := gameHotfix["value"].(map[string]any)["taskConfig"].(map[string]any)
	withdrawConfig := gameHotfix["value"].(map[string]any)["withdrawConfig"].(map[string]any)
	shareContentConfig := gameHotfix["value"].(map[string]any)["shareContentConfig"].(map[string]any)
	invitePeopleConfig := gameHotfix["value"].(map[string]any)["invitePeopleConfig"].(map[string]any)
	inviteRewardConfig := gameHotfix["value"].(map[string]any)["inviteRewardConfig"].(map[string]any)
	delete(hotfixConfig, "build")
	delete(webConfig, "build")
	delete(taskConfig, "build")
	delete(withdrawConfig, "build")
	delete(shareContentConfig, "build")
	delete(invitePeopleConfig, "build")
	delete(inviteRewardConfig, "build")
	return map[string]any{
		"hotfixConfig":       hotfixConfig,
		"webConfig":          webConfig,
		"taskConfig":         taskConfig,
		"withdrawConfig":     withdrawConfig,
		"shareContentConfig": shareContentConfig,
		"invitePeopleConfig": invitePeopleConfig,
		"inviteRewardConfig": inviteRewardConfig,
	}
}

// 返回某个游戏的所有热更版本
func (a *gameHotfixService) List(gameId string) []*entity.GameHotfix {
	ctx, cursor := db.Find(bson.D{{"gameId", gameId}}, entity.GameHotfix{}, &options.FindOptions{Sort: bson.D{{"version", -1}}})
	return mongodb.DecodeList(ctx, cursor, entity.GameHotfix{})
}

// 游戏版本绑定热更
func (a *gameHotfixService) AddVersion(gameId, version, appVersion string) {
	db.DB.Collection("GameHotfix").UpdateMany(context.Background(), bson.D{{"gameId", gameId}, {"version", version}}, bson.M{"$addToSet": bson.M{"appVersion": appVersion}})
	GameService.UpdateGame(gameId, bson.D{})
}

// 开启热更版本
func (a *gameHotfixService) EnableVersion(gameId, version string) {
	db.UpdateMany(entity.GameHotfix{}, bson.D{{"gameId", gameId}, {"version", version}}, bson.D{{"enable", true}})
	GameService.RecordGameOperateHistory(gameId, "mailer", "gameConfig", fmt.Sprintf("把热更版本 %s 的热更状态修改为 开启", version))
	GameService.UpdateGame(gameId, bson.D{})
}
