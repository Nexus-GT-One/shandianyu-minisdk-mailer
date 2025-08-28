package mail_parser

import (
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"shandianyu-minisdk-mailer/util/arrayutil"
	"strings"
)

// 封号
type accountTerminatedMailParser struct{}

func init() {
	registerImplement(&accountTerminatedMailParser{})
}

func (o *accountTerminatedMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "Google Play Support")
}

func (o *accountTerminatedMailParser) checkTitle(title string) bool {
	return strings.Contains(title, "Notification from Google Play")
}

func (o *accountTerminatedMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "Account Terminated")
}

func (o *accountTerminatedMailParser) parse(from, to, bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := arrayutil.First(service.GameService.GetByDeveloperEmail(to))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: oneGame.PublishVersion,
		Status:     "封号",
		Content:    bodyText,
	}
}

func (o *accountTerminatedMailParser) after(game *entity.Game, gameMail *entity.GameMail) {}
