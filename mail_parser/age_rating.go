package mail_parser

import (
	"shandianyu-minisdk-mailer/entity"
	"strings"
)

// 年龄分级提示
type AgeRatingMailParser struct{}

func init() {
	registerImplement(&AgeRatingMailParser{})
}

func (o *AgeRatingMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "Apple Developer")
}

func (o *AgeRatingMailParser) checkTitle(title string) bool {
	return strings.Contains(title, "Answer the updated age rating questions")
}

func (o *AgeRatingMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "updated the age rating system for apps")
}

func (o *AgeRatingMailParser) parse(bodyText string) (*entity.Game, *entity.GameMail) {
	return &entity.Game{}, &entity.GameMail{
		Status:  "年龄分级提示",
		Content: bodyText,
	}
}

func (o *AgeRatingMailParser) after(game *entity.Game, gameMail *entity.GameMail) {}
