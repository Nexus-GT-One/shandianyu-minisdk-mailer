package mail_parser

import (
	"regexp"
	"runtime/debug"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"shandianyu-minisdk-mailer/util/arrayutil"
	"shandianyu-minisdk-mailer/util/secretutil"
	"strings"
	"time"
)

// 邮件样例文档
// https://shimo.im/sheets/KlkKvKD25wcLG4qd/MODOC

var mailParserImplementMap = make(map[string]IMailParser)
var storeMap = map[string]string{"iOS": "AppStore", "GooglePlay": "GooglePlay", "Samsung": "galaxyStore"}

type IMailParser interface {
	checkFrom(from string) bool
	checkTitle(title string) bool
	checkKeyword(bodyText string) bool
	parse(from, to, bodyText string) (*entity.Game, *entity.GameMail)
	after(*entity.Game, *entity.GameMail)
}

func ParseMail(title, from, to, receiveTime, bodyText string) (*entity.Game, *entity.GameMail) {
	for _, handler := range mailParserImplementMap {
		oneGame, gameMail := baseParseMail(handler, title, from, to, receiveTime, bodyText)
		if oneGame == nil || gameMail == nil {
			continue
		}
		return oneGame, gameMail
	}

	return nil, nil
}

func ParseOtherMail(title, from, to, receiveTime, bodyText string) (*entity.Game, *entity.GameMail) {
	return baseParseMail(&OtherAppstoreMailParser{}, title, from, to, receiveTime, bodyText)
}

func baseParseMail(handler IMailParser, title, from, to, receiveTime, bodyText string) (*entity.Game, *entity.GameMail) {
	if !handler.checkFrom(from) || !handler.checkTitle(title) || !handler.checkKeyword(bodyText) {
		return nil, nil
	}

	oneGame, gameMail := handler.parse(from, to, bodyText)
	if gameMail == nil {
		return nil, nil
	}
	if oneGame == nil {
		oneGame = arrayutil.First(service.GameService.GetByDeveloperEmail(to))
	}
	gameMail.Title = title
	gameMail.MD5 = secretutil.MD5(bodyText)
	beijingLocation, _ := time.LoadLocation("Asia/Shanghai")
	createTime, _ := time.ParseInLocation(time.DateTime, receiveTime, beijingLocation)
	gameMail.ReceiveTime = createTime.UnixMilli()

	// 换行符换成<br>方便在后台显示
	gameMail.Content = strings.TrimSpace(strings.ReplaceAll(gameMail.Content, "\n", "<br>"))

	// 执行邮件解析完成之后要执行的逻辑
	handler.after(oneGame, gameMail)
	return oneGame, gameMail
}

func registerImplement(instance any) {
	mailParserImplementMap[getName()] = instance.(IMailParser)
}

func getName() string {
	programStack := strings.Split(string(debug.Stack()), "\n")
	programPath := strings.TrimSpace(programStack[len(programStack)-2])
	programPath = strings.TrimSpace(strings.Split(programPath, " ")[0])
	programPath = strings.TrimSpace(programPath[strings.LastIndex(programPath, "/")+1 : strings.LastIndex(programPath, ":")])
	return programPath[:len(programPath)-3]
}

func extractDeveloperEmail(body string) string {
	re := regexp.MustCompile(`(?i)email:\s*([\w\.-]+@[\w\.-]+\.\w+)`)
	match := re.FindStringSubmatch(body)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}

	return ""
}

func extractBundleId(body string) string {
	re := regexp.MustCompile(`\(([a-zA-Z][a-zA-Z0-9_]*(?:\.[a-zA-Z0-9_]+)+)\)`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}

	return ""
}
