package mail_parser

import (
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"shandianyu-minisdk-mailer/util/arrayutil"
	"strings"
)

// API级别更新
type apiLevelRequirementsMailParser struct{}

func init() {
	registerImplement(&apiLevelRequirementsMailParser{})
}

func (o *apiLevelRequirementsMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "Google Play Support")
}

func (o *apiLevelRequirementsMailParser) checkTitle(title string) bool {
	return strings.Contains(title, "API level requirements")
}

func (o *apiLevelRequirementsMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "API level requirements")
}

func (o *apiLevelRequirementsMailParser) parse(from, to, bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := arrayutil.First(service.GameService.GetByDeveloperEmail(to))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: oneGame.PublishVersion,
		Status:     "API级别更新",
		Content:    bodyText,
	}
}

func (o *apiLevelRequirementsMailParser) after(game *entity.Game, gameMail *entity.GameMail) {}
