package mail_parser

import (
	"regexp"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"strings"
)

// 过审
type reviewApprovedMailParser struct{}

func init() {
	registerImplement(&reviewApprovedMailParser{})
}

func (o *reviewApprovedMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "App Store Connect")
}

func (o *reviewApprovedMailParser) checkTitle(title string) bool {
	return strings.Contains(title, "Review of your submission is complete")
}

func (o *reviewApprovedMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "eligible for distribution")
}

func (o *reviewApprovedMailParser) parse(bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := service.GameService.GetByName(o.extractAppName(bodyText))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: findAuditingVersion(oneGame),
		Status:     "过审",
		Content:    bodyText,
	}
}

func (o *reviewApprovedMailParser) extractAppName(body string) string {
	re := regexp.MustCompile(`(?m)^App\sName:\s([A-Za-z\s]+)$`)
	match := re.FindStringSubmatch(body)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}
