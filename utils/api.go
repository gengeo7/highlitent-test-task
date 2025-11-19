package utils

import (
	"encoding/json"
	"net/http"

	"github.com/gengeo7/highlitent/apierror"
	"github.com/gengeo7/highlitent/storage"
	"github.com/gengeo7/highlitent/types/common"
)

type Response struct {
	Data   any
	Status int
}

func SendResponse(response *Response, err error, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		apierror.SendError(w, r, err)
		return
	}

	if response == nil {
		response = &Response{
			Data: common.MessageDto{
				Message: "ok",
			},
			Status: http.StatusOK,
		}
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(response.Status)
	if err := json.NewEncoder(w).Encode(response.Data); err != nil {
		apierror.SendError(w, r, err)
	}
}

type ErrorCreator = func(err error) *apierror.ApiError

type ErrDbCase struct {
	Func     func(error) bool
	Creator  ErrorCreator
	CheckErr bool
}

func DeadlineDbError(err error) *apierror.ApiError {
	return apierror.NewApiError(http.StatusRequestTimeout, "достигнут лимит по времени", err)
}

func UnhandledError(err error) *apierror.ApiError {
	return apierror.NewApiError(http.StatusInternalServerError, "непредвиденная ошибка", err)
}

func EmptyDto(err error) *apierror.ApiError {
	return apierror.NewApiError(http.StatusInternalServerError, "отсутсвует dto", err)
}

func QuestionNotFound(err error) *apierror.ApiError {
	return apierror.NewApiError(http.StatusNotFound, "вопрос не найден", err)
}

func AnswerNotFound(err error) *apierror.ApiError {
	return apierror.NewApiError(http.StatusNotFound, "ответ не найден", err)
}

func TestDbErr(err error, cases ...*ErrDbCase) error {
	for _, c := range cases {
		if c.Func(err) {
			if c.CheckErr {
				return c.Creator(err)
			} else {
				return c.Creator(nil)
			}
		}
	}
	if storage.IsErrDeadline(err) {
		return DeadlineDbError(err)
	}

	return UnhandledError(err)
}
