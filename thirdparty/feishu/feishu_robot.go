package feishu

import (
	"encoding/json"
	"fmt"
	queues "github.com/eapache/queue"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"runtime"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/provider/logger"
	"shandianyu-minisdk-mailer/provider/mongodb"
	"shandianyu-minisdk-mailer/util/arrayutil"
	"shandianyu-minisdk-mailer/util/httputil"
	"shandianyu-minisdk-mailer/util/systemutil"
	"strings"
	"sync"
)

//AWS报警群（AWS的监听机的都发去这里，公司的大佬全部在这个群）：
//https://open.feishu.cn/open-apis/bot/v2/hook/6f9574ec-633b-422d-8b50-af658afa45e3
//R报警群（R的监听机的都发去这里，这个群只有我跟你）：
//https://open.feishu.cn/open-apis/bot/v2/hook/4882de36-de41-4fcc-a4c4-a1b186dae668
//需求分配群机器人：
//https://open.feishu.cn/open-apis/bot/v2/hook/5aa4b733-115c-4aac-9eb4-b78ef599568a
//需求开发机器人
//https://open.feishu.cn/open-apis/bot/v2/hook/611fd705-4a88-4c6e-8f6b-aff9fb2abffb
//前期埋点告警机器人
//https://open.feishu.cn/open-apis/bot/v2/hook/d8020653-9958-432d-ba9b-b72dbe0c7ea4
//归因率报警群
//https://open.feishu.cn/open-apis/bot/v2/hook/c7588d71-bcdb-47a1-ac61-06fe307b358a
//iOS邮件通知群
//https://open.feishu.cn/open-apis/bot/v2/hook/484a4255-fe36-426e-a789-7c9a0c99d1c5
//GP邮件通知群
//https://open.feishu.cn/open-apis/bot/v2/hook/4ae33f52-6e9e-451c-a463-a6084cace9af

type FeishuRobot struct{ url string }

type queueMessage struct {
	url        string
	message    string
	userOpenId string
}

type UserItem struct {
	Mobile string `json:"mobile"`
	OpenId string `json:"user_id"`
}
type UserListResponse struct {
	Code int `json:"code"`
	Data struct {
		UserList []UserItem `json:"user_list"`
	} `json:"data"`
	Msg string `json:"msg"`
}

var queue = queues.New()
var lock = sync.NewCond(&sync.Mutex{})

func init() {
	systemutil.Goroutine(func() {
		for {
			obj := dequeue()
			data := map[string]any{
				"msg_type": "text",
				"content":  map[string]string{"text": obj.message},
			}
			if len(obj.userOpenId) > 0 {
				data["content"] = map[string]string{"text": fmt.Sprintf(`%s <at user_id="%s"></at>`, obj.message, obj.userOpenId)}
			}
			body, _ := json.Marshal(data)
			response := make(map[string]any)
			httputil.Post(obj.url, make(map[string]string), body, &response)
		}
	})
}

func NewFeishuRobot() *FeishuRobot {
	return &FeishuRobot{url: "https://open.feishu.cn/open-apis/bot/v2/hook/6f9574ec-633b-422d-8b50-af658afa45e3"}
}

func AttributeRobot() *FeishuRobot {
	return &FeishuRobot{url: "https://open.feishu.cn/open-apis/bot/v2/hook/c7588d71-bcdb-47a1-ac61-06fe307b358a"}
}

func AdminRobot() *FeishuRobot {
	return &FeishuRobot{url: "https://open.feishu.cn/open-apis/bot/v2/hook/4882de36-de41-4fcc-a4c4-a1b186dae668"}
}

func DemandRobot() *FeishuRobot {
	return &FeishuRobot{url: "https://open.feishu.cn/open-apis/bot/v2/hook/5aa4b733-115c-4aac-9eb4-b78ef599568a"}
}

func CheckPointRobot() *FeishuRobot {
	return &FeishuRobot{url: "https://open.feishu.cn/open-apis/bot/v2/hook/d8020653-9958-432d-ba9b-b72dbe0c7ea4"}
}

