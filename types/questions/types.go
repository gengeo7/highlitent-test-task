package questions

import (
	"time"

	"github.com/gengeo7/highlitent/types/answers"
)

type Question struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Text      string    `json:"text" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

type QuestionWithAnswers struct {
	Question Question         `json:"question"`
	Answers  []answers.Answer `json:"answers"`
}
