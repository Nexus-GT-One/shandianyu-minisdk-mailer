package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"shandianyu-minisdk-mailer/provider/mongodb"
)

type GameMail struct {
	Id          primitive.ObjectID `bson:"_id"`
	Symbol      string             `bson:"symbol" index:"{'name':'symbol','keys':{'symbol':-1,'appVersion':-1}}"`
	AppVersion  string             `bson:"appVersion"`
	Developer   string             `bson:"developer"`
	Title       string             `bson:"title"`
	Status      string             `bson:"status"`
	Content     string             `bson:"content"`
	MD5         string             `bson:"md5" index:"{'name':'md5','keys':{'developer':-1,'md5':-1},'unique':true}"`
	ReceiveTime int64              `bson:"receiveTime"`
}

func init() {
	mongodb.EnsureIndex(GameMail{})
}
