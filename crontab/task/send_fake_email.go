package task

import (
	"fmt"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"shandianyu-minisdk-mailer/thirdparty/aigcbest"
	"shandianyu-minisdk-mailer/util/arrayutil"
	"shandianyu-minisdk-mailer/util/randomutil"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

func init() {
	registerTask("5 */5 * * * *", func(param map[string]any) {
		start := time.Now()

		// 判断达到了发送时间的就发邮件
		now := time.Now().UnixMilli()
		for _, schedule := range service.EmailSchedule.ListAll() {
			// 还没到发送时间
			if schedule.SendTime > now {
				continue
			}

			// 到点发送
			sendFakeEmail(service.GameService.GetBySymbol(schedule.Symbol))
		}

		// 发送完成后，设定下次发送时间
		for _, game := range service.GameService.ListAll() {
			if !game.Enable || len(game.DeveloperEmail) <= 0 {
				service.EmailSchedule.DeleteSchedule(game)
				continue
			}

			// 新邮箱注册7天内每天都要发一封邮件；7天后随机一个时间偶尔发一封
			sendTime := int64(0)
			registerDays := int((now - game.CreateTime) / 86400000)
			schedule := service.EmailSchedule.GetGameSendEmailSchedule(game)
			if schedule == nil {
				schedule = service.EmailSchedule.SaveScheduleNextSendTime(game, sendTime)
			} else {
				sendTime = schedule.SendTime
			}
			if sendTime < now {
				if registerDays <= 0 {
					sendTime = now + int64(randomutil.Int(10800000, 21600000))
				} else if registerDays <= 7 {
					sendTime = now + int64(randomutil.Int(10800000, 86400000))
				} else if registerDays <= 14 {
					sendTime = now + int64(randomutil.Int(86400000*3, 86400000*5))
				} else {
					sendTime = now + int64(randomutil.Int(86400000*6, 86400000*14))
				}
				service.EmailSchedule.SaveScheduleNextSendTime(game, sendTime)
			}
			logger.Info("[%s] (%s) 注册距离现在 %d 天，下次发送时间 %s", game.Symbol, game.DeveloperEmail, registerDays, time.UnixMilli(sendTime).Format(time.DateTime))
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
	title := randomTitleByAI()
	body := randomEmailByAI()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", randomNameByAI(), from))
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

func randomEmailByAI() string {
	letters := []int{300, 400, 500, 600, 800, 1000, 1500, 2000, 2500, 3000}
	content := fmt.Sprintf("帮我生成一封邮件，要求内容完整，语句通顺，有一定的实际意义，英文版，%d个单词，要求只返回html代码内容。", letters[randomutil.Int(0, len(letters)-1)])
	html := arrayutil.First(aigcbest.ChatCompletion(content).Choices).Message.Content
	if len(html) <= 0 {
		return ""
	}
	html = strings.TrimSpace(html[strings.Index(html, "`")+7 : strings.LastIndex(html, "`")-3])
	return html
}

func randomNameByAI() string {
	return arrayutil.First(aigcbest.ChatCompletion("帮我生成一个英文名，有姓有名，既可以美式，也可以欧式，西班牙语、法语，都可以，直接返回即可").Choices).Message.Content
}

func randomTitleByAI() string {
	return strings.ReplaceAll(arrayutil.First(aigcbest.ChatCompletion("帮我生成一个英文邮件标题，直接返回即可").Choices).Message.Content, `"`, "")
}
