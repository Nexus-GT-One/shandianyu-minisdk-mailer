package mail_parser

import (
	"regexp"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"strings"
)

// 过审
type reviewApprovedMailParser struct{}

func init() {
	registerImplement(&reviewApprovedMailParser{})
}

func (o *reviewApprovedMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "App Store Connect")
}

func (o *reviewApprovedMailParser) checkTitle(title string) bool {
	return strings.Contains(title, "Review of your submission is complete")
}

func (o *reviewApprovedMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "eligible for distribution")
}

func (o *reviewApprovedMailParser) parse(from, to, bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := service.GameService.GetByName(o.extractAppName(bodyText))
	if oneGame == nil {
		oneGame = service.GameService.GetByAppId(o.extractAppId(bodyText))
		if oneGame == nil {
			return nil, nil
		}
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: service.GameService.GetAuditingVersion(oneGame),
		BuildNum:   service.GameService.GetLastSubmitBuildNum(oneGame),
		Status:     "过审",
		Content:    bodyText,
	}
}

func (o *reviewApprovedMailParser) extractAppName(body string) string {
	re := regexp.MustCompile(`(?m)^App\sName:\s([A-Za-z\s]+)$`)
	match := re.FindStringSubmatch(body)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func (o *reviewApprovedMailParser) extractAppId(s string) string {
	// 忽略大小写，匹配 /id 后的一串数字
	re := regexp.MustCompile(`(?i)/id(\d+)\b`)
	m := re.FindStringSubmatch(s)
	if len(m) >= 2 {
		return m[1]
	}
	return ""
}

func (o *reviewApprovedMailParser) after(game *entity.Game, gameMail *entity.GameMail) {
	service.ApplicationService.CheckApplicationNewVersion(game, gameMail)
}
