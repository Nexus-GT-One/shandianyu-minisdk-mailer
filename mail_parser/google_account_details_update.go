package mail_parser

import (
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"shandianyu-minisdk-mailer/util/arrayutil"
	"strings"
)

// 账号信息更新
type accountDetailsUpdateMailParser struct{}

func init() {
	registerImplement(&accountDetailsUpdateMailParser{})
}

func (o *accountDetailsUpdateMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "Google Play Support")
}

func (o *accountDetailsUpdateMailParser) checkTitle(title string) bool {
	return strings.Contains(title, "Your account details have been successfully updated")
}

func (o *accountDetailsUpdateMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "Your account details have been successfully updated")
}

func (o *accountDetailsUpdateMailParser) parse(from, to, bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := arrayutil.First(service.GameService.GetByDeveloperEmail(to))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: oneGame.PublishVersion,
		Status:     "账号信息更新",
		Content:    bodyText,
	}
}

func (o *accountDetailsUpdateMailParser) after(game *entity.Game, gameMail *entity.GameMail) {}
