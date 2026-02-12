package store

import (
	"context"
	"fmt"
	"time"

	"github.com/daaingkaryaad/F3_LaptopStore/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type Store struct {
	db       *mongo.Database
	users    *mongo.Collection
	products *mongo.Collection
	carts    *mongo.Collection
	orders   *mongo.Collection
	reviews  *mongo.Collection
	sessions *mongo.Collection
}

func NewStore(db *mongo.Database) *Store {
	return &Store{
		db:       db,
		users:    db.Collection("users"),
		products: db.Collection("laptops"),
		carts:    db.Collection("carts"),
		orders:   db.Collection("orders"),
		reviews:  db.Collection("reviews"),
		sessions: db.Collection("sessions"),
	}
}

func (s *Store) ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

type ProductFilter struct {
	BrandID         string
	CategoryID      string
	CPU             string
	RAM             string
	GPU             string
	StorageType     string
	PriceMin        float64
	PriceMax        float64
	Sort            string
	IncludeInactive bool
}

func (s *Store) RegisterUser(email, fullName, password, role string) (model.User, error) {
	if email == "" || password == "" {
		return model.User{}, fmt.Errorf("email and password required")
	}
	if role == "" {
		role = "user"
	}

	ctx, cancel := s.ctx()
	defer cancel()

	count, err := s.users.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return model.User{}, err
	}
	if count > 0 {
		return model.User{}, fmt.Errorf("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to hash password")
	}

	user := model.User{
		ID:           primitive.NewObjectID(),
		Email:        email,
		FullName:     fullName,
		PasswordHash: string(hash),
		Role:         role,
		CreatedAt:    time.Now(),
	}

	if _, err := s.users.InsertOne(ctx, user); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *Store) AuthenticateUser(email, password string) (model.User, error) {
	ctx, cancel := s.ctx()
	defer cancel()

	var user model.User
	if err := s.users.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return model.User{}, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return model.User{}, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

func (s *Store) ListProducts(filter ProductFilter) ([]model.Laptop, error) {
	ctx, cancel := s.ctx()
	defer cancel()

	q := bson.M{}
	if filter.BrandID != "" {
		q["brand_id"] = filter.BrandID
	}
	if filter.CategoryID != "" {
		q["category_id"] = filter.CategoryID
	}
	if filter.CPU != "" {
		q["specs.cpu"] = filter.CPU
	}
	if filter.RAM != "" {
		q["specs.ram"] = filter.RAM
	}
	if filter.GPU != "" {
		q["specs.gpu"] = filter.GPU
	}
	if filter.StorageType != "" {
		q["specs.storage_type"] = filter.StorageType
	}
	if filter.PriceMin > 0 || filter.PriceMax > 0 {
		price := bson.M{}
		if filter.PriceMin > 0 {
			price["$gte"] = filter.PriceMin
		}
		if filter.PriceMax > 0 {
			price["$lte"] = filter.PriceMax
		}
		q["price"] = price
	}
	if !filter.IncludeInactive {
		q["is_active"] = true
	}

	opts := options.Find()
	switch filter.Sort {
	case "price_asc":
		opts.SetSort(bson.D{{Key: "price", Value: 1}})
	case "price_desc":
		opts.SetSort(bson.D{{Key: "price", Value: -1}})
	case "newest":
		opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	}

	cur, err := s.products.Find(ctx, q, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []model.Laptop
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Store) GetProductByID(id string) (model.Laptop, bool) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Laptop{}, false
	}

	ctx, cancel := s.ctx()
	defer cancel()

	var p model.Laptop
	if err := s.products.FindOne(ctx, bson.M{"_id": oid}).Decode(&p); err != nil {
		return model.Laptop{}, false
	}
	return p, true
}

func (s *Store) CreateProduct(p model.Laptop) (model.Laptop, error) {
	ctx, cancel := s.ctx()
	defer cancel()

	p.ID = primitive.NewObjectID()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	if !p.IsActive {
		p.IsActive = true
	}

	if _, err := s.products.InsertOne(ctx, p); err != nil {
		return model.Laptop{}, err
	}

	return p, nil
}

