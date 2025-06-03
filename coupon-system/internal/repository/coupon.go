package repository

import (
	"context"
	"coupon-system/internal/cache"
	"coupon-system/internal/models"
	"database/sql"
	"errors"
	"sync"
	"time"
)

// In-memory repository (for testing/dev)
type InMemoryCouponRepository struct {
	mu      sync.RWMutex
	coupons map[string]*models.Coupon
}

func NewCouponRepository() *InMemoryCouponRepository {
	return &InMemoryCouponRepository{
		coupons: make(map[string]*models.Coupon),
	}
}

func (r *InMemoryCouponRepository) CreateCoupon(ctx context.Context, coupon *models.Coupon) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.coupons[coupon.CouponCode]; exists {
		return errors.New("coupon already exists")
	}
	r.coupons[coupon.CouponCode] = coupon
	return nil
}

func (r *InMemoryCouponRepository) GetCouponByCode(ctx context.Context, code string) (*models.Coupon, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	coupon, exists := r.coupons[code]
	if !exists {
		return nil, errors.New("coupon not found")
	}
	return coupon, nil
}

func (r *InMemoryCouponRepository) UpdateCoupon(ctx context.Context, coupon *models.Coupon) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.coupons[coupon.CouponCode]; !exists {
		return errors.New("coupon not found")
	}
	r.coupons[coupon.CouponCode] = coupon
	return nil
}

func (r *InMemoryCouponRepository) DeleteCoupon(ctx context.Context, code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.coupons[code]; !exists {
		return errors.New("coupon not found")
	}
	delete(r.coupons, code)
	return nil
}

func (r *InMemoryCouponRepository) ListAllCoupons(ctx context.Context) ([]*models.Coupon, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	coupons := make([]*models.Coupon, 0, len(r.coupons))
	for _, c := range r.coupons {
		coupons = append(coupons, c)
	}
	return coupons, nil
}

// Persistent SQLite repository
type SQLiteCouponRepository struct {
	db *sql.DB
}

func NewSQLiteCouponRepository(db *sql.DB) *SQLiteCouponRepository {
	return &SQLiteCouponRepository{db: db}
}

func (r *SQLiteCouponRepository) CreateCoupon(ctx context.Context, coupon *models.Coupon) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO coupons (coupon_code, expiry_date, usage_type, discount_type, discount_value, min_order_value)
        VALUES (?, ?, ?, ?, ?, ?)`,
		coupon.CouponCode, coupon.ExpiryDate, coupon.UsageType, coupon.DiscountType, coupon.DiscountValue, coupon.MinOrderValue)
	return err
}

func (r *SQLiteCouponRepository) GetCouponByCode(ctx context.Context, code string) (*models.Coupon, error) {
	row := r.db.QueryRowContext(ctx, `SELECT coupon_code, expiry_date, usage_type, discount_type, discount_value, min_order_value
        FROM coupons WHERE coupon_code = ?`, code)

	var coupon models.Coupon
	var expiry string
	err := row.Scan(
		&coupon.CouponCode,
		&expiry,
		&coupon.UsageType,
		&coupon.DiscountType,
		&coupon.DiscountValue,
		&coupon.MinOrderValue,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("coupon not found")
		}
		return nil, err
	}
	// Parse expiry date
	coupon.ExpiryDate, _ = time.Parse(time.RFC3339, expiry)
	return &coupon, nil
}

func (r *SQLiteCouponRepository) UpdateCoupon(ctx context.Context, coupon *models.Coupon) error {
	_, err := r.db.ExecContext(ctx, `UPDATE coupons SET expiry_date=?, usage_type=?, discount_type=?, discount_value=?, min_order_value=?
        WHERE coupon_code=?`,
		coupon.ExpiryDate, coupon.UsageType, coupon.DiscountType, coupon.DiscountValue, coupon.MinOrderValue, coupon.CouponCode)
	return err
}

func (r *SQLiteCouponRepository) DeleteCoupon(ctx context.Context, code string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM coupons WHERE coupon_code=?`, code)
	return err
}

func (r *SQLiteCouponRepository) ListAllCoupons(ctx context.Context) ([]*models.Coupon, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT coupon_code, expiry_date, usage_type, discount_type, discount_value, min_order_value FROM coupons`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coupons []*models.Coupon
	for rows.Next() {
		var coupon models.Coupon
		var expiry string
		err := rows.Scan(
			&coupon.CouponCode,
			&expiry,
			&coupon.UsageType,
			&coupon.DiscountType,
			&coupon.DiscountValue,
			&coupon.MinOrderValue,
		)
		if err == nil {
			coupon.ExpiryDate, _ = time.Parse(time.RFC3339, expiry)
			coupons = append(coupons, &coupon)
		}
	}
	return coupons, nil
}

// CouponRepository interface
type CouponRepository interface {
	CreateCoupon(ctx context.Context, coupon *models.Coupon) error
	GetCouponByCode(ctx context.Context, code string) (*models.Coupon, error)
	UpdateCoupon(ctx context.Context, coupon *models.Coupon) error
	DeleteCoupon(ctx context.Context, code string) error
	ListAllCoupons(ctx context.Context) ([]*models.Coupon, error)
}

// Caching wrapper for CouponRepository
type CachedCouponRepository struct {
	inner CouponRepository
	cache *cache.Cache
	ttl   time.Duration
}

func NewCachedCouponRepository(inner CouponRepository, ttl time.Duration) *CachedCouponRepository {
	return &CachedCouponRepository{
		inner: inner,
		cache: cache.NewCache(),
		ttl:   ttl,
	}
}

func (r *CachedCouponRepository) GetCouponByCode(ctx context.Context, code string) (*models.Coupon, error) {
	if val, ok := r.cache.Get(code); ok {
		return val.(*models.Coupon), nil
	}
	coupon, err := r.inner.GetCouponByCode(ctx, code)
	if err == nil {
		r.cache.Set(code, coupon, r.ttl)
	}
	return coupon, err
}

// Delegate other methods to inner
func (r *CachedCouponRepository) CreateCoupon(ctx context.Context, coupon *models.Coupon) error {
	return r.inner.CreateCoupon(ctx, coupon)
}
func (r *CachedCouponRepository) UpdateCoupon(ctx context.Context, coupon *models.Coupon) error {
	return r.inner.UpdateCoupon(ctx, coupon)
}
func (r *CachedCouponRepository) DeleteCoupon(ctx context.Context, code string) error {
	return r.inner.DeleteCoupon(ctx, code)
}
func (r *CachedCouponRepository) ListAllCoupons(ctx context.Context) ([]*models.Coupon, error) {
	return r.inner.ListAllCoupons(ctx)
}
