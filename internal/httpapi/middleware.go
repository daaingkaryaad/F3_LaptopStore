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
	CtxRoleID ctxKey = "roleID"
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
		ctx = context.WithValue(ctx, CtxRoleID, claims.RoleID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireRoleID(roleID int, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid, ok := RoleIDFromContext(r.Context())
		if !ok || rid != roleID {
			writeError(w, 403, "forbidden")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func UserIDFromContext(ctx context.Context) (int, bool) {
	v := ctx.Value(CtxUserID)
	id, ok := v.(int)
	return id, ok
}

func RoleIDFromContext(ctx context.Context) (int, bool) {
	v := ctx.Value(CtxRoleID)
	id, ok := v.(int)
	return id, ok
}
