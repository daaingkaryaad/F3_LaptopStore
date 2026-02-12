package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/daaingkaryaad/F3_LaptopStore/internal/db"
	"github.com/daaingkaryaad/F3_LaptopStore/internal/httpapi"
	"github.com/daaingkaryaad/F3_LaptopStore/internal/store"
)

func main() {
	client, database, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	st := store.NewStore(database)

	_ = st.EnsureAdminUser(os.Getenv("ADMIN_EMAIL"), os.Getenv("ADMIN_FULL_NAME"), os.Getenv("ADMIN_PASSWORD"))

	mux := http.NewServeMux()

	authH := httpapi.NewAuthHandlers(st)
	mux.Handle("/api/auth/register", http.HandlerFunc(authH.Register))
	mux.Handle("/api/auth/login", http.HandlerFunc(authH.Login))

	prodH := httpapi.NewProductHandler(st)
	mux.Handle("/api/laptops/compare", http.HandlerFunc(prodH.HandleCompare))
	mux.Handle("/api/laptops/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			prodH.HandleLaptopByID(w, r)
			return
		}
		httpapi.AuthRequiredWithSession(st, httpapi.AdminOnly(http.HandlerFunc(prodH.HandleLaptopByID))).ServeHTTP(w, r)
	}))
	mux.Handle("/api/laptops", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			prodH.HandleLaptops(w, r)
			return
		}
		httpapi.AuthRequiredWithSession(st, httpapi.AdminOnly(http.HandlerFunc(prodH.HandleLaptops))).ServeHTTP(w, r)
	}))

	reviewH := httpapi.NewReviewHandlers(st)
	mux.Handle("/api/reviews/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			httpapi.AuthRequiredWithSession(st, httpapi.AdminOnly(http.HandlerFunc(reviewH.HandleReviewByID))).ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	mux.Handle("/api/reviews", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			httpapi.AuthOptionalWithSession(st, http.HandlerFunc(reviewH.HandleReviews)).ServeHTTP(w, r)
			return
		}
		if r.Method == http.MethodPost {
			httpapi.AuthRequiredWithSession(st, http.HandlerFunc(reviewH.HandleReviews)).ServeHTTP(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))

	cartH := httpapi.NewCartHandlers(st)
	orderH := httpapi.NewOrderHandlers(st)
	mux.Handle("/api/cart/items", httpapi.AuthRequiredWithSession(st, http.HandlerFunc(cartH.AddToCart)))
	mux.Handle("/api/cart", httpapi.AuthRequiredWithSession(st, http.HandlerFunc(cartH.GetCart)))
	mux.Handle("/api/orders", httpapi.AuthRequiredWithSession(st, http.HandlerFunc(orderH.HandleOrders)))

	fs := http.FileServer(http.Dir("frontend"))
	mux.Handle("/", fs)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("server :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
