package service

import (
	_ "embed"
	"go.mongodb.org/mongo-driver/bson"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/provider/mongodb"
)

type gameService struct{}

var GameService = newGameService()

func newGameService() *gameService {
	return &gameService{}
}

func (a *gameService) GetByName(name string) *entity.Game {
	query := bson.D{{"name", name}}
	ctx, cursor := db.FindOne(query, entity.Game{})
	return mongodb.DecodeOne(ctx, cursor, entity.Game{})
}

func (a *gameService) GetByDeveloperEmail(developerEmail string) *entity.Game {
	query := bson.D{{"developerEmail", developerEmail}}
	ctx, cursor := db.FindOne(query, entity.Game{})
	return mongodb.DecodeOne(ctx, cursor, entity.Game{})
}
