package mail_parser

import (
	"regexp"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"strings"
)

// 疑似封号
type seemsBanAccountMailParser struct{}

func init() {
	registerImplement(&seemsBanAccountMailParser{})
}

func (o *seemsBanAccountMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "App Store Connect")
}

func (o *seemsBanAccountMailParser) checkTitle(title string) bool {
	return strings.Contains(title, `You have a message from App Review`)
}

func (o *seemsBanAccountMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "app Apple ID: ")
}

func (o *seemsBanAccountMailParser) parse(bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := service.GameService.GetByName(o.extractAppName(bodyText))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: oneGame.PublishVersion,
		Status:     "疑似封号",
		Content:    bodyText,
	}
}

func (o *seemsBanAccountMailParser) extractAppName(body string) string {
	re := regexp.MustCompile(`app,\s*(.+?),\s*app Apple ID`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func (o *seemsBanAccountMailParser) after(game *entity.Game, gameMail *entity.GameMail) {}
