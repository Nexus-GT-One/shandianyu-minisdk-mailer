package mail_parser

import (
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"shandianyu-minisdk-mailer/util/arrayutil"
	"strings"
)

// 页面大小限制
type pageSizeRequirementsMailParser struct{}

func init() {
	registerImplement(&pageSizeRequirementsMailParser{})
}

func (o *pageSizeRequirementsMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "Google Play Support")
}

func (o *pageSizeRequirementsMailParser) checkTitle(title string) bool {
	return strings.Contains(title, "page size requirements")
}

func (o *pageSizeRequirementsMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "page size requirements")
}

func (o *pageSizeRequirementsMailParser) parse(from, to, bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := arrayutil.First(service.GameService.GetByDeveloperEmail(to))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: oneGame.PublishVersion,
		Status:     "页面大小限制",
		Content:    bodyText,
	}
}

func (o *pageSizeRequirementsMailParser) after(game *entity.Game, gameMail *entity.GameMail) {}
