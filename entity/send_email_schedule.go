package entity

import (
	"shandianyu-minisdk-mailer/provider/mongodb"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SendEmailSchedule struct {
	Id       primitive.ObjectID `bson:"_id"`
	Symbol   string             `bson:"symbol" index:"{'name':'symbol','keys':{'symbol':-1},'unique':true}"`
	Email    string             `bson:"email"`
	SendTime int64              `bson:"sendTime"`
}

func init() {
	mongodb.EnsureIndex(SendEmailSchedule{})
}
