package apierror

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gengeo7/highlitent/logger"
	"github.com/gengeo7/highlitent/types/common"
)

type ErrorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

type ApiError struct {
	StatusCode    int
	OriginalError error
	Msg           string
}

type ValidationError struct {
	Fields map[string]string
	Msg    string
}

func (v ValidationError) Error() string {
	return v.Msg
}

func (a ApiError) Error() string {
	return a.Msg
}

func NewApiError(code int, msg string, err error) *ApiError {
	return &ApiError{
		StatusCode:    code,
		Msg:           msg,
		OriginalError: err,
	}
}

func NewValidationError(msg string, fields map[string]string) *ValidationError {
	return &ValidationError{
		Fields: fields,
		Msg:    msg,
	}
}

func SendError(w http.ResponseWriter, r *http.Request, err error) {
	var response ErrorResponse
	var statusCode int
	var ae *ApiError
	var ve *ValidationError
	if errors.As(err, &ae) {
		statusCode = ae.StatusCode
		response.Error = ae.Msg
		if r != nil && ae.OriginalError != nil {
			logger.Error("internal error",
				"route", r.RequestURI,
				"error", ae.OriginalError.Error(),
				"id", r.Context().Value(common.RequestIdKey{}),
			)
		}

	} else if errors.As(err, &ve) {
		statusCode = http.StatusBadRequest
		response.Fields = ve.Fields
		response.Error = ve.Msg

	} else {
		statusCode = http.StatusInternalServerError
		response.Error = "unhandled internal error"
		if r != nil {
			logger.Error("unhandled internal error",
				"route", r.RequestURI,
				"error", err.Error(),
				"id", r.Context().Value(common.RequestIdKey{}),
			)
		} else {
			logger.Error("unhandled internal error", "error", err.Error())
		}
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
