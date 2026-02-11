package main

import (
	"log"
	"net/http"
	"time"

	"github.com/daaingkaryaad/F3_LaptopStore/internal/httpapi"
	"github.com/daaingkaryaad/F3_LaptopStore/internal/store"
)

func main() {
	st := store.NewStore()
	go st.CleanupInactiveProducts()

	mux := http.NewServeMux()

	authH := httpapi.NewAuthHandlers(st)
	mux.Handle("/api/auth/register", http.HandlerFunc(authH.Register))
	mux.Handle("/api/auth/login", http.HandlerFunc(authH.Login))

	prodH := httpapi.NewProductHandler(st)
	mux.Handle("/api/laptops", httpapi.AuthRequired(http.HandlerFunc(prodH.HandleLaptops)))
	mux.Handle("/api/laptops/", httpapi.AuthRequired(http.HandlerFunc(prodH.HandleLaptopByID)))
	mux.Handle("/api/laptops/compare", httpapi.AuthRequired(http.HandlerFunc(prodH.HandleCompare)))

	cartH := httpapi.NewCartHandlers(st)
	orderH := httpapi.NewOrderHandlers(st)

	mux.Handle("/api/cart/items", httpapi.AuthRequired(http.HandlerFunc(cartH.AddToCart)))
	mux.Handle("/api/cart", httpapi.AuthRequired(http.HandlerFunc(cartH.GetCart)))
	mux.Handle("/api/orders", httpapi.AuthRequired(http.HandlerFunc(orderH.HandleOrders)))

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
