package mail_parser

import (
	"regexp"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"strings"
)

// 正在审核
type appInReviewMailParser struct{}

func init() {
	registerImplement(&appInReviewMailParser{})
}

func (o *appInReviewMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "App Store Connect")
}

func (o *appInReviewMailParser) checkTitle(title string) bool {
	return strings.Contains(title, `is now "In Review"`)
}

func (o *appInReviewMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "In Review")
}

func (o *appInReviewMailParser) parse(bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := service.GameService.GetByName(o.extractAppName(bodyText))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: service.GameService.GetAuditingVersion(oneGame),
		BuildNum:   service.GameService.GetLastSubmitBuildNum(oneGame),
		Status:     "正在审核",
		Content:    bodyText,
	}
}

func (o *appInReviewMailParser) extractAppName(body string) string {
	re := regexp.MustCompile(`App Name:\s*(.*?)\s*App Version Number:`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func (o *appInReviewMailParser) after(game *entity.Game, gameMail *entity.GameMail) {}
