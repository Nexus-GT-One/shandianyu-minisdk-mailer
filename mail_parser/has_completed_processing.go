package mail_parser

import (
	"regexp"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/service"
	"strings"
)

// 上传App成功（可以手动选择Build了）
type hasCompletedProcessingMailParser struct{}

func init() {
	registerImplement(&hasCompletedProcessingMailParser{})
}

func (o *hasCompletedProcessingMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "App Store Connect")
}

func (o *hasCompletedProcessingMailParser) checkTitle(title string) bool {
	return strings.Contains(title, "has completed processing")
}

func (o *hasCompletedProcessingMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "has completed processing")
}

func (o *hasCompletedProcessingMailParser) parse(bodyText string) (*entity.Game, *entity.GameMail) {
	oneGame := service.GameService.GetByName(o.extractAppName(bodyText))
	if oneGame == nil {
		return nil, nil
	}
	return oneGame, &entity.GameMail{
		Symbol:     oneGame.Symbol,
		AppVersion: findAuditingVersion(oneGame),
		Status:     "上传App成功",
		Content:    bodyText,
	}
}

func (o *hasCompletedProcessingMailParser) extractAppName(body string) string {
	re := regexp.MustCompile(`App Name:\s*(.*?)\s*Build Number:`)
	matches := re.FindStringSubmatch(body)
	if len(matches) > 1 {
		return strings.TrimSpace(matches[1])
	}
	return ""
}
