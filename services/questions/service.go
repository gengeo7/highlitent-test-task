package questions

import (
	"context"

	"github.com/gengeo7/highlitent/storage"
	"github.com/gengeo7/highlitent/types/questions"
	"github.com/gengeo7/highlitent/utils"
)

type QuestionsGetter interface {
	QuestionsGet(ctx context.Context) ([]questions.Question, error)
}

type QuestionCreater interface {
	QuestionCreate(ctx context.Context, dto *questions.QuestionDto) (*questions.Question, error)
}

type QuestionGetter interface {
	QuestionGet(ctx context.Context, id int) (*questions.QuestionWithAnswers, error)
}

type QuestionDeleter interface {
	QuestionDelete(ctx context.Context, id int) error
}

func GetAllQuestions(ctx context.Context, questionsGetter QuestionsGetter) ([]questions.Question, error) {
	questions, err := questionsGetter.QuestionsGet(ctx)
	if err != nil {
		return nil, utils.TestDbErr(err)
	}
	return questions, nil
}

func CreateQuestion(ctx context.Context, questionCreater QuestionCreater, dto *questions.QuestionDto) (*questions.Question, error) {
	if dto == nil {
		return nil, utils.EmptyDto(nil)
	}
	question, err := questionCreater.QuestionCreate(ctx, dto)
	if err != nil {
		return nil, utils.TestDbErr(err)
	}
	return question, nil
}

func GetQuestionWithAnswers(ctx context.Context, questionGetter QuestionGetter, id int) (*questions.QuestionWithAnswers, error) {
	questionWithAnswer, err := questionGetter.QuestionGet(ctx, id)
	if err != nil {
		return nil, utils.TestDbErr(err, &utils.ErrDbCase{Func: storage.IsErrNotFound, Creator: utils.QuestionNotFound, CheckErr: false})
	}
	return questionWithAnswer, nil
}

func DeleteQuestion(ctx context.Context, questionDeleter QuestionDeleter, id int) error {
	err := questionDeleter.QuestionDelete(ctx, id)
	if err != nil {
		return utils.TestDbErr(err, &utils.ErrDbCase{Func: storage.IsErrNotFound, Creator: utils.QuestionNotFound, CheckErr: false})
	}
	return nil
}
