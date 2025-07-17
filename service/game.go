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
	ctx, cursor := db.Find(query, entity.Game{})
	return mongodb.DecodeOne(ctx, cursor, entity.Game{})
}
