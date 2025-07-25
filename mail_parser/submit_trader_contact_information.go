package mail_parser

import (
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"shandianyu-minisdk-mailer/util/arrayutil"
	"strings"
)

// 提交欧盟认证
type submitTraderContactInformationMailParser struct{}

func init() {
	registerImplement(&submitTraderContactInformationMailParser{})
}

func (o *submitTraderContactInformationMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "App Store Connect")
}

func (o *submitTraderContactInformationMailParser) checkTitle(title string) bool {
	return strings.Contains(title, "We received your trader contact information")
}

func (o *submitTraderContactInformationMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "providing your trader contact information")
}

func (o *submitTraderContactInformationMailParser) parse(bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := arrayutil.Last(service.GameService.GetByDeveloperEmail(extractDeveloperEmail(bodyText)))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: findAuditingVersion(oneGame),
		Status:     "已提交欧盟认证",
		Content:    bodyText,
	}
}
