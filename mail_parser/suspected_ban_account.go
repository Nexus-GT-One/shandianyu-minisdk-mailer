package mail_parser

import (
	"regexp"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"strings"
)

// 疑似封号
type suspectedBanAccountMailParser struct{}

func init() {
	registerImplement(&suspectedBanAccountMailParser{})
}

func (o *suspectedBanAccountMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "App Store Connect")
}

func (o *suspectedBanAccountMailParser) checkTitle(title string) bool {
	return strings.Contains(title, `You have a message from App Review`)
}

func (o *suspectedBanAccountMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "app Apple ID: ")
}

func (o *suspectedBanAccountMailParser) parse(bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := service.GameService.GetByName(o.extractAppName(bodyText))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: findAuditingVersion(oneGame),
		Status:     "疑似封号",
		Content:    bodyText,
	}
}

func (o *suspectedBanAccountMailParser) extractAppName(body string) string {
	re := regexp.MustCompile(`app,\s*(.+?),\s*app Apple ID`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}
