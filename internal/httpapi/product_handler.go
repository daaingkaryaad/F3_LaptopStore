package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/daaingkaryaad/F3_LaptopStore/internal/model"
	"github.com/daaingkaryaad/F3_LaptopStore/internal/store"
)

type ProductHandler struct {
	store *store.Store
}

func NewProductHandler(st *store.Store) *ProductHandler {
	return &ProductHandler{store: st}
}

func (h *ProductHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/laptops", h.HandleLaptops)
	mux.HandleFunc("/api/laptops/", h.HandleLaptopByID)
	mux.HandleFunc("/api/laptops/compare", h.HandleCompare)
}

func (h *ProductHandler) HandleLaptops(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		writeJSON(w, 200, h.store.ListProducts())

	case http.MethodPost:
		var p model.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			writeError(w, 400, "bad json")
			return
		}
		if err := validateProduct(p); err != nil {
			writeError(w, 400, err.Error())
			return
		}
		created := h.store.CreateProduct(p)
		writeJSON(w, 201, created)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) HandleLaptopByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/laptops/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeError(w, 400, "bad id")
		return
	}

	switch r.Method {

	case http.MethodGet:
		p, ok := h.store.GetProductByID(id)
		if !ok {
			writeError(w, 404, "not found")
			return
		}
		writeJSON(w, 200, p)

	case http.MethodPut:
		var p model.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			writeError(w, 400, "bad json")
			return
		}
		if err := validateProduct(p); err != nil {
			writeError(w, 400, err.Error())
			return
		}
		updated, ok := h.store.UpdateProduct(id, p)
		if !ok {
			writeError(w, 404, "not found")
			return
		}
		writeJSON(w, 200, updated)

	case http.MethodDelete:
		if ok := h.store.DeleteProduct(id); !ok {
			writeError(w, 404, "not found")
			return
		}
		writeJSON(w, 200, map[string]string{"message": "deleted"})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) HandleCompare(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	firstID, err1 := strconv.Atoi(r.URL.Query().Get("first"))
	secondID, err2 := strconv.Atoi(r.URL.Query().Get("second"))
	if err1 != nil || err2 != nil || firstID <= 0 || secondID <= 0 {
		writeError(w, 400, "query params first and second required")
		return
	}

	first, ok := h.store.GetProductByID(firstID)
	if !ok {
		writeError(w, 404, "first laptop not found")
		return
	}
	second, ok := h.store.GetProductByID(secondID)
	if !ok {
		writeError(w, 404, "second laptop not found")
		return
	}

	resp := map[string]any{
		"first":      first,
		"second":     second,
		"price_diff": first.Price - second.Price,
	}
	writeJSON(w, 200, resp)
}

func validateProduct(p model.Product) error {
	if strings.TrimSpace(p.ModelName) == "" {
		return httpError("model_name required")
	}
	if p.BrandID <= 0 || p.CategoryID <= 0 {
		return httpError("brand_id and category_id required")
	}
	if p.Price < 0 || p.Stock < 0 {
		return httpError("price and stock must be >= 0")
	}
	return nil
}

type httpError string

func (e httpError) Error() string { return string(e) }
