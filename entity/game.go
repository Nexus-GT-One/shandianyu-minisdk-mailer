package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskConfigOld struct {
	// hide-隐藏条件;display-显示条件
	Attribute string `json:"attribute"`
	// watchVideo-看视频;mission-通关
	Type string `json:"type"`
	// 结合type来计算达到的数值
	Value int `json:"value"`
}

type PushCert struct {
	// 证书文件名
	Name string `bson:"name" json:"name"`
	// 推送证书内容，用app本身的aesKey加密后存储
	Content string `bson:"content" json:"content"`
	// 推送证书过期时间
	Expired string `bson:"expired" json:"expired"`
	// 推送监控机器人
	PushRobot string `bson:"pushRobot" json:"pushRobot"`
	// 是否开启监控
	PushRobotEnable bool `bson:"pushRobotEnable" json:"pushRobotEnable"`
}

type GmConfig struct {
	// GM开关
	Enable bool `bson:"enable" json:"enable"`
	// GM密码
	Password string `bson:"password" json:"password"`
	// 覆盖版本
	Version []string `bson:"version" json:"version"`
}

type WithdrawStage struct {
	// 货币符号
	Symbol string `json:"symbol"`
	// 可以提现金额
	Amount float64 `json:"amount"`
}

type DesAbTestConfig struct {
	AbVer     int64    `bson:"ab_ver" json:"ab_ver"`
	AbUid     []int    `bson:"ab_uid" json:"ab_uid"`
	AbRgar    []string `bson:"ab_rgar" json:"ab_rgar"`
	AbCnan    []string `bson:"ab_cnan" json:"ab_cnan"`
	AbOldtMin int      `bson:"ab_oldt_min" json:"ab_oldt_min"`
	AbOldtMax int      `bson:"ab_oldt_max" json:"ab_oldt_max"`
	AbDes     string   `bson:"ab_des" json:"ab_des"`
}

