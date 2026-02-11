package httpapi

import (
	"context"
	"net/http"
	"strings"

	"github.com/daaingkaryaad/F3_LaptopStore/internal/auth"
)

type ctxKey string

const (
	CtxUserID ctxKey = "userID"
	CtxRole   ctxKey = "role"
)

func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := strings.TrimSpace(r.Header.Get("Authorization"))
		if raw == "" || !strings.HasPrefix(raw, "Bearer ") {
			writeError(w, 401, "missing bearer token")
			return
		}

		token := strings.TrimPrefix(raw, "Bearer ")
		claims, err := auth.ParseToken(token)
		if err != nil {
			writeError(w, 401, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), CtxUserID, claims.UserID)
		ctx = context.WithValue(ctx, CtxRole, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := RoleFromContext(r.Context())
		if !ok || role != "admin" {
			writeError(w, 403, "forbidden")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(CtxUserID)
	id, ok := v.(string)
	return id, ok
}

func RoleFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(CtxRole)
	role, ok := v.(string)
	return role, ok
}
