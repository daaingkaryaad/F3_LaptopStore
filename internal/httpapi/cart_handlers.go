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
	LaptopID int `json:"laptop_id"`
	Quantity int `json:"quantity"`
}

func (h *CartHandlers) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, 401, "no user")
		return
	}

	var req addCartReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid json")
		return
	}

	cart, err := h.store.AddToCart(userID, req.LaptopID, req.Quantity)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}

	writeJSON(w, 200, cart)
}

func (h *CartHandlers) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, 401, "no user")
		return
	}

	cart := h.store.GetCart(userID)
	writeJSON(w, 200, cart)
}
