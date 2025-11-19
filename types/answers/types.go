package answers

import (
	"time"

	"github.com/google/uuid"
)

type Answer struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	QuestionID int       `json:"questionID" gorm:"index;not null"`
	UserID     uuid.UUID `json:"userID" gorm:"uuid;not null"`
	Text       string    `json:"text" gorm:"type:text;not null"`
	CreatedAt  time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
