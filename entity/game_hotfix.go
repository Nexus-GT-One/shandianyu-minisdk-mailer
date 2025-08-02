package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotfixConfig struct {
	Build   int64    `bson:"build" json:"build"`
	Uid     []int    `bson:"uid" json:"uid"`
	Country []string `bson:"country" json:"country"`
	Channel []string `bson:"channel" json:"channel"`
	DayMin  int      `bson:"dayMin" json:"dayMin"`
	DayMax  int      `bson:"dayMax" json:"dayMax"`
	Content string   `bson:"content" json:"content"`
}

type WebConfig struct {
	Build            int64  `bson:"build" json:"build"`
	BasicConfig      string `bson:"basicConfig" json:"basicConfig"`
	WebConfig        string `bson:"webConfig" json:"webConfig"`
	LaunchInfoConfig string `bson:"launchInfoConfig" json:"launchInfoConfig"`
}

type TaskConfig struct {
	Build      int64  `bson:"build" json:"build"`
	TaskConfig string `bson:"taskConfig" json:"taskConfig"`
}

type WithdrawConfig struct {
	Build          int64  `bson:"build" json:"build"`
	WithdrawConfig string `bson:"withdrawConfig" json:"withdrawConfig"`
}

type ShareContentConfig struct {
	Build              int64  `bson:"build" json:"build"`
	ShareContentConfig string `bson:"shareContentConfig" json:"shareContentConfig"`
}

type InvitePeopleConfig struct {
	Build              int64  `bson:"build" json:"build"`
	InvitePeopleConfig string `bson:"invitePeopleConfig" json:"invitePeopleConfig"`
}

type InvitePeopleConfigItem struct {
	Minimum     float64 `json:"minimum"`
	Maximum     float64 `json:"maximum"`
	Probability int     `json:"probability"`
}

type InviteRewardConfig struct {
	Build              int64  `bson:"build" json:"build"`
	InviteRewardConfig string `bson:"inviteRewardConfig" json:"inviteRewardConfig"`
}

type GameHotfix struct {
	Id primitive.ObjectID `bson:"_id" json:"id"`
	// 关联游戏id
	GameId string `bson:"gameId" json:"gameId" index:"{'name':'gameIdAndVersion','keys':{'gameId':-1,'version':-1},'unique':true}"`
	// 热更版本号
	Version string `bson:"version" json:"version"`
	// 关联应用版本号
	AppVersion []string `bson:"appVersion" json:"appVersion"`
	// 热更是否开启
	Enable bool `bson:"enable" json:"enable"`
	// 热更配置内容 (配置项目 => 配置内容)
	Value map[string]any `bson:"value" json:"value"`
}
