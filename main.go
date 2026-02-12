package main

import (
	"context"
	"log"
	"net/http"
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

	mux := http.NewServeMux()

	authH := httpapi.NewAuthHandlers(st)
	mux.Handle("/api/auth/admin/register", http.HandlerFunc(authH.AdminRegister))
	mux.Handle("/api/auth/register", http.HandlerFunc(authH.Register))
	mux.Handle("/api/auth/login", http.HandlerFunc(authH.Login))

	prodH := httpapi.NewProductHandler(st)
	mux.Handle("/api/laptops/compare", http.HandlerFunc(prodH.HandleCompare))
	mux.Handle("/api/laptops/", http.HandlerFunc(prodH.HandleLaptopByID))
	mux.Handle("/api/laptops", http.HandlerFunc(prodH.HandleLaptops))

	reviewH := httpapi.NewReviewHandlers(st)
	mux.Handle("/api/reviews/", httpapi.AuthRequired(http.HandlerFunc(reviewH.HandleReviewByID)))
	mux.Handle("/api/reviews", httpapi.AuthRequired(http.HandlerFunc(reviewH.HandleReviews)))

	cartH := httpapi.NewCartHandlers(st)
	orderH := httpapi.NewOrderHandlers(st)
	mux.Handle("/api/cart/items", httpapi.AuthRequired(http.HandlerFunc(cartH.AddToCart)))
	mux.Handle("/api/cart", httpapi.AuthRequired(http.HandlerFunc(cartH.GetCart)))
	mux.Handle("/api/orders", httpapi.AuthRequired(http.HandlerFunc(orderH.HandleOrders)))

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
