package mail_parser

import (
	"regexp"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"strings"
)

// 被拒
type noticedIssueMailParser struct{}

func init() {
	registerImplement(&noticedIssueMailParser{})
}

func (o *noticedIssueMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "App Store Connect")
}

func (o *noticedIssueMailParser) checkTitle(title string) bool {
	return strings.Contains(title, `We noticed an issue with your submission`)
}

func (o *noticedIssueMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "noticed an issue")
}

func (o *noticedIssueMailParser) parse(bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := service.GameService.GetByName(o.extractAppName(bodyText))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: service.GameService.GetAuditingVersion(oneGame),
		Status:     "被拒",
		Content:    bodyText,
	}
}

func (o *noticedIssueMailParser) extractAppName(body string) string {
	re := regexp.MustCompile(`App Name:\s*(.*?)\s*Submission ID`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func (o *noticedIssueMailParser) after(game *entity.Game, gameMail *entity.GameMail) {
	service.GameService.RecordRejected(game.BundleId, service.GameService.GetAuditingVersion(game))
}
