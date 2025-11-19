package questions

import (
	"context"
	"errors"
	"testing"

	"github.com/gengeo7/highlitent/apierror"
	"github.com/gengeo7/highlitent/storage"
	"github.com/gengeo7/highlitent/types/answers"
	"github.com/gengeo7/highlitent/types/questions"
	"github.com/gengeo7/highlitent/utils"
	"github.com/google/go-cmp/cmp"
)

type mockQuestionsGetter struct {
	ReturnedValue []questions.Question
	ReturnedError error
}

func (m *mockQuestionsGetter) QuestionsGet(ctx context.Context) ([]questions.Question, error) {
	return m.ReturnedValue, m.ReturnedError
}

func TestGetAllQuestions(t *testing.T) {
	tests := []struct {
		name            string
		questionsGetter QuestionsGetter
		want            []questions.Question
		wantErr         *apierror.ApiError
	}{
		{
			name: "ok",
			questionsGetter: &mockQuestionsGetter{
				ReturnedValue: []questions.Question{
					{
						Text: "test1",
					},
					{
						Text: "test2",
					},
					{
						Text: "test3",
					},
				},
				ReturnedError: nil,
			},
			want: []questions.Question{
				{
					Text: "test1",
				},
				{
					Text: "test2",
				},
				{

					Text: "test3",
				},
			},
			wantErr: nil,
		},
		{
			name: "empty list",
			questionsGetter: &mockQuestionsGetter{
				ReturnedValue: []questions.Question{},
				ReturnedError: nil,
			},
			want:    []questions.Question{},
			wantErr: nil,
		},
		{
			name: "not found is unexpected error",
			questionsGetter: &mockQuestionsGetter{
				ReturnedValue: nil,
				ReturnedError: storage.ErrDbNotFound,
			},
			want:    nil,
			wantErr: utils.UnhandledError(nil),
		},
		{
			name: "deadline exceeded",
			questionsGetter: &mockQuestionsGetter{
				ReturnedValue: nil,
				ReturnedError: context.DeadlineExceeded,
			},
			want:    nil,
			wantErr: utils.DeadlineDbError(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := GetAllQuestions(context.Background(), tt.questionsGetter)
			if gotErr != nil {
				if tt.wantErr == nil {
					t.Fatalf("GetAllQuestions() failed: %v", gotErr)
				}
				var gotApiError *apierror.ApiError
				if errors.As(gotErr, &gotApiError) {
					if gotApiError.Msg != tt.wantErr.Msg || gotApiError.StatusCode != tt.wantErr.StatusCode {
						t.Fatalf("GetAllQuestions(): %v, want: %v", gotErr, tt.wantErr)
					}
				} else {
					t.Fatalf("GetAllQuestions() expected error of type ApiError: %v", gotErr)
				}
				return
			}

			if tt.wantErr != nil {
				t.Fatal("GetAllQuestions() succeeded unexpectedly")
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetAllQuestions() mismatch:\n %s", diff)
			}
		})
	}
}

type mockQuestionCreater struct {
	ReturnedValue *questions.Question
	ReturnedError error
}

func (m *mockQuestionCreater) QuestionCreate(ctx context.Context, dto *questions.QuestionDto) (*questions.Question, error) {
	return m.ReturnedValue, m.ReturnedError
}

func TestCreateQuestion(t *testing.T) {
	tests := []struct {
		name            string
		questionCreater QuestionCreater
		dto             *questions.QuestionDto
		want            *questions.Question
		wantErr         *apierror.ApiError
	}{
		{
			name: "created",
			questionCreater: &mockQuestionCreater{
				ReturnedValue: &questions.Question{
					Text: "test",
				},
				ReturnedError: nil,
			},
			dto: &questions.QuestionDto{
				Text: "test",
			},
			want: &questions.Question{
				Text: "test",
			},
			wantErr: nil,
		},
		{
			name: "deadline exceeded",
			questionCreater: &mockQuestionCreater{
				ReturnedValue: nil,
				ReturnedError: context.DeadlineExceeded,
			},
			dto: &questions.QuestionDto{
				Text: "test",
			},
			want:    nil,
			wantErr: utils.DeadlineDbError(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := CreateQuestion(context.Background(), tt.questionCreater, tt.dto)
			if gotErr != nil {
				if tt.wantErr == nil {
					t.Fatalf("CreateQuestion() failed: %v", gotErr)
				}
				var gotApiError *apierror.ApiError
				if errors.As(gotErr, &gotApiError) {
					if gotApiError.Msg != tt.wantErr.Msg || gotApiError.StatusCode != tt.wantErr.StatusCode {
						t.Fatalf("CreateQuestion(): %v, want: %v", gotErr, tt.wantErr)
					}
				} else {
					t.Fatalf("CreateQuestion() expected error of type ApiError: %v", gotErr)
				}
				return
			}

			if tt.wantErr != nil {
				t.Fatal("CreateQuestion() succeeded unexpectedly")
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("CreateQuestion() mismatch:\n %s", diff)
			}
		})
	}
}

