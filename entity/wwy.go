package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Wwy struct {
	Id              primitive.ObjectID `bson:"_id" index:"{'name':'updateTime','keys':{'updateTime':-1}}"`
	BundleId        string             `bson:"bundleId" index:"{'name':'bundleIdAndDevice','keys':{'bundleId':-1,'device':-1},'unique':true}"`
	Symbol          string             `bson:"symbol" index:"{'name':'symbol','keys':{'symbol':-1,'uid':-1}}"`
	Device          string             `bson:"device" index:"{'name':'device','keys':{'device':-1}}"`
	Uid             int64              `bson:"uid" index:"{'name':'uid','keys':{'uid':-1}}"`
	Idfa            string             `bson:"idfa"`
	Idfv            string             `bson:"idfv"`
	Token           string             `bson:"token"`
	Channel         string             `bson:"channel"`
	AdjustId        string             `bson:"adjustId"`
	MobileModel     string             `bson:"mobileModel"`
	InviteCode      string             `bson:"inviteCode" index:"{'name':'inviteCode','keys':{'inviteCode':-1}}"`
	InviteUrl       string             `bson:"inviteUrl"`
	AppVersion      string             `bson:"appVersion"`
	FirstActiveTime int64              `bson:"firstActiveTime"`
	RegTime         int64              `bson:"regTime"`
	CreateTime      int64              `bson:"createTime"`
	UpdateTime      int64              `bson:"updateTime"`
}
