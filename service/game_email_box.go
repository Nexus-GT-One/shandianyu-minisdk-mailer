package service

import (
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/provider/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type gameMailBox struct{}

var GameMailBox = newGameMailBox()

func newGameMailBox() *gameMailBox {
	return &gameMailBox{}
}

func (o *gameMailBox) ListAll() []*entity.GameMailBox {
	ctx, cursor := db.Find(bson.D{}, entity.GameMailBox{})
	return mongodb.DecodeList(ctx, cursor, entity.GameMailBox{})
}

func (o *gameMailBox) GetGameGameMailBox(game *entity.Game) *entity.GameMailBox {
	query := bson.D{{"symbol", game.Symbol}}
	ctx, cursor := db.FindOne(query, entity.GameMailBox{})
	return mongodb.DecodeOne(ctx, cursor, entity.GameMailBox{})
}

func (o *gameMailBox) SaveScheduleNextSendTime(game *entity.Game, sendTime int64) *entity.GameMailBox {
	query := bson.D{{"symbol", game.Symbol}}
	update := bson.D{
		{"symbol", game.Symbol},
		{"email", game.DeveloperEmail},
		{"sendTime", sendTime},
		{"registerTime", game.CreateTime},
	}
	db.UpdateMany(entity.GameMailBox{}, query, update, options.Update().SetUpsert(true))
	return o.GetGameGameMailBox(game)
}

func (o *gameMailBox) DeleteSchedule(game *entity.Game) {
	query := bson.D{{"symbol", game.Symbol}}
	db.DeleteOne(entity.GameMailBox{}, query)
}
