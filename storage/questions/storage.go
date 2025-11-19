package questions

import (
	"context"

	"github.com/gengeo7/highlitent/types/questions"
)

type Storage interface {
	QuestionsGet(ctx context.Context) ([]questions.Question, error)
	QuestionCreate(ctx context.Context, dto *questions.QuestionDto) (*questions.Question, error)
	QuestionGet(ctx context.Context, id int) (*questions.QuestionWithAnswers, error)
	QuestionDelete(ctx context.Context, id int) error
}
