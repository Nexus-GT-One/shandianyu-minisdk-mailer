package mail_parser

import (
	"regexp"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"strings"
)

// 提审app成功（海外账号）
type waitingForReviewMailParser struct{}

func init() {
	registerImplement(&waitingForReviewMailParser{})
}

func (o *waitingForReviewMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "App Store Connect")
}

func (o *waitingForReviewMailParser) checkTitle(title string) bool {
	return strings.Contains(title, `is now "Waiting for Review"`)
}

func (o *waitingForReviewMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "Waiting for Review")
}

func (o *waitingForReviewMailParser) parse(bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := service.GameService.GetByName(o.extractAppName(bodyText))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: findAuditingVersion(oneGame),
		Status:     "提审app成功 (海外账号)",
		Content:    bodyText,
	}
}

func (o *waitingForReviewMailParser) extractAppName(body string) string {
	re := regexp.MustCompile(`App Name:\s*(.*?)\s*App Version Number:`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}
