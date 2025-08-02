package mail_parser

import (
	"regexp"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"strings"
)

// 过审
type readyForDistributionMailParser struct{}

func init() {
	registerImplement(&readyForDistributionMailParser{})
}

func (o *readyForDistributionMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "App Store Connect")
}

func (o *readyForDistributionMailParser) checkTitle(title string) bool {
	return strings.Contains(title, "Ready for Distribution")
}

func (o *readyForDistributionMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "ready for distribution")
}

func (o *readyForDistributionMailParser) parse(bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := service.GameService.GetByName(o.extractAppName(bodyText))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: service.GameService.GetAuditingVersion(oneGame),
		Status:     "过审",
		Content:    bodyText,
	}
}

func (o *readyForDistributionMailParser) extractAppName(body string) string {
	re := regexp.MustCompile(`(?m)^App\sName:\s([A-Za-z\s]+)$`)
	match := re.FindStringSubmatch(body)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func (o *readyForDistributionMailParser) after(game *entity.Game, gameMail *entity.GameMail) {
	service.ApplicationService.CheckApplicationNewVersion(game, gameMail)
}
