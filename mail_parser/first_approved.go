package mail_parser

import (
	"regexp"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"strings"
)

// 首次过审
type firstApprovedMailParser struct{}

func init() {
	registerImplement(&firstApprovedMailParser{})
}

func (o *firstApprovedMailParser) checkTitle(title string) bool {
	return strings.Contains(title, "Welcome to the App Store")
}

func (o *firstApprovedMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "Distributing on the App Store")
}

func (o *firstApprovedMailParser) parse(bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := service.GameService.GetByName(o.extractAppName(bodyText))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: findAuditingVersion(oneGame),
		Status:     "首次过审",
		Content:    bodyText,
	}
}

func (o *firstApprovedMailParser) extractAppName(body string) string {
	re := regexp.MustCompile(`(?i)Congratulations!([A-Z][a-zA-Z]*(?: [A-Z][a-zA-Z]*)*)iOS`)
	match := re.FindStringSubmatch(body)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}
