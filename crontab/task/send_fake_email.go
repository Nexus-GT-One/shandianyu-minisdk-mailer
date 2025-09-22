package task

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"shandianyu-minisdk-mailer/util/randomutil"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/gomail.v2"
)

//go:embed file/names.json
var namesByte []byte

//go:embed file/emails.json
var emailsByte []byte

var namesArray []string
var emailsArray []string

func init() {
	json.Unmarshal(namesByte, &namesArray)
	json.Unmarshal(emailsByte, &emailsArray)

	registerTask("5 */5 * * * *", func(param map[string]any) {
		start := time.Now()

		// 判断达到了发送时间的就发邮件
		now := time.Now().UnixMilli()
		for _, schedule := range service.GameMailBox.ListAll() {
			// 还没到发送时间
			if schedule.SendTime <= 0 || schedule.SendTime > now {
				continue
			}

			// 到点发送
			sendFakeEmail(service.GameService.GetBySymbol(schedule.Symbol))
		}

		// 发送完成后，设定下次发送时间
		for _, game := range service.GameService.ListAll() {
			if !game.Enable || len(game.DeveloperEmail) <= 0 {
				service.GameMailBox.DeleteSchedule(game)
				continue
			}

			// 新邮箱注册7天内每天都要发一封邮件；7天后不需要再发
			sendTime := int64(0)
			registerDays := int((now - game.CreateTime) / 86400000)
			schedule := service.GameMailBox.GetGameGameMailBox(game)
			if schedule == nil {
				schedule = service.GameMailBox.SaveScheduleNextSendTime(game, sendTime)
			} else {
				sendTime = schedule.SendTime
			}
			if sendTime < now {
				if registerDays <= 0 {
					sendTime = now + int64(randomutil.Int(10800000, 21600000))
				} else if registerDays <= 7 {
					sendTime = now + int64(randomutil.Int(10800000, 86400000))
				} else {
					sendTime = 0
				}
				service.GameMailBox.SaveScheduleNextSendTime(game, sendTime)
			}

			if sendTime <= 0 {
				logger.Info("[%s] (%s) 注册距离现在 %d 天，不需要再发送邮件", game.Symbol, game.DeveloperEmail, registerDays)
			} else {
				logger.Info("[%s] (%s) 注册距离现在 %d 天，下次发送时间 %s", game.Symbol, game.DeveloperEmail, registerDays, time.UnixMilli(sendTime).Format(time.DateTime))
			}
		}
		logger.Info("本次邮件发送执行时间：%v", time.Since(start))
	})
}

func sendFakeEmail(game *entity.Game) {
	// 发件人信息
	from := "fishflash_mail@sina.com"
	password := "f3da51a030dfb484"

	// 收件人
	to := game.DeveloperEmail

	// SMTP服务器配置
	smtpHost := "smtp.sina.com"
	smtpPort := 465 // SSL 端口

	// 创建邮件
	m := gomail.NewMessage()
	body := emailsArray[randomutil.Int(0, len(emailsArray)-1)]
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(body))
	title := strings.TrimSpace(doc.Find("title").Text())
	m.SetHeader("From", fmt.Sprintf("%s <%s>", namesArray[randomutil.Int(0, len(namesArray)-1)], from))
	m.SetHeader("To", to)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(smtpHost, smtpPort, from, password)
	d.SSL = true
	logger.Info("发件人：%s\n收件人：%s\n标题：%s\n正文：%s", from, to, title, body)
	if err := d.DialAndSend(m); err != nil {
		logger.Info("游戏 [%s] 邮件发送失败：%v", game.Symbol, err)
	} else {
		logger.Info("游戏 [%s] 邮件发送成功！", game.Symbol)
	}
}
