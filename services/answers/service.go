package answers

import (
	"context"

	"github.com/gengeo7/highlitent/storage"
	"github.com/gengeo7/highlitent/types/answers"
	"github.com/gengeo7/highlitent/utils"
)

type AnswerGetter interface {
	AnswerGet(ctx context.Context, id int) (*answers.Answer, error)
}

type AnswerCreater interface {
	AnswerCreate(ctx context.Context, dto *answers.AnswerDto, questionID int) (*answers.Answer, error)
}

type AnswerDeleter interface {
	AnswerDelete(ctx context.Context, id int) error
}

func GetAnswer(ctx context.Context, answerGetter AnswerGetter, id int) (*answers.Answer, error) {
	answer, err := answerGetter.AnswerGet(ctx, id)
	if err != nil {
		return nil, utils.TestDbErr(err, &utils.ErrDbCase{Func: storage.IsErrNotFound, Creator: utils.AnswerNotFound, CheckErr: false})
	}
	return answer, nil
}

func CreateAnswer(ctx context.Context, answerCreater AnswerCreater, dto *answers.AnswerDto, questionID int) (*answers.Answer, error) {
	answer, err := answerCreater.AnswerCreate(ctx, dto, questionID)
	if err != nil {
		return nil, utils.TestDbErr(err, &utils.ErrDbCase{Func: storage.IsErrNotFound, Creator: utils.QuestionNotFound, CheckErr: false})
	}
	return answer, nil
}

func DeleteAnswer(ctx context.Context, answerDeleter AnswerDeleter, id int) error {
	err := answerDeleter.AnswerDelete(ctx, id)
	if err != nil {
		return utils.TestDbErr(err, &utils.ErrDbCase{Func: storage.IsErrNotFound, Creator: utils.AnswerNotFound, CheckErr: false})
	}
	return nil
}
