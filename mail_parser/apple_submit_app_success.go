package mail_parser

import (
	"regexp"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"strings"
)

// 提审app成功（国内账号）
type submitAppSuccessMailParser struct{}

func init() {
	registerImplement(&submitAppSuccessMailParser{})
}

func (o *submitAppSuccessMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "App Store Connect")
}

func (o *submitAppSuccessMailParser) checkTitle(title string) bool {
	return strings.Contains(title, "Thank You for Submitting Your App")
}

func (o *submitAppSuccessMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "We've received your app for review")
}

func (o *submitAppSuccessMailParser) parse(from, to, bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := service.GameService.GetByName(o.extractAppName(bodyText))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: service.GameService.GetAuditingVersion(oneGame),
		BuildNum:   service.GameService.GetLastSubmitBuildNum(oneGame),
		Status:     "提审app成功 (国内账号)",
		Content:    bodyText,
	}
}

func (o *submitAppSuccessMailParser) extractAppName(body string) string {
	re := regexp.MustCompile(`(?m)^We've received your app for review\.\s*\n(.+?)\n`)
	match := re.FindStringSubmatch(body)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func (o *submitAppSuccessMailParser) after(game *entity.Game, gameMail *entity.GameMail) {}
