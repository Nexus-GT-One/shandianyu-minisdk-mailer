package mail_parser

import (
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"shandianyu-minisdk-mailer/util/arrayutil"
	"strings"
)

// 政策更新
type policyUpdateMailParser struct{}

func init() {
	registerImplement(&policyUpdateMailParser{})
}

func (o *policyUpdateMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "Google Play Support")
}

func (o *policyUpdateMailParser) checkTitle(title string) bool {
	return strings.Contains(title, "Policy Update")
}

func (o *policyUpdateMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "Policy Update")
}

func (o *policyUpdateMailParser) parse(from, to, bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := arrayutil.First(service.GameService.GetByDeveloperEmail(to))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: oneGame.PublishVersion,
		Status:     "政策更新",
		Content:    bodyText,
	}
}

func (o *policyUpdateMailParser) after(game *entity.Game, gameMail *entity.GameMail) {}
