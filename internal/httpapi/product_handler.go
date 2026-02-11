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

func (h *ProductHandler) HandleLaptops(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		filter := productFilterFromQuery(r)
		products, err := h.store.ListProducts(filter)
		if err != nil {
			writeError(w, 500, "failed to list products")
			return
		}
		writeJSON(w, 200, products)

	case http.MethodPost:
		role, _ := RoleFromContext(r.Context())
		if role != "admin" {
			writeError(w, 403, "forbidden")
			return
		}

		var p model.Laptop
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			writeError(w, 400, "bad json")
			return
		}
		if err := validateProduct(p); err != nil {
			writeError(w, 400, err.Error())
			return
		}
		created, err := h.store.CreateProduct(p)
		if err != nil {
			writeError(w, 500, "create failed")
			return
		}
		writeJSON(w, 201, created)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) HandleLaptopByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/laptops/")
	if idStr == "" {
		writeError(w, 400, "bad id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		p, ok := h.store.GetProductByID(idStr)
		if !ok {
			writeError(w, 404, "not found")
			return
		}
		writeJSON(w, 200, p)

	case http.MethodPut:
		role, _ := RoleFromContext(r.Context())
		if role != "admin" {
			writeError(w, 403, "forbidden")
			return
		}
		var p model.Laptop
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			writeError(w, 400, "bad json")
			return
		}
		if err := validateProduct(p); err != nil {
			writeError(w, 400, err.Error())
			return
		}
		updated, ok := h.store.UpdateProduct(idStr, p)
		if !ok {
			writeError(w, 404, "not found")
			return
		}
		writeJSON(w, 200, updated)

	case http.MethodDelete:
		role, _ := RoleFromContext(r.Context())
		if role != "admin" {
			writeError(w, 403, "forbidden")
			return
		}
		if ok := h.store.DeleteProduct(idStr); !ok {
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

	firstID := r.URL.Query().Get("first")
	secondID := r.URL.Query().Get("second")
	if firstID == "" || secondID == "" {
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

func validateProduct(p model.Laptop) error {
	if strings.TrimSpace(p.ModelName) == "" {
		return httpError("model_name required")
	}
	if p.BrandID == "" || p.CategoryID == "" {
		return httpError("brand_id and category_id required")
	}
	if p.Price < 0 || p.Stock < 0 {
		return httpError("price and stock must be >= 0")
	}
	return nil
}

type httpError string

func (e httpError) Error() string { return string(e) }

func productFilterFromQuery(r *http.Request) store.ProductFilter {
	q := r.URL.Query()

	priceMin, _ := strconv.ParseFloat(q.Get("price_min"), 64)
	priceMax, _ := strconv.ParseFloat(q.Get("price_max"), 64)

	includeInactive := false
	if q.Get("include_inactive") == "true" {
		includeInactive = true
	}

	return store.ProductFilter{
		BrandID:         q.Get("brand"),
		CategoryID:      q.Get("category"),
		CPU:             q.Get("cpu"),
		RAM:             q.Get("ram"),
		GPU:             q.Get("gpu"),
		StorageType:     q.Get("storage_type"),
		PriceMin:        priceMin,
		PriceMax:        priceMax,
		Sort:            q.Get("sort"),
		IncludeInactive: includeInactive,
	}
}
