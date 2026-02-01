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
	mux.HandleFunc("/api/laptops", h.handleLaptops)     // GET, POST
	mux.HandleFunc("/api/laptops/", h.handleLaptopByID) // GET, PUT, DELETE
}

func (h *ProductHandler) handleLaptops(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		writeJSON(w, 200, h.store.ListProducts())

	case http.MethodPost:
		var p model.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			writeError(w, 400, "bad json")
			return
		}
		created := h.store.CreateProduct(p)
		writeJSON(w, 201, created)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) handleLaptopByID(w http.ResponseWriter, r *http.Request) {
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
