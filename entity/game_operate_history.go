package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameOperateHistory struct {
	Id         primitive.ObjectID `bson:"_id"`
	GameId     string             `bson:"gameId" index:"{'name':'gameId','keys':{'gameId':-1}}"`
	Operator   string             `bson:"operator"`
	Type       string             `bson:"type"`
	Remark     string             `bson:"remark"`
	CreateTime int64              `bson:"createTime"`
}
