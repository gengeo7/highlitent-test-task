package middleware

import (
	"net/http"
	"time"

	"github.com/gengeo7/highlitent/logger"
	"github.com/gengeo7/highlitent/types/common"
)

func TimeElapsed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			elapsed := time.Since(start)
			logger.Info(
				"request completed",
				"elapsed_ms", elapsed.Milliseconds(),
				"path", r.URL.Path,
				"id", r.Context().Value(common.RequestIdKey{}),
			)
		}()
		next.ServeHTTP(w, r)
	})
}
