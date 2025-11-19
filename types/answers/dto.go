package answers

import "github.com/google/uuid"

type AnswerDto struct {
	UserID uuid.UUID `json:"userID" validate:"required,uuid"`
	Text   string    `json:"text" validate:"required"`
}
