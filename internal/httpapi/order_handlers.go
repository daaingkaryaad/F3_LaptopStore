package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/daaingkaryaad/F3_LaptopStore/internal/store"
)

type createOrderReq struct {
	ItemIDs []string `json:"item_ids"`
}

type OrderHandlers struct {
	store *store.Store
}

func NewOrderHandlers(s *store.Store) *OrderHandlers {
	return &OrderHandlers{store: s}
}

func (h *OrderHandlers) HandleOrders(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateOrder(w, r)
	case http.MethodGet:
		h.ListOrders(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *OrderHandlers) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, 401, "no user")
		return
	}

	var req createOrderReq
	_ = json.NewDecoder(r.Body).Decode(&req)

	order, err := h.store.CreateOrderFromCart(userID, req.ItemIDs)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}

	writeJSON(w, 201, order)
}

func (h *OrderHandlers) ListOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, 401, "no user")
		return
	}

	orders, err := h.store.ListOrders(userID)
	if err != nil {
		writeError(w, 500, "failed to list orders")
		return
	}

	writeJSON(w, 200, orders)
}
