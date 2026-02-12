package httpapi

import (
	"context"
	"net/http"
	"strings"

	"github.com/daaingkaryaad/F3_LaptopStore/internal/auth"
	"github.com/daaingkaryaad/F3_LaptopStore/internal/store"
)

type ctxKey string

const (
	CtxUserID ctxKey = "userID"
	CtxRole   ctxKey = "role"
)


// AuthRequiredWithSession validates JWT and also checks that the token exists in MongoDB sessions.
func AuthRequiredWithSession(st *store.Store, next http.Handler) http.Handler {
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
		ok, err := st.IsSessionValid(token)
		if err != nil {
			writeError(w, 500, "session check failed")
			return
		}
		if !ok {
			writeError(w, 401, "session expired or revoked")
			return
		}
		ctx := context.WithValue(r.Context(), CtxUserID, claims.UserID)
		ctx = context.WithValue(ctx, CtxRole, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthOptionalWithSession parses a bearer token if present; if present it must be valid and exist in sessions.
func AuthOptionalWithSession(st *store.Store, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := strings.TrimSpace(r.Header.Get("Authorization"))
		if raw == "" {
			next.ServeHTTP(w, r)
			return
		}
		if !strings.HasPrefix(raw, "Bearer ") {
			writeError(w, 401, "invalid authorization header")
			return
		}
		token := strings.TrimPrefix(raw, "Bearer ")
		claims, err := auth.ParseToken(token)
		if err != nil {
			writeError(w, 401, "invalid token")
			return
		}
		ok, err := st.IsSessionValid(token)
		if err != nil {
			writeError(w, 500, "session check failed")
			return
		}
		if !ok {
			writeError(w, 401, "session expired or revoked")
			return
		}
		ctx := context.WithValue(r.Context(), CtxUserID, claims.UserID)
		ctx = context.WithValue(ctx, CtxRole, claims.Role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
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

// AuthOptional parses a bearer token if present and, if valid, stores claims in context.
// If no token is present, it simply passes the request through.
// If a token is present but invalid, it returns 401.
func AuthOptional(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw := strings.TrimSpace(r.Header.Get("Authorization"))
		if raw == "" {
			next.ServeHTTP(w, r)
			return
		}
		if !strings.HasPrefix(raw, "Bearer ") {
			writeError(w, 401, "invalid authorization header")
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
