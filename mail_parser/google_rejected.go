package mail_parser

import (
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"strings"
)

// 应用拒审
type rejectedMailParser struct{}

func init() {
	registerImplement(&rejectedMailParser{})
}

func (o *rejectedMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "Google Play Support")
}

func (o *rejectedMailParser) checkTitle(title string) bool {
	return strings.Contains(title, "Your app is not compliant")
}

func (o *rejectedMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "Rejected")
}

func (o *rejectedMailParser) parse(from, to, bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := service.GameService.GetByBundleId(extractBundleId(bodyText))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: oneGame.PublishVersion,
		Status:     "拒审",
		Content:    bodyText,
	}
}

func (o *rejectedMailParser) after(game *entity.Game, gameMail *entity.GameMail) {}
