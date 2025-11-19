package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gengeo7/highlitent/apierror"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New()

type ValidateJsonKey struct{}

func ValidateJson[T any]() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body := r.Body
			defer body.Close()

			var v T
			decoder := json.NewDecoder(body)
			decoder.DisallowUnknownFields()
			err := decoder.Decode(&v)
			if err == io.EOF {
				apierror.SendError(w, r, apierror.NewApiError(http.StatusBadRequest, "пустое тело запроса", nil))
				return
			}
			if err != nil {
				apierror.SendError(w, r, apierror.NewApiError(http.StatusBadRequest, "ошибка чтения тела запроса", nil))
				return
			}

			err = validate.Struct(v)
			if err != nil {
				var invalidValidationError *validator.InvalidValidationError
				if errors.As(err, &invalidValidationError) {
					apierror.SendError(w, r, err)
					return
				}
				validationErrors := make(map[string]string)
				for _, e := range err.(validator.ValidationErrors) {
					field := e.Field()
					field = strings.ToLower(string(field[0])) + field[1:]
					tag := e.Tag()
					if e.Param() != "" {
						tag += ": " + e.Param()
					}
					validationErrors[field] = tag
				}

				apierror.SendError(w, r, apierror.NewValidationError("ошибка валидации", validationErrors))
				return
			}

			ctx := context.WithValue(r.Context(), ValidateJsonKey{}, v)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func DtoFromContext[T any](ctx context.Context) *T {
	val := ctx.Value(ValidateJsonKey{})
	if res, ok := val.(T); !ok {
		return nil
	} else {
		return &res
	}
}