func IOSMailRobot() *FeishuRobot {
	return &FeishuRobot{url: "https://open.feishu.cn/open-apis/bot/v2/hook/484a4255-fe36-426e-a789-7c9a0c99d1c5"}
}
func GPMailRobot() *FeishuRobot {
	return &FeishuRobot{url: "https://open.feishu.cn/open-apis/bot/v2/hook/4ae33f52-6e9e-451c-a463-a6084cace9af"}
}

func DevRobot() *FeishuRobot {
	return &FeishuRobot{url: "https://open.feishu.cn/open-apis/bot/v2/hook/611fd705-4a88-4c6e-8f6b-aff9fb2abffb"}
}

func (f *FeishuRobot) SendRobotMessage(message string, at ...string) {
	if reflect.DeepEqual(runtime.GOOS, "windows") {
		logger.GetLogger().Info("发送飞书消息：%s", message)
		return
	}

	lock.L.Lock()
	defer lock.L.Unlock()
	msg := queueMessage{url: f.url, message: message}
	if len(at) > 0 {
		msg.userOpenId = GetUserOpenId(strings.TrimSpace(at[0])).OpenId
	}
	queue.Add(msg)
	lock.Signal()
}

func dequeue() queueMessage {
	lock.L.Lock()
	defer lock.L.Unlock()
	for queue.Length() == 0 {
		lock.Wait()
	}

	return queue.Remove().(queueMessage)
}

func (f *FeishuRobot) SendRobotInteractive(title, content string) {
	body, _ := json.Marshal(map[string]any{
		"msg_type": "interactive",
		"card": map[string]any{
			"schema": "2.0",
			"header": map[string]any{
				"title": map[string]any{
					"tag":     "plain_text",
					"content": title,
				},
				"template": "blue",
				"padding":  "12px 12px 12px 12px",
			},
			"body": map[string]any{
				"direction": "vertical",
				"padding":   "12px 12px 12px 12px",
				"elements": []map[string]any{{
					"tag": "div",
					"text": map[string]any{
						"tag":     "lark_md",
						"content": content,
					},
				}},
			},
		},
	})
	httputil.Post(f.url, make(map[string]string), body, &map[string]any{})
}

// https://open.feishu.cn/document/server-docs/contact-v3/user/batch_get_id?appId=cli_a72420c4753fd00e
func GetUserOpenId(email string) UserItem {
	url := "https://open.feishu.cn/open-apis/contact/v3/users/batch_get_id"
	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", getSdyAppToken())}
	body, _ := json.Marshal(map[string]any{"mobiles": []string{getAllEmployees()[email].Mobile}})
	response := &UserListResponse{}
	httputil.Post(url, headers, body, &response)
	return arrayutil.First(response.Data.UserList)
}

// https://open.feishu.cn/document/server-docs/authentication-management/access-token/tenant_access_token_internal
func getSdyAppToken() string {
	ctx, cursor := mongodb.GetInstance().FindOne(bson.D{{"module", "feishu"}, {"key", "sdyApp"}}, entity.SystemConfig{})
	conf := mongodb.DecodeOne(ctx, cursor, entity.SystemConfig{})
	b, _ := json.Marshal(conf.Value)
	x := bson.D{}
	json.Unmarshal(b, &x)
	url := "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
	body, _ := json.Marshal(map[string]any{"app_id": x.Map()["appId"], "app_secret": x.Map()["appSecret"]})
	response := make(map[string]any)
	httputil.Post(url, make(map[string]string), body, &response)
	return strings.TrimSpace(fmt.Sprintf("%v", response["tenant_access_token"]))
}

// 获取所有管理员
func getAllEmployees() map[string]*entity.Admin {
	result := make(map[string]*entity.Admin)
	ctx, cursor := mongodb.GetInstance().Find(bson.D{}, entity.Admin{}, &options.FindOptions{Sort: bson.D{{"email", 1}}})
	for _, admin := range mongodb.DecodeList(ctx, cursor, entity.Admin{}) {
		result[admin.Email] = admin
	}
	return result
}
