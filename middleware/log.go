package middleware

import (
	"context"
	"net/http"

	"github.com/gengeo7/highlitent/logger"
	"github.com/gengeo7/highlitent/types/common"
	"github.com/google/uuid"
)

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		logger.Info("new request", "ip", r.RemoteAddr, "route", r.RequestURI, "id", id)
		ctx := context.WithValue(r.Context(), common.RequestIdKey{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