func (s *Store) UpdateProduct(id string, p model.Laptop) (model.Laptop, bool) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Laptop{}, false
	}

	ctx, cancel := s.ctx()
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"model_name":  p.ModelName,
			"brand_id":    p.BrandID,
			"category_id": p.CategoryID,
			"price":       p.Price,
			"stock":       p.Stock,
			"description": p.Description,
			"is_active":   p.IsActive,
			"specs":       p.Specs,
			"updated_at":  time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated model.Laptop
	if err := s.products.FindOneAndUpdate(ctx, bson.M{"_id": oid}, update, opts).Decode(&updated); err != nil {
		return model.Laptop{}, false
	}
	return updated, true
}

func (s *Store) DeleteProduct(id string) bool {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false
	}

	ctx, cancel := s.ctx()
	defer cancel()

	res, err := s.products.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return false
	}
	return res.DeletedCount > 0
}

func (s *Store) AddToCart(userID, laptopID string, qty int) (model.Cart, error) {
	if qty <= 0 {
		qty = 1
	}

	if _, ok := s.GetProductByID(laptopID); !ok {
		return model.Cart{}, fmt.Errorf("laptop not found")
	}

	ctx, cancel := s.ctx()
	defer cancel()

	var cart model.Cart
	err := s.carts.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		cart = model.Cart{
			ID:        primitive.NewObjectID(),
			UserID:    userID,
			Items:     []model.CartItem{},
			UpdatedAt: time.Now(),
		}
	}

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

	cart.UpdatedAt = time.Now()

	_, err = s.carts.ReplaceOne(ctx, bson.M{"user_id": userID}, cart, options.Replace().SetUpsert(true))
	if err != nil {
		return model.Cart{}, err
	}

	return cart, nil
}

func (s *Store) GetCart(userID string) (model.Cart, error) {
	ctx, cancel := s.ctx()
	defer cancel()

	var cart model.Cart
	err := s.carts.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart)
	if err != nil {
		return model.Cart{
			UserID: userID,
			Items:  []model.CartItem{},
		}, nil
	}
	return cart, nil
}

func (s *Store) CreateOrderFromCart(userID string, itemIDs []string) (model.Order, error) {
	ctx, cancel := s.ctx()
	defer cancel()

	var cart model.Cart
	if err := s.carts.FindOne(ctx, bson.M{"user_id": userID}).Decode(&cart); err != nil {
		return model.Order{}, fmt.Errorf("cart empty")
	}
	if len(cart.Items) == 0 {
		return model.Order{}, fmt.Errorf("cart empty")
	}

	selected := cart.Items
	if len(itemIDs) > 0 {
		set := map[string]struct{}{}
		for _, id := range itemIDs {
			set[id] = struct{}{}
		}
		filtered := make([]model.CartItem, 0, len(cart.Items))
		for _, it := range cart.Items {
			if _, ok := set[it.LaptopID]; ok {
				filtered = append(filtered, it)
			}
		}
		if len(filtered) == 0 {
			return model.Order{}, fmt.Errorf("no selected items")
		}
		selected = filtered
	}

	var items []model.OrderItem
	var total float64

	for _, it := range selected {
		p, ok := s.GetProductByID(it.LaptopID)
		if !ok || !p.IsActive {
			return model.Order{}, fmt.Errorf("laptop not found")
		}
		if it.Quantity > p.Stock {
			return model.Order{}, fmt.Errorf("insufficient stock")
		}
		items = append(items, model.OrderItem{
			LaptopID: it.LaptopID,
			Quantity: it.Quantity,
			Price:    p.Price,
		})
		total += p.Price * float64(it.Quantity)
	}

	for _, it := range selected {
		oid, err := primitive.ObjectIDFromHex(it.LaptopID)
		if err != nil {
			return model.Order{}, fmt.Errorf("invalid laptop id")
		}
		_, err = s.products.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{
			"$inc": bson.M{"stock": -it.Quantity},
			"$set": bson.M{"updated_at": time.Now()},
		})
		if err != nil {
			return model.Order{}, err
		}
	}

	order := model.Order{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		Items:     items,
		Total:     total,
		Status:    "created",
		CreatedAt: time.Now(),
	}

	if _, err := s.orders.InsertOne(ctx, order); err != nil {
		return model.Order{}, err
	}

	if len(itemIDs) == 0 {
		_, _ = s.carts.DeleteOne(ctx, bson.M{"user_id": userID})
	} else {
		remaining := []model.CartItem{}
		set := map[string]bool{}
		for _, id := range itemIDs {
			set[id] = true
		}
		for _, it := range cart.Items {
			if !set[it.LaptopID] {
				remaining = append(remaining, it)
			}
		}
		cart.Items = remaining
		cart.UpdatedAt = time.Now()
		_, _ = s.carts.ReplaceOne(ctx, bson.M{"user_id": userID}, cart, options.Replace().SetUpsert(true))
	}

	return order, nil
}

