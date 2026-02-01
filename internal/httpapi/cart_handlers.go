package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/daaingkaryaad/F3_LaptopStore/internal/store"
)

type CartHandlers struct {
	store *store.Store
}

func NewCartHandlers(s *store.Store) *CartHandlers {
	return &CartHandlers{store: s}
}

type addCartReq struct {
	UserID   int `json:"user_id"`
	LaptopID int `json:"laptop_id"`
	Quantity int `json:"quantity"`
}

func (h *CartHandlers) AddToCart(w http.ResponseWriter, r *http.Request) {
	var req addCartReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid json")
		return
	}

	cart := h.store.AddToCart(req.UserID, req.LaptopID, req.Quantity)
	writeJSON(w, 200, cart)
}

func (h *CartHandlers) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(CtxUserID).(int)
	if !ok {
		writeError(w, 401, "no user")
		return
	}

	cart := h.store.GetCart(userID)
	writeJSON(w, 200, cart)
}
