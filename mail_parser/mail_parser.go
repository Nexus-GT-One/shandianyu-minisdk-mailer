package mail_parser

import (
	"regexp"
	"runtime/debug"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"shandianyu-minisdk-mailer/util/secretutil"
	"strings"
	"time"
)

// 邮件样例文档
// https://shimo.im/sheets/KlkKvKD25wcLG4qd/MODOC

var mailParserImplementMap = make(map[string]IMailParser)

type IMailParser interface {
	checkFrom(from string) bool
	checkTitle(title string) bool
	checkKeyword(bodyText string) bool
	parse(bodyText string) (*entity.Game, *entity.GameMail)
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

	oneGame, gameMail := handler.parse(bodyText)
	if oneGame == nil {
		oneGame = service.GameService.GetByDeveloperEmail(to)
	}
	gameMail.Title = title
	gameMail.MD5 = secretutil.MD5(bodyText)
	beijingLocation, _ := time.LoadLocation("Asia/Shanghai")
	createTime, _ := time.ParseInLocation(time.DateTime, receiveTime, beijingLocation)
	gameMail.ReceiveTime = createTime.UnixMilli()
	gameMail.Content = strings.TrimSpace(strings.ReplaceAll(gameMail.Content, "\n", "<br>")) // 换行符换成<br>方便在后台显示
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

func findAuditingVersion(oneGame *entity.Game) string {
	for appVersion, isAudit := range oneGame.Audit {
		if !isAudit {
			continue
		}
		return appVersion
	}

	return ""
}

func extractDeveloperEmail(body string) string {
	re := regexp.MustCompile(`(?i)email:\s*([\w\.-]+@[\w\.-]+\.\w+)`)
	match := re.FindStringSubmatch(body)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}

	return ""
}
