package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SystemConfig struct {
	Id     primitive.ObjectID `bson:"_id"`
	Module string             `bson:"module" index:"{'name':'moduleAndKey','keys':{'module':-1,'key':-1},'unique':true}"`
	Key    string             `bson:"key"`
	Value  any                `bson:"value"`
}
