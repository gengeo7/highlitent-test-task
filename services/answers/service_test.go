package answers

import (
	"context"
	"errors"
	"testing"

	"github.com/gengeo7/highlitent/apierror"
	"github.com/gengeo7/highlitent/storage"
	"github.com/gengeo7/highlitent/types/answers"
	"github.com/gengeo7/highlitent/utils"
	"github.com/google/go-cmp/cmp"
)

type mockAnswerGetter struct {
	ReturnedValue *answers.Answer
	ReturnedError error
}

func (m *mockAnswerGetter) AnswerGet(ctx context.Context, id int) (*answers.Answer, error) {
	return m.ReturnedValue, m.ReturnedError
}

func TestGetAnswer(t *testing.T) {
	tests := []struct {
		name         string
		answerGetter AnswerGetter
		id           int
		want         *answers.Answer
		wantErr      *apierror.ApiError
	}{
		{
			name: "ok",
			id:   1,
			answerGetter: &mockAnswerGetter{
				ReturnedValue: &answers.Answer{
					Text: "test",
				},
				ReturnedError: nil,
			},
			want: &answers.Answer{
				Text: "test",
			},
			wantErr: nil,
		},
		{
			name: "not found",
			id:   1,
			answerGetter: &mockAnswerGetter{
				ReturnedValue: nil,
				ReturnedError: storage.ErrDbNotFound,
			},
			want:    nil,
			wantErr: utils.AnswerNotFound(nil),
		},
		{
			name: "deadline exceeded",
			id:   1,
			answerGetter: &mockAnswerGetter{
				ReturnedValue: nil,
				ReturnedError: context.DeadlineExceeded,
			},
			want:    nil,
			wantErr: utils.DeadlineDbError(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := GetAnswer(context.Background(), tt.answerGetter, tt.id)
			if gotErr != nil {
				if tt.wantErr == nil {
					t.Fatalf("GetAnswer() failed: %v", gotErr)
				}
				var gotApiError *apierror.ApiError
				if errors.As(gotErr, &gotApiError) {
					if gotApiError.Msg != tt.wantErr.Msg || gotApiError.StatusCode != tt.wantErr.StatusCode {
						t.Fatalf("GetAnswer(): %v, want: %v", gotErr, tt.wantErr)
					}
				} else {
					t.Fatalf("GetAnswer() expected error of type ApiError: %v", gotErr)
				}
				return
			}

			if tt.wantErr != nil {
				t.Fatal("GetAnswer() succeeded unexpectedly")
			}

			if diff := cmp.Diff(*tt.want, *got); diff != "" {
				t.Errorf("GetAnswer() mismatch:\n %s", diff)
			}
		})
	}
}

type mockAnswerCreater struct {
	ReturnedValue *answers.Answer
	ReturnedError error
}

func (m *mockAnswerCreater) AnswerCreate(ctx context.Context, dto *answers.AnswerDto, questionID int) (*answers.Answer, error) {
	return m.ReturnedValue, m.ReturnedError
}

func TestCreateAnswer(t *testing.T) {
	tests := []struct {
		name          string
		answerCreater AnswerCreater
		dto           *answers.AnswerDto
		questionID    int
		want          *answers.Answer
		wantErr       *apierror.ApiError
	}{
		{
			name: "created",
			answerCreater: &mockAnswerCreater{
				ReturnedValue: &answers.Answer{
					Text: "test",
				},
				ReturnedError: nil,
			},
			dto:        &answers.AnswerDto{},
			questionID: 1,
			want: &answers.Answer{
				Text: "test",
			},
			wantErr: nil,
		},
		{
			name: "question not found",
			answerCreater: &mockAnswerCreater{
				ReturnedValue: nil,
				ReturnedError: storage.ErrDbNotFound,
			},
			questionID: 1,
			want:       nil,
			wantErr:    utils.QuestionNotFound(nil),
		},
		{
			name: "deadline exceeded",
			answerCreater: &mockAnswerCreater{
				ReturnedValue: nil,
				ReturnedError: context.DeadlineExceeded,
			},
			questionID: 1,
			want:       nil,
			wantErr:    utils.DeadlineDbError(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := CreateAnswer(context.Background(), tt.answerCreater, tt.dto, tt.questionID)
			if gotErr != nil {
				if tt.wantErr == nil {
					t.Fatalf("CreateAnswer() failed: %v", gotErr)
				}
				var gotApiError *apierror.ApiError
				if errors.As(gotErr, &gotApiError) {
					if gotApiError.Msg != tt.wantErr.Msg || gotApiError.StatusCode != tt.wantErr.StatusCode {
						t.Fatalf("CreateAnswer(): %v, want: %v", gotErr, tt.wantErr)
					}
				} else {
					t.Fatalf("CreateAnswer() expected error of type ApiError: %v", gotErr)
				}
				return
			}

			if tt.wantErr != nil {
				t.Fatal("CreateAnswer() succeeded unexpectedly")
			}

			if diff := cmp.Diff(*tt.want, *got); diff != "" {
				t.Errorf("CreateAnswer() mismatch:\n %s", diff)
			}
		})
	}
}

type mockAnswerDeleter struct {
	ReturnedError error
}

func (m *mockAnswerDeleter) AnswerDelete(ctx context.Context, id int) error {
	return m.ReturnedError
}

func TestDeleteAnswer(t *testing.T) {
	tests := []struct {
		name          string
		answerDeleter AnswerDeleter
		id            int
		wantErr       *apierror.ApiError
	}{
		{
			name: "ok",
			answerDeleter: &mockAnswerDeleter{
				ReturnedError: nil,
			},
			id:      1,
			wantErr: nil,
		},
		{
			name: "not found",
			answerDeleter: &mockAnswerDeleter{
				ReturnedError: storage.ErrDbNotFound,
			},
			id:      1,
			wantErr: utils.AnswerNotFound(nil),
		},
		{
			name: "deadline exceeded",
			answerDeleter: &mockAnswerDeleter{
				ReturnedError: context.DeadlineExceeded,
			},
			id:      1,
			wantErr: utils.DeadlineDbError(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := DeleteAnswer(context.Background(), tt.answerDeleter, tt.id)
			if gotErr != nil {
				if tt.wantErr == nil {
					t.Fatalf("DeleteAnswer() failed: %v", gotErr)
				}
				var gotApiError *apierror.ApiError
				if errors.As(gotErr, &gotApiError) {
					if gotApiError.Msg != tt.wantErr.Msg || gotApiError.StatusCode != tt.wantErr.StatusCode {
						t.Fatalf("DeleteAnswer(): %v, want: %v", gotErr, tt.wantErr)
					}
				} else {
					t.Fatalf("DeleteAnswer() expected error of type ApiError: %v", gotErr)
				}
				return
			}

			if tt.wantErr != nil {
				t.Fatal("DeleteAnswer() succeeded unexpectedly")
			}
		})
	}
}
