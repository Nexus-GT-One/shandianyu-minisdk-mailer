package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"shandianyu-minisdk-mailer/provider/mongodb"
)

type GameDispatchedVersion struct {
	Id         primitive.ObjectID `bson:"_id" json:"id"`
	Symbol     string             `bson:"symbol" index:"{'name':'symbol','keys':{'symbol':-1},'unique':true}"`
	AppVersion string             `bson:"appVersion"`
}

func init() {
	mongodb.EnsureIndex(GameDispatchedVersion{})
}
