package store

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/daaingkaryaad/F3_LaptopStore/internal/model"
	"golang.org/x/crypto/bcrypt"
)

const (
	RoleAdmin    = 1
	RoleCustomer = 2
)

type userRecord struct {
	User         model.User
	PasswordHash []byte
}

type Store struct {
	mu sync.RWMutex

	Products map[int]model.Product
	Carts    map[int]model.Cart
	Orders   map[int][]model.Order
	Users    map[string]userRecord

	nextProductID int
	nextOrderID   int
	nextUserID    int
}

func NewStore() *Store {
	return &Store{
		Products:      make(map[int]model.Product),
		Carts:         make(map[int]model.Cart),
		Orders:        make(map[int][]model.Order),
		Users:         make(map[string]userRecord),
		nextProductID: 0,
		nextOrderID:   0,
		nextUserID:    0,
	}
}

/* ---------------- Users ---------------- */

func (s *Store) RegisterUser(email, fullName, password string, roleID int) (model.User, error) {
	if email == "" || password == "" {
		return model.User{}, fmt.Errorf("email and password required")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Users[email]; exists {
		return model.User{}, fmt.Errorf("user already exists")
	}

	if roleID == 0 {
		roleID = RoleCustomer
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to hash password")
	}

	s.nextUserID++
	user := model.User{
		ID:        s.nextUserID,
		Email:     email,
		FullName:  fullName,
		RoleID:    roleID,
		CreatedAt: time.Now(),
	}

	s.Users[email] = userRecord{User: user, PasswordHash: hash}
	return user, nil
}

func (s *Store) AuthenticateUser(email, password string) (model.User, error) {
	s.mu.RLock()
	rec, ok := s.Users[email]
	s.mu.RUnlock()

	if !ok {
		return model.User{}, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword(rec.PasswordHash, []byte(password)); err != nil {
		return model.User{}, fmt.Errorf("invalid credentials")
	}

	return rec.User, nil
}

/* ---------------- Products ---------------- */

func (s *Store) ListProducts() []model.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]model.Product, 0, len(s.Products))
	for _, p := range s.Products {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
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

	if p.Stock < 0 {
		p.Stock = 0
	}
	if p.Price < 0 {
		p.Price = 0
	}
	if !p.IsActive {
		p.IsActive = true
	}

	s.nextProductID++
	p.ID = s.nextProductID
	s.Products[p.ID] = p
	return p
}

func (s *Store) UpdateProduct(id int, p model.Product) (model.Product, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.Products[id]
	if !ok {
		return model.Product{}, false
	}

	p.ID = id
	if p.Stock < 0 {
		p.Stock = 0
	}
	if p.Price < 0 {
		p.Price = 0
	}
	if p.ModelName == "" {
		p.ModelName = existing.ModelName
	}
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

func (s *Store) AddToCart(userID, laptopID, qty int) (model.Cart, error) {
	if qty <= 0 {
		qty = 1
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	product, ok := s.Products[laptopID]
	if !ok || !product.IsActive {
		return model.Cart{}, fmt.Errorf("laptop not found")
	}
	if product.Stock <= 0 {
		return model.Cart{}, fmt.Errorf("laptop out of stock")
	}

	cart := s.Carts[userID]
	cart.UserID = userID

	found := false
	for i := range cart.Items {
		if cart.Items[i].LaptopID == laptopID {
			cart.Items[i].Quantity += qty
			found = true
			break
		}
	}
	if !found {
		cart.Items = append(cart.Items, model.CartItem{
			LaptopID: laptopID,
			Quantity: qty,
		})
	}

	s.Carts[userID] = cart
	return cart, nil
}

func (s *Store) GetCart(userID int) model.Cart {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cart, ok := s.Carts[userID]
	if !ok {
		return model.Cart{UserID: userID, Items: []model.CartItem{}}
	}
	return cart
}

/* ---------------- Orders ---------------- */

func (s *Store) CreateOrderFromCart(userID int) (model.Order, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cart, ok := s.Carts[userID]
	if !ok || len(cart.Items) == 0 {
		return model.Order{}, fmt.Errorf("cart empty")
	}

	var items []model.OrderItem
	var total float64

	for _, it := range cart.Items {
		p, ok := s.Products[it.LaptopID]
		if !ok || !p.IsActive {
			return model.Order{}, fmt.Errorf("laptop %d not found", it.LaptopID)
		}
		if it.Quantity > p.Stock {
			return model.Order{}, fmt.Errorf("insufficient stock for laptop %d", it.LaptopID)
		}
		items = append(items, model.OrderItem{
			LaptopID: it.LaptopID,
			Quantity: it.Quantity,
			Price:    p.Price,
		})
		total += p.Price * float64(it.Quantity)
	}

	for _, it := range cart.Items {
		p := s.Products[it.LaptopID]
		p.Stock -= it.Quantity
		s.Products[it.LaptopID] = p
	}

	s.nextOrderID++
	order := model.Order{
		ID:        s.nextOrderID,
		UserID:    userID,
		Items:     items,
		Total:     total,
		Status:    "created",
		CreatedAt: time.Now(),
	}

	s.Orders[userID] = append(s.Orders[userID], order)
	delete(s.Carts, userID)

	return order, nil
}

func (s *Store) ListOrders(userID int) []model.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]model.Order, len(s.Orders[userID]))
	copy(out, s.Orders[userID])
	return out
}

func (s *Store) CleanupInactiveProducts() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()

		for id, product := range s.Products {
			if !product.IsActive {
				delete(s.Products, id)
			}
		}

		for userID, cart := range s.Carts {
			if len(cart.Items) == 0 {
				delete(s.Carts, userID)
			}
		}

		s.mu.Unlock()
	}
}
