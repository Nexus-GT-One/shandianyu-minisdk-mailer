package service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/provider/mongodb"
)

type gameHotfixService struct{}

var defaultValue any

var GameHotfixService = newGameHotfixService()

func newGameHotfixService() *gameHotfixService {
	return &gameHotfixService{}
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
	GameService.RecordGameOperateHistory(gameId, "crawler", "gameConfig", fmt.Sprintf("把热更版本 %s 的热更状态修改为 开启", version))
	GameService.UpdateGame(gameId, bson.D{})
}
