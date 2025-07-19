package mail_parser

import (
	"shandianyu-minisdk-mailer/entity"
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
	return nil, &entity.GameMail{
		Status:  "欧盟认证成功",
		Content: bodyText,
	}
}
