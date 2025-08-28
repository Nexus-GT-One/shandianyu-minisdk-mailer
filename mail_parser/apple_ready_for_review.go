package mail_parser

import (
	"regexp"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"strings"
)

// 准备审核
type readyForReviewMailParser struct{}

func init() {
	registerImplement(&readyForReviewMailParser{})
}

func (o *readyForReviewMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "App Store Connect")
}

func (o *readyForReviewMailParser) checkTitle(title string) bool {
	return strings.Contains(title, `is now "Ready For Review"`)
}

func (o *readyForReviewMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "Ready For Review")
}

func (o *readyForReviewMailParser) parse(from, to, bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := service.GameService.GetByName(o.extractAppName(bodyText))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: service.GameService.GetAuditingVersion(oneGame),
		BuildNum:   service.GameService.GetLastSubmitBuildNum(oneGame),
		Status:     "准备审核",
		Content:    bodyText,
	}
}

func (o *readyForReviewMailParser) extractAppName(body string) string {
	re := regexp.MustCompile(`App Name:\s*(.*?)\s*App Version Number:`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func (o *readyForReviewMailParser) after(game *entity.Game, gameMail *entity.GameMail) {}
