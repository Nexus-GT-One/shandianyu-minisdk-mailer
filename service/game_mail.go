package service

import (
	"go.mongodb.org/mongo-driver/bson"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/provider/mongodb"
)

type gameMailService struct{}

var GameMailService = newGameMailService()

func newGameMailService() *gameMailService {
	return &gameMailService{}
}

func (s *gameMailService) FindByMd5(md5 string) *entity.GameMail {
	ctx, cursor := db.FindOne(bson.D{{"md5", md5}}, entity.GameMail{})
	return mongodb.DecodeOne(ctx, cursor, entity.GameMail{})
}
