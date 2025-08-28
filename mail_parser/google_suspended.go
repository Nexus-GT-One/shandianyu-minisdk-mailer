package mail_parser

import (
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"strings"
)

// 应用下架
type suspendedMailParser struct{}

func init() {
	registerImplement(&suspendedMailParser{})
}

func (o *suspendedMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "Google Play Support")
}

func (o *suspendedMailParser) checkTitle(title string) bool {
	return strings.Contains(title, "Your app is not compliant")
}

func (o *suspendedMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "Suspended")
}

func (o *suspendedMailParser) parse(from, to, bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := service.GameService.GetByName(extractBundleId(bodyText))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: oneGame.PublishVersion,
		Status:     "应用下架",
		Content:    bodyText,
	}
}

func (o *suspendedMailParser) after(game *entity.Game, gameMail *entity.GameMail) {}
