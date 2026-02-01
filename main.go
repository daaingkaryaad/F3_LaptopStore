package main

import (
	"fmt"
	"net/http"

	"github.com/daaingkaryaad/F3_LaptopStore/internal/httpapi"
	"github.com/daaingkaryaad/F3_LaptopStore/internal/store"
)

func main() {
	st := store.NewStore()
	mux := http.NewServeMux()

	prodH := httpapi.NewProductHandler(st)
	prodH.Register(mux)

	cartH := httpapi.NewCartHandlers(st)
	orderH := httpapi.NewOrderHandlers(st)

	mux.Handle("/api/cart/items", httpapi.WithUser(http.HandlerFunc(cartH.AddToCart)))
	mux.Handle("/api/cart", httpapi.WithUser(http.HandlerFunc(cartH.GetCart)))
	mux.Handle("/api/orders", httpapi.WithUser(http.HandlerFunc(orderH.CreateOrder)))

	fmt.Println("server :8080")
	http.ListenAndServe(":8080", mux)

}