type mockQuestionGetter struct {
	ReturnedValue *questions.QuestionWithAnswers
	ReturnedError error
}

func (m *mockQuestionGetter) QuestionGet(ctx context.Context, id int) (*questions.QuestionWithAnswers, error) {
	return m.ReturnedValue, m.ReturnedError
}

func TestGetQuestionWithAnswers(t *testing.T) {
	tests := []struct {
		name           string
		questionGetter QuestionGetter
		id             int
		want           *questions.QuestionWithAnswers
		wantErr        *apierror.ApiError
	}{
		{
			name: "ok",
			questionGetter: &mockQuestionGetter{
				ReturnedValue: &questions.QuestionWithAnswers{
					Question: questions.Question{
						Text: "test",
					},
					Answers: []answers.Answer{
						{
							Text: "test1",
						},
					},
				},
			},
			id: 1,
			want: &questions.QuestionWithAnswers{
				Question: questions.Question{
					Text: "test",
				},
				Answers: []answers.Answer{
					{
						Text: "test1",
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "not found",
			questionGetter: &mockQuestionGetter{
				ReturnedValue: nil,
				ReturnedError: storage.ErrDbNotFound,
			},
			id:      1,
			want:    nil,
			wantErr: utils.QuestionNotFound(nil),
		},
		{
			name: "deadline exceeded",
			questionGetter: &mockQuestionGetter{
				ReturnedValue: nil,
				ReturnedError: context.DeadlineExceeded,
			},
			id:      1,
			want:    nil,
			wantErr: utils.DeadlineDbError(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := GetQuestionWithAnswers(context.Background(), tt.questionGetter, tt.id)
			if gotErr != nil {
				if tt.wantErr == nil {
					t.Fatalf("CreateQuestion() failed: %v", gotErr)
				}
				var gotApiError *apierror.ApiError
				if errors.As(gotErr, &gotApiError) {
					if gotApiError.Msg != tt.wantErr.Msg || gotApiError.StatusCode != tt.wantErr.StatusCode {
						t.Fatalf("CreateQuestion(): %v, want: %v", gotErr, tt.wantErr)
					}
				} else {
					t.Fatalf("CreateQuestion() expected error of type ApiError: %v", gotErr)
				}
				return
			}

			if tt.wantErr != nil {
				t.Fatal("CreateQuestion() succeeded unexpectedly")
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("CreateQuestion() mismatch:\n %s", diff)
			}
		})
	}
}

type mockQuestionDeleter struct {
	ReturnedError error
}

func (m *mockQuestionDeleter) QuestionDelete(ctx context.Context, id int) error {
	return m.ReturnedError
}

func TestDeleteQuestion(t *testing.T) {
	tests := []struct {
		name            string
		questionDeleter QuestionDeleter
		id              int
		wantErr         *apierror.ApiError
	}{
		{
			name: "ok",
			questionDeleter: &mockQuestionDeleter{
				ReturnedError: nil,
			},
			id:      1,
			wantErr: nil,
		},
		{
			name: "not found",
			questionDeleter: &mockQuestionDeleter{
				ReturnedError: storage.ErrDbNotFound,
			},
			id:      1,
			wantErr: utils.QuestionNotFound(nil),
		},
		{
			name: "deadline exceeded",
			questionDeleter: &mockQuestionDeleter{
				ReturnedError: context.DeadlineExceeded,
			},
			id:      1,
			wantErr: utils.DeadlineDbError(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := DeleteQuestion(context.Background(), tt.questionDeleter, tt.id)
			if gotErr != nil {
				if tt.wantErr == nil {
					t.Fatalf("CreateQuestion() failed: %v", gotErr)
				}
				var gotApiError *apierror.ApiError
				if errors.As(gotErr, &gotApiError) {
					if gotApiError.Msg != tt.wantErr.Msg || gotApiError.StatusCode != tt.wantErr.StatusCode {
						t.Fatalf("CreateQuestion(): %v, want: %v", gotErr, tt.wantErr)
					}
				} else {
					t.Fatalf("CreateQuestion() expected error of type ApiError: %v", gotErr)
				}
				return
			}

			if tt.wantErr != nil {
				t.Fatal("CreateQuestion() succeeded unexpectedly")
			}
		})
	}
}
