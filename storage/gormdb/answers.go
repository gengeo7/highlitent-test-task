package gormdb

import (
	"context"
	"errors"
	"strings"

	"github.com/gengeo7/highlitent/storage"
	"github.com/gengeo7/highlitent/types/answers"
	"gorm.io/gorm"
)

func (d *Db) AnswerGet(ctx context.Context, id int) (*answers.Answer, error) {
	var a answers.Answer
	err := d.Db.WithContext(ctx).
		Where("id = ?", id).
		First(&a).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, storage.ErrDbNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (d *Db) AnswerCreate(ctx context.Context, dto *answers.AnswerDto, questionID int) (*answers.Answer, error) {
	a := &answers.Answer{
		QuestionID: questionID,
		UserID:     dto.UserID,
		Text:       dto.Text,
	}

	err := d.Db.WithContext(ctx).Create(a).Error
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "foreign key") {
			return nil, storage.ErrDbNotFound
		}
		return nil, err
	}
	return a, nil
}

func (d *Db) AnswerDelete(ctx context.Context, id int) error {
	res := d.Db.WithContext(ctx).Unscoped().
		Delete(&answers.Answer{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return storage.ErrDbNotFound
	}

	return nil
}
