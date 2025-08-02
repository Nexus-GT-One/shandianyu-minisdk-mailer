package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuditTrack struct {
	Id              primitive.ObjectID `bson:"_id"`
	BundleId        string             `bson:"bundleId"`
	Device          string             `bson:"device"`
	ActionType      string             `bson:"actionType"`
	Value           string             `bson:"value"`
	AppVersion      string             `bson:"appVersion"`
	BuildNum        int                `bson:"buildNum,omitempty"`
	MobileOsVersion string             `bson:"mobileOsVersion"`
	MobileLanguage  string             `bson:"MobileLanguage"`
	NetType         string             `bson:"netType"`
	CountryCode     string             `bson:"countryCode"`
	Country         string             `bson:"country"`
	City            string             `bson:"city"`
	Idfv            string             `bson:"idfv"`
	VpnOrProxy      int                `bson:"vpnOrProxy"`
	AutoReleaseDays int                `bson:"autoReleaseDays"`
	MobileModel     string             `bson:"mobileModel"`
	Ip              string             `bson:"ip"`
	Remark          string             `bson:"remark"`
	CreateTime      int64              `bson:"createTime"`
}
