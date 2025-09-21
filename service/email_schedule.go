package service

import (
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/provider/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type emailSchedule struct{}

var EmailSchedule = newEmailSchedule()

func newEmailSchedule() *emailSchedule {
	return &emailSchedule{}
}

func (o *emailSchedule) ListAll() []*entity.SendEmailSchedule {
	ctx, cursor := db.Find(bson.D{}, entity.SendEmailSchedule{})
	return mongodb.DecodeList(ctx, cursor, entity.SendEmailSchedule{})
}

func (o *emailSchedule) GetGameSendEmailSchedule(game *entity.Game) *entity.SendEmailSchedule {
	query := bson.D{{"symbol", game.Symbol}}
	ctx, cursor := db.FindOne(query, entity.SendEmailSchedule{})
	return mongodb.DecodeOne(ctx, cursor, entity.SendEmailSchedule{})
}

func (o *emailSchedule) SaveScheduleNextSendTime(game *entity.Game, sendTime int64) *entity.SendEmailSchedule {
	query := bson.D{{"symbol", game.Symbol}}
	update := bson.D{
		{"symbol", game.Symbol},
		{"email", game.DeveloperEmail},
		{"sendTime", sendTime},
	}
	db.UpdateMany(entity.SendEmailSchedule{}, query, update, options.Update().SetUpsert(true))
	return o.GetGameSendEmailSchedule(game)
}

func (o *emailSchedule) DeleteSchedule(game *entity.Game) {
	query := bson.D{{"symbol", game.Symbol}}
	db.DeleteOne(entity.SendEmailSchedule{}, query)
}
