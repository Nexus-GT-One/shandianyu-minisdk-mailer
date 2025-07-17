package mail_parser

import (
	"runtime/debug"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/util/secretutil"
	"strings"
	"time"
)

// 邮件样例文档
// https://shimo.im/sheets/KlkKvKD25wcLG4qd/MODOC

var mailParserImplementMap = make(map[string]IMailParser)

type IMailParser interface {
	checkTitle(title string) bool
	checkKeyword(bodyText string) bool
	parse(bodyText string) (*entity.Game, *entity.GameMail)
}

func ParseMail(title, receiveTime, bodyText string) (*entity.Game, *entity.GameMail) {
	for _, handler := range mailParserImplementMap {
		if !handler.checkTitle(title) || !handler.checkKeyword(bodyText) {
			continue
		}
		oneGame, gameMail := handler.parse(bodyText)
		gameMail.Title = title
		gameMail.MD5 = secretutil.MD5(bodyText)
		beijingLocation, _ := time.LoadLocation("Asia/Shanghai")
		createTime, _ := time.ParseInLocation(time.DateTime, receiveTime, beijingLocation)
		gameMail.ReceiveTime = createTime.UnixMilli()
		gameMail.Content = strings.TrimSpace(strings.ReplaceAll(gameMail.Content, "\n", "<br>")) // 换行符换成<br>方便在后台显示
		return oneGame, gameMail
	}

	return nil, nil
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
