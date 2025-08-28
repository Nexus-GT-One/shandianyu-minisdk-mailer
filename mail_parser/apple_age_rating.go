package mail_parser

import (
	"shandianyu-minisdk-mailer/entity"
	"strings"
)

// 年龄分级提示
type ageRatingMailParser struct{}

func init() {
	registerImplement(&ageRatingMailParser{})
}

func (o *ageRatingMailParser) checkFrom(from string) bool {
	return strings.Contains(from, "Apple Developer")
}

func (o *ageRatingMailParser) checkTitle(title string) bool {
	return strings.Contains(title, "Answer the updated age rating questions")
}

func (o *ageRatingMailParser) checkKeyword(bodyText string) bool {
	return strings.Contains(bodyText, "updated the age rating system for apps")
}

func (o *ageRatingMailParser) parse(from, to, bodyText string) (*entity.Game, *entity.GameMail) {
	return &entity.Game{}, &entity.GameMail{
		Status:  "年龄分级提示",
		Content: bodyText,
	}
}

func (o *ageRatingMailParser) after(game *entity.Game, gameMail *entity.GameMail) {}