func (s *Store) ListOrders(userID string) ([]model.Order, error) {
	ctx, cancel := s.ctx()
	defer cancel()

	cur, err := s.orders.Find(ctx, bson.M{"user_id": userID}, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []model.Order
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Store) CreateReview(userID, laptopID string, rating int, comment string) (model.Review, error) {
	if rating < 1 || rating > 5 {
		return model.Review{}, fmt.Errorf("rating must be 1-5")
	}
	if _, ok := s.GetProductByID(laptopID); !ok {
		return model.Review{}, fmt.Errorf("laptop not found")
	}

	ctx, cancel := s.ctx()
	defer cancel()

	review := model.Review{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		LaptopID:  laptopID,
		Rating:    rating,
		Comment:   comment,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	if _, err := s.reviews.InsertOne(ctx, review); err != nil {
		return model.Review{}, err
	}

	return review, nil
}

func (s *Store) ListReviews(laptopID string, includePending bool) ([]model.Review, error) {
	ctx, cancel := s.ctx()
	defer cancel()

	filter := bson.M{"laptop_id": laptopID}
	if !includePending {
		filter["status"] = "approved"
	}

	cur, err := s.reviews.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var out []model.Review
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *Store) SetReviewStatus(id string, status string) (model.Review, bool) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Review{}, false
	}

	ctx, cancel := s.ctx()
	defer cancel()

	if status == "" {
		status = "approved"
	}

	update := bson.M{"$set": bson.M{"status": status}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var out model.Review
	if err := s.reviews.FindOneAndUpdate(ctx, bson.M{"_id": oid}, update, opts).Decode(&out); err != nil {
		return model.Review{}, false
	}
	return out, true
}

func (s *Store) ApproveReview(id string) (model.Review, bool) {
	return s.SetReviewStatus(id, "approved")
}

func (s *Store) EnsureAdminUser(email, fullName, password string) error {
	if email == "" {
		email = "admin@rapidtech.local"
	}
	if fullName == "" {
		fullName = "Admin"
	}
	if password == "" {
		password = "Admin123!"
	}

	ctx, cancel := s.ctx()
	defer cancel()

	count, err := s.users.CountDocuments(ctx, bson.M{"role": "admin"})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	_, err = s.RegisterUser(email, fullName, password, "admin")
	return err
}

func (s *Store) CreateSession(token, userID, role string, expiresAt time.Time) error {
	ctx, cancel := s.ctx()
	defer cancel()

	_, err := s.sessions.InsertOne(ctx, model.Session{
		ID:        primitive.NewObjectID(),
		Token:     token,
		UserID:    userID,
		Role:      role,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	})
	return err
}

func (s *Store) IsSessionValid(token string) (bool, error) {
	ctx, cancel := s.ctx()
	defer cancel()

	var sess model.Session
	err := s.sessions.FindOne(ctx, bson.M{"token": token}).Decode(&sess)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	if time.Now().After(sess.ExpiresAt) {
		_, _ = s.sessions.DeleteOne(ctx, bson.M{"token": token})
		return false, nil
	}
	return true, nil
}
