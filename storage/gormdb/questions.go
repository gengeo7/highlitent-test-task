package gormdb

import (
	"context"
	"errors"

	"github.com/gengeo7/highlitent/storage"
	"github.com/gengeo7/highlitent/types/answers"
	"github.com/gengeo7/highlitent/types/questions"
	"gorm.io/gorm"
)

func (d *Db) QuestionsGet(ctx context.Context) ([]questions.Question, error) {
	var qs []questions.Question
	err := d.Db.WithContext(ctx).
		Find(&qs).Error
	if err != nil {
		return nil, err
	}
	return qs, nil
}

func (d *Db) QuestionCreate(ctx context.Context, dto *questions.QuestionDto) (*questions.Question, error) {
	q := &questions.Question{
		Text: dto.Text,
	}

	err := d.Db.WithContext(ctx).Create(q).Error
	if err != nil {
		return nil, err
	}
	return q, nil
}

func (d *Db) QuestionGet(ctx context.Context, id int) (*questions.QuestionWithAnswers, error) {
	var result questions.QuestionWithAnswers

	err := d.Db.WithContext(ctx).
		Where("questions.id = ?", id).
		First(&result.Question).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, storage.ErrDbNotFound
		}
		return nil, err
	}

	var answers []answers.Answer
	err = d.Db.WithContext(ctx).
		Where("question_id = ?", id).
		Order("created_at asc").
		Find(&answers).Error

	if err != nil {
		return nil, err
	}

	result.Answers = answers
	return &result, nil
}

func (d *Db) QuestionDelete(ctx context.Context, id int) error {
	res := d.Db.WithContext(ctx).Unscoped().
		Delete(&questions.Question{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return storage.ErrDbNotFound
	}

	return nil
}
