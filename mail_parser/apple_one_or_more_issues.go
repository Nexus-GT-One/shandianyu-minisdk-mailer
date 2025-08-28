package mail_parser

import (
	"regexp"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"strings"
)

// 有问题待纠正
type OneOrMoreIssuesMailParser struct{}

func init() {
	registerImplement(&OneOrMoreIssuesMailParser{})
}

func (o *OneOrMoreIssuesMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "App Store Connect")
}

func (o *OneOrMoreIssuesMailParser) checkTitle(title string) bool {
	return strings.Contains(title, "one or more issues")
}

func (o *OneOrMoreIssuesMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "one or more issues")
}

func (o *OneOrMoreIssuesMailParser) parse(from, to, bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := service.GameService.GetByName(o.extractAppName(bodyText))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: service.GameService.GetAuditingVersion(oneGame),
		BuildNum:   service.GameService.GetLastSubmitBuildNum(oneGame),
		Status:     "有问题待纠正",
		Content:    bodyText,
	}
}

func (o *OneOrMoreIssuesMailParser) extractAppName(body string) string {
	re := regexp.MustCompile(`(?m)^(.+)\s*\nApp Apple ID\s+\d+\s*\nVersion\s+([^\s]+)`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func (o *OneOrMoreIssuesMailParser) after(game *entity.Game, gameMail *entity.GameMail) {}
