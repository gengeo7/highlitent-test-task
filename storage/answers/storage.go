package answers

import (
	"context"

	"github.com/gengeo7/highlitent/types/answers"
)

type Storage interface {
	AnswerGet(ctx context.Context, id int) (*answers.Answer, error)
	AnswerCreate(ctx context.Context, dto *answers.AnswerDto, questionID int) (*answers.Answer, error)
	AnswerDelete(ctx context.Context, id int) error
}
