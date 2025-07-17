package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Admin struct {
	Id        primitive.ObjectID `bson:"_id" json:"id"`
	Email     string             `bson:"email" json:"email" index:"{'name':'email','keys':{'email':-1},'unique':true}"`
	Name      string             `bson:"name" json:"name"`
	Mobile    string             `bson:"mobile" json:"mobile"`
	Password  string             `bson:"password" json:"password,omitempty"`
	Authority []string           `bson:"authority" json:"authority"`
}
