package mail_parser

import (
	"regexp"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"strings"
)

// 欧盟认证成功
type verifiedEuropeanUnionSuccessMailParser struct{}

func init() {
	registerImplement(&verifiedEuropeanUnionSuccessMailParser{})
}

func (o *verifiedEuropeanUnionSuccessMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "App Store Connect")
}

func (o *verifiedEuropeanUnionSuccessMailParser) checkTitle(title string) bool {
	return strings.Contains(title, "Your trader contact information was verified")
}

func (o *verifiedEuropeanUnionSuccessMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "verified your trader contact") && strings.Contains(bodyText, "European Union")
}

func (o *verifiedEuropeanUnionSuccessMailParser) parse(bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := service.GameService.GetByName(o.extractAppName(bodyText))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: findAuditingVersion(oneGame),
		Status:     "欧盟认证成功",
		Content:    bodyText,
	}
}

func (o *verifiedEuropeanUnionSuccessMailParser) extractAppName(body string) string {
	re := regexp.MustCompile(`(?m)^We've received your app for review\.\s*\n(.+?)\n`)
	match := re.FindStringSubmatch(body)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}
