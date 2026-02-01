package httpapi

import (
	"context"
	"net/http"
	"strconv"
)

type ctxKey string

const CtxUserID ctxKey = "userID"

func WithUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := r.Header.Get("X-User-Id")
		if raw == "" {
			writeError(w, 401, "missing X-User-Id")
			return
		}

		id, err := strconv.Atoi(raw)
		if err != nil {
			writeError(w, 400, "bad user id")
			return
		}

		ctx := context.WithValue(r.Context(), CtxUserID, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRole(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Role") != role {
			writeError(w, 403, "forbidden")
			return
		}
		next.ServeHTTP(w, r)
	})
}
