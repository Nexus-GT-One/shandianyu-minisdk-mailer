package mail_parser

import (
	"regexp"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"strings"
)

// 机审4.3被拒
type wasRejectedMailParser struct{}

func init() {
	registerImplement(&wasRejectedMailParser{})
}

func (o *wasRejectedMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "App Store Connect")
}

func (o *wasRejectedMailParser) checkTitle(title string) bool {
	return strings.Contains(title, "Your App Review Feedback")
}

func (o *wasRejectedMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "was rejected")
}

func (o *wasRejectedMailParser) parse(bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := service.GameService.GetByName(o.extractAppName(bodyText))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: service.GameService.GetAuditingVersion(oneGame),
		Status:     "机审4.3被拒",
		Content:    bodyText,
	}
}

func (o *wasRejectedMailParser) extractAppName(body string) string {
	re := regexp.MustCompile(`(?m)^Changes\sneeded\.\s*\n([A-Za-z\s]+)\s*$`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return strings.TrimSpace(strings.ReplaceAll(matches[1], "iOS", ""))
	}
	return ""
}

func (o *wasRejectedMailParser) after(game *entity.Game, gameMail *entity.GameMail) {
	service.GameService.RecordRejected(game.BundleId, service.GameService.GetAuditingVersion(game))
}
