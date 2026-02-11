package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/daaingkaryaad/F3_LaptopStore/internal/auth"
	"github.com/daaingkaryaad/F3_LaptopStore/internal/store"
)

type AuthHandlers struct {
	store *store.Store
}

func NewAuthHandlers(s *store.Store) *AuthHandlers {
	return &AuthHandlers{store: s}
}

type registerReq struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
	RoleID   int    `json:"role_id,omitempty"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResp struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "bad json")
		return
	}

	user, err := h.store.RegisterUser(req.Email, req.FullName, req.Password, req.RoleID)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}

	token, err := auth.GenerateToken(user.ID, user.RoleID)
	if err != nil {
		writeError(w, 500, "token error")
		return
	}

	writeJSON(w, 201, authResp{Token: token, User: user})
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "bad json")
		return
	}

	user, err := h.store.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		writeError(w, 401, "invalid credentials")
		return
	}

	token, err := auth.GenerateToken(user.ID, user.RoleID)
	if err != nil {
		writeError(w, 500, "token error")
		return
	}

	writeJSON(w, 200, authResp{Token: token, User: user})
}
