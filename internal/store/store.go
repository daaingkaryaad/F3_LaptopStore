package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/daaingkaryaad/F3_LaptopStore/internal/model"
)

type Store struct {
	mu       sync.RWMutex
	Products map[int]model.Product
	Carts    map[int]model.Cart
	Orders   map[int][]model.Order

	nextProductID int
	nextOrderID   int
}

func NewStore() *Store {
	return &Store{
		Products:      make(map[int]model.Product),
		Carts:         make(map[int]model.Cart),
		Orders:        make(map[int][]model.Order),
		nextProductID: 0,
		nextOrderID:   0,
	}
}

/* ---------------- Products ---------------- */

func (s *Store) ListProducts() []model.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]model.Product, 0, len(s.Products))
	for _, p := range s.Products {
		out = append(out, p)
	}
	return out
}

func (s *Store) GetProductByID(id int) (model.Product, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.Products[id]
	return p, ok
}

func (s *Store) CreateProduct(p model.Product) model.Product {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextProductID++
	p.ID = s.nextProductID
	s.Products[p.ID] = p
	return p
}

func (s *Store) UpdateProduct(id int, p model.Product) (model.Product, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Products[id]; !ok {
		return model.Product{}, false
	}
	p.ID = id
	s.Products[id] = p
	return p, true
}

func (s *Store) DeleteProduct(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Products[id]; !ok {
		return false
	}
	delete(s.Products, id)
	return true
}

/* ---------------- Cart ---------------- */

func (s *Store) AddToCart(userID, laptopID, qty int) model.Cart {
	if qty <= 0 {
		qty = 1
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	c := s.Carts[userID]
	c.UserID = userID
	c.Items = append(c.Items, model.CartItem{LaptopID: laptopID, Quantity: qty})
	s.Carts[userID] = c
	return c
}

func (s *Store) GetCart(userID int) model.Cart {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Carts[userID]
}

/* ---------------- Orders ---------------- */

func (s *Store) CreateOrderFromCart(userID int) (model.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cart := s.Carts[userID]
	if len(cart.Items) == 0 {
		return model.Order{}, fmt.Errorf("cart empty")
	}

	var items []model.OrderItem
	var total float64

	for _, it := range cart.Items {
		p, ok := s.Products[it.LaptopID]
		if !ok {
			return model.Order{}, fmt.Errorf("laptop %d not found", it.LaptopID)
		}
		items = append(items, model.OrderItem{
			LaptopID: it.LaptopID,
			Quantity: it.Quantity,
			Price:    p.Price,
		})
		total += p.Price * float64(it.Quantity)
	}

	s.nextOrderID++
	order := model.Order{
		ID:     s.nextOrderID,
		UserID: userID,
		Items:  items,
		Total:  total,
		Time:   time.Now(),
	}

	s.Orders[userID] = append(s.Orders[userID], order)
	delete(s.Carts, userID)

	return order, nil
}

func (s *Store) ListOrders(userID int) []model.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Orders[userID]
}