type Game struct {
	Id                                           primitive.ObjectID                    `bson:"_id"`
	Unique                                       string                                `bson:"unique"`
	BundleId                                     string                                `bson:"bundleId" index:"{'name':'bundleId','keys':{'bundleId':-1},'unique':true}"`
	AppId                                        string                                `bson:"appId"`
	AppStoreUrl                                  string                                `bson:"appStoreUrl"`
	Name                                         string                                `bson:"name"`
	Tag                                          []string                              `bson:"tag,omitempty"`
	Channel                                      string                                `bson:"channel"`
	Symbol                                       string                                `bson:"symbol" index:"{'name':'symbol','keys':{'symbol':-1},'unique':true}"`
	Enable                                       bool                                  `bson:"enable"`
	AesKey                                       string                                `bson:"aesKey"`
	AesIV                                        string                                `bson:"aesIV"`
	Developer                                    []string                              `bson:"developer"`
	Tester                                       []string                              `bson:"tester"`
	Producter                                    []string                              `bson:"producter"`
	DeveloperEmail                               string                                `bson:"developerEmail"`
	DeveloperEmailConfirm                        bool                                  `bson:"developerEmailConfirm"`
	CocosScriptKey                               string                                `bson:"cocosScriptKey"`
	CocosEncryptKey                              string                                `bson:"cocosEncryptKey"`
	CocosEncryptCode                             string                                `bson:"cocosEncryptCode"`
	SourcePin                                    string                                `bson:"sourcePin"`
	ProductionDomain                             string                                `bson:"productionDomain"`
	TestingDomain                                string                                `bson:"testingDomain"`
	IsUseLocalSource                             bool                                  `bson:"isUseLocalSource"`
	IsLogReportAllCocosHU                        bool                                  `bson:"isLogReportAllCocosHU"`
	IsLogReportAllCocosNormal                    bool                                  `bson:"isLogReportAllCocosNormal"`
	IsLogReportAllSDKIAP                         bool                                  `bson:"isLogReportAllSDKIAP"`
	IsLogReportAllSDKNormal                      bool                                  `bson:"isLogReportAllSDKNormal"`
	Webhook                                      string                                `bson:"webhook"`
	ReportWebhook                                string                                `bson:"reportWebhook"`
	ReportWebhookEnable                          bool                                  `bson:"reportWebhookEnable"`
	ReportTestingDevice                          string                                `bson:"reportTestingDevice,omitempty"`
	TermUrl                                      string                                `bson:"termUrl"`
	PolicyUrl                                    string                                `bson:"policyUrl"`
	PublishVersion                               string                                `bson:"publishVersion"`
	PublishTime                                  int64                                 `bson:"publishTime"`
	Header                                       map[string]string                     `bson:"header"`
	Report                                       map[string]string                     `bson:"report"`
	InitRequired                                 map[string]string                     `bson:"initRequired"`
	InitColumn                                   []map[string]string                   `bson:"initColumn"`
	Analytics                                    map[string]string                     `bson:"analytics"`
	AnalyticsResult                              map[string]string                     `bson:"analyticsResult"`
	Api                                          map[string]string                     `bson:"api"`
	ApiStatus                                    bool                                  `bson:"apiStatus"`
	DesVer                                       map[string]int64                      `bson:"desVer"`
	DesConfig                                    map[string]string                     `bson:"desConfig"`
	DesAbTestConfig                              map[string]map[string]DesAbTestConfig `bson:"desAbTestConfig"`
	H5BasicConfig                                map[string]string                     `bson:"h5BasicConfig"`
	H5WebConfig                                  map[string]string                     `bson:"h5WebConfig"`
	WebConfigVer                                 map[string]int64                      `bson:"webConfigVer"`
	LaunchInfoVer                                map[string]int64                      `bson:"launchInfoVer"`
	LaunchInfoConfig                             map[string]string                     `bson:"launchInfoConfig"`
	TaskConfig                                   []TaskConfigOld                       `bson:"taskConfig"`
	TaskConfigVer                                int64                                 `bson:"taskConfigVer"`
	WithdrawStage                                []WithdrawStage                       `bson:"withdrawStage"`
	WithdrawStageVer                             int64                                 `bson:"withdrawStageVer"`
	PushCert                                     PushCert                              `bson:"pushCert"`
	Hotfix                                       map[string]bool                       `bson:"hotfix"`
	Audit                                        map[string]bool                       `bson:"audit"`
	WwyEnable                                    bool                                  `bson:"wwyEnable"`
	WwyApiKey                                    string                                `bson:"wwyApiKey"`
	WwyApiHost                                   string                                `bson:"wwyApiHost"`
	WwySdkVer                                    string                                `bson:"wwySdkVer"`
	WwyApi                                       map[string]string                     `bson:"wwyApi"`
	WwyHeader                                    map[string]string                     `bson:"wwyHeader"`
	WwyColumn                                    map[string]string                     `bson:"wwyColumn"`
	WwyReport                                    map[string]string                     `bson:"wwyReport"`
	WwyRiskResult                                map[string]string                     `bson:"wwyRiskResult"`
	WwyWebhook                                   string                                `bson:"wwyWebhook"`
	WwyWebhookEnable                             bool                                  `bson:"wwyWebhookEnable"`
	WwyTestingDevice                             string                                `bson:"wwyTestingDevice,omitempty"`
	WwyUploadParameters                          map[string]any                        `bson:"wwyUploadParameters"`
	AdParametersVer                              int64                                 `bson:"adParametersVer"`
	AdjustKey                                    string                                `bson:"adjustKey"`
	ApplovinKey                                  string                                `bson:"applovinKey"`
	ApplovinOpenADID                             string                                `bson:"applovinOpenADID"`
	ApplovinBannerADID                           string                                `bson:"applovinBannerADID"`
	ApplovinInterstitialADID                     string                                `bson:"applovinInterstitialADID"`
	ApplovinRewardADID                           string                                `bson:"applovinRewardADID"`
	MonitorEnable                                bool                                  `bson:"monitorEnable"`
	BSideSuccess                                 string                                `bson:"bSideSuccess"`
	RewardAd                                     string                                `bson:"rewardAd"`
	H5Ad                                         string                                `bson:"h5Ad"`
	GmConfig                                     GmConfig                              `bson:"gmConfig"`
	AppPurchaseEnable                            bool                                  `bson:"appPurchaseEnable"`
	AppPurchaseApi                               map[string]string                     `bson:"appPurchaseApi"`
	AppPurchaseColumn                            map[string]string                     `bson:"appPurchaseColumn"`
	AppPurchaseStoreKitV2KeyId                   string                                `bson:"appPurchaseStoreKitV2KeyId"`
	AppPurchaseStoreKitV2IssuerId                string                                `bson:"appPurchaseStoreKitV2IssuerId"`
	AppPurchaseStoreKitV2SubscriptionKey         string                                `bson:"appPurchaseStoreKitV2SubscriptionKey"`
	AppPurchaseStoreKitV2SubscriptionKeyFileName string                                `bson:"appPurchaseStoreKitV2SubscriptionKeyFileName"`
	CreateTime                                   int64                                 `bson:"createTime"`
	UpdateTime                                   int64                                 `bson:"updateTime"`
}
