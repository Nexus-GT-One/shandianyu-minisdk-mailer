package aigcbest

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/provider/logger"
	"shandianyu-minisdk-mailer/provider/mongodb"
	"shandianyu-minisdk-mailer/thirdparty/feishu"
	"shandianyu-minisdk-mailer/util/httputil"
	"shandianyu-minisdk-mailer/util/maputil"
	"shandianyu-minisdk-mailer/util/systemutil"
	"strconv"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatCompletionResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

type UserInfo struct {
	Data struct {
		Id               int         `json:"id"`
		Username         string      `json:"username"`
		Password         string      `json:"password"`
		DisplayName      string      `json:"display_name"`
		Role             int         `json:"role"`
		Status           int         `json:"status"`
		Email            string      `json:"email"`
		GithubId         string      `json:"github_id"`
		WechatId         string      `json:"wechat_id"`
		TelegramId       string      `json:"telegram_id"`
		VerificationCode string      `json:"verification_code"`
		AccessToken      interface{} `json:"access_token"`
		Quota            int         `json:"quota"`
		UsedQuota        int         `json:"used_quota"`
		RequestCount     int         `json:"request_count"`
		Group            string      `json:"group"`
		AffCode          string      `json:"aff_code"`
		AffCount         int         `json:"aff_count"`
		AffQuota         int         `json:"aff_quota"`
		AffHistoryQuota  int         `json:"aff_history_quota"`
		InviterId        int         `json:"inviter_id"`
		DeletedAt        interface{} `json:"DeletedAt"`
	} `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

func ChatCompletion(content string) *ChatCompletionResponse {
	res := &ChatCompletionResponse{}
	if reflect.DeepEqual(runtime.GOOS, "windows") {
		return res
	}

	headers := map[string]string{"Authorization": "Bearer sk-lwZJ0mBUKGuphTpdF1dj2IwIt9pnkpr5wVI0MRRnUcTan47k"}
	body, _ := json.Marshal(map[string]any{
		"model": "gpt-4o",
		"messages": []map[string]any{{
			"role":    "user",
			"content": content,
		}}})
	httputil.Post("https://api2.aigcbest.top/v1/chat/completions", headers, body, &res)

	// 检查余额，<=$1发报警
	systemutil.Goroutine(func() {
		balance := getBalance()
		if balance > 0 && balance <= 1 {
			feishu.NewFeishuRobot().SendRobotMessage(fmt.Sprintf("钱多多API账号剩余$%f，请尽快充值！", balance))
		}
	})
	return res
}

func getBalance() float64 {
	userInfo := UserInfo{}
	rate, _ := decimal.NewFromString("500000")
	ctx, cursor := mongodb.GetInstance().FindOne(bson.D{{"module", "ai"}, {"key", "aigcbestApi"}}, entity.SystemConfig{})
	conf := mongodb.DecodeOne(ctx, cursor, entity.SystemConfig{})
	url := fmt.Sprintf("%v", maputil.GetValueFromMap(conf.Value.(primitive.D).Map(), "url", ""))
	cookie := fmt.Sprintf("%v", maputil.GetValueFromMap(conf.Value.(primitive.D).Map(), "cookie", ""))
	newApiUser := fmt.Sprintf("%v", maputil.GetValueFromMap(conf.Value.(primitive.D).Map(), "new-api-user", ""))
	headers := map[string]string{"new-api-user": newApiUser, "Cookie": cookie, "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36"}
	json.Unmarshal(httputil.Get(url, headers), &userInfo)
	if !userInfo.Success || strconv.Itoa(userInfo.Data.Id) != newApiUser {
		go feishu.NewFeishuRobot().SendRobotMessage("钱多多API账号监控异常，请尽快更新cookie和new-api-user！")
		return -1
	}
	quota, _ := decimal.NewFromString(strconv.Itoa(userInfo.Data.Quota))
	balance := quota.Div(rate).Round(2).InexactFloat64()
	logger.GetLogger().Info("aigcbest balance is: $%f", balance)
	return balance
}
