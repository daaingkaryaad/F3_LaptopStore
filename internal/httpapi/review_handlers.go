package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/daaingkaryaad/F3_LaptopStore/internal/store"
)

type ReviewHandlers struct {
	store *store.Store
}

func NewReviewHandlers(s *store.Store) *ReviewHandlers {
	return &ReviewHandlers{store: s}
}

type moderateReviewReq struct {
	Status string `json:"status"`
	Approved *bool `json:"approved,omitempty"`
}

type createReviewReq struct {
	LaptopID string `json:"laptop_id"`
	Rating   int    `json:"rating"`
	Comment  string `json:"comment"`
}

func (h *ReviewHandlers) HandleReviews(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateReview(w, r)
	case http.MethodGet:
		h.ListReviews(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *ReviewHandlers) HandleReviewByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/reviews/")
	if strings.HasSuffix(path, "/approve") {
		id := strings.TrimSuffix(path, "/approve")
		if id == "" {
			writeError(w, 400, "bad id")
			return
		}
		h.ApproveReview(w, r, id)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (h *ReviewHandlers) CreateReview(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		writeError(w, 401, "no user")
		return
	}

	var req createReviewReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "bad json")
		return
	}

	review, err := h.store.CreateReview(userID, req.LaptopID, req.Rating, req.Comment)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}

	writeJSON(w, 201, review)
}

func (h *ReviewHandlers) ListReviews(w http.ResponseWriter, r *http.Request) {
	productID := r.URL.Query().Get("product_id")
	if productID == "" {
		writeError(w, 400, "product_id required")
		return
	}

	role, _ := RoleFromContext(r.Context())
	includePending := false
	if role == "admin" && r.URL.Query().Get("all") == "true" {
		includePending = true
	}

	reviews, err := h.store.ListReviews(productID, includePending)
	if err != nil {
		writeError(w, 500, "failed to list reviews")
		return
	}

	writeJSON(w, 200, reviews)
}

func (h *ReviewHandlers) ApproveReview(w http.ResponseWriter, r *http.Request, id string) {
	role, _ := RoleFromContext(r.Context())
	if role != "admin" {
		writeError(w, 403, "forbidden")
		return
	}

	var req moderateReviewReq
	_ = json.NewDecoder(r.Body).Decode(&req)
	status := "approved"
	if req.Status != "" {
		status = req.Status
	} else if req.Approved != nil && !*req.Approved {
		status = "rejected"
	}

	updated, ok := h.store.SetReviewStatus(id, status)
	if !ok {
		writeError(w, 404, "not found")
		return
	}

	writeJSON(w, 200, updated)
}
