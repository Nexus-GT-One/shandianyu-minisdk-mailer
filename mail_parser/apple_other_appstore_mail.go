package mail_parser

import (
	"regexp"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"shandianyu-minisdk-mailer/util/arrayutil"
	"strings"
)

// 未能识别的苹果邮件类型
type OtherAppstoreMailParser struct{}

func (o *OtherAppstoreMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "App Store") || strings.Contains(from, "Apple")
}

func (o *OtherAppstoreMailParser) checkTitle(title string) bool {
	return true
}

func (o *OtherAppstoreMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "App Store")
}

func (o *OtherAppstoreMailParser) parse(from, to, bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := service.GameService.GetByName(o.extractAppName(bodyText))
	if oneGame == nil {
		oneGame = arrayutil.First(service.GameService.GetByDeveloperEmail(extractDeveloperEmail(bodyText)))
		if oneGame == nil {
			oneGame = &entity.Game{}
		}
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: o.findAuditingVersion(bodyText),
		Content:    bodyText,
	}
}

func (o *OtherAppstoreMailParser) extractAppName(body string) string {
	re := regexp.MustCompile(`(?m)^(.+)\s*\nApp Apple ID\s+\d+\s*\nVersion\s+([^\s]+)`)
	match := re.FindStringSubmatch(body)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func (o *OtherAppstoreMailParser) findAuditingVersion(body string) string {
	re := regexp.MustCompile(`(?m)^(.+)\s*\nApp Apple ID\s+\d+\s*\nVersion\s+([^\s]+)`)
	match := re.FindStringSubmatch(body)
	if len(match) > 1 {
		return strings.TrimSpace(match[2])
	}
	return ""
}

func (o *OtherAppstoreMailParser) after(game *entity.Game, gameMail *entity.GameMail) {}
