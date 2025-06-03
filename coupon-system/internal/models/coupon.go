package models

import "time"

type Coupon struct {
	ID                    string    `json:"id"`
	CouponCode            string    `json:"coupon_code"`
	ExpiryDate            time.Time `json:"expiry_date"`
	UsageType             string    `json:"usage_type"` // "one_time", "multi_use", or "time_based"
	ApplicableMedicineIDs []string  `json:"applicable_medicine_ids"`
	ApplicableCategories  []string  `json:"applicable_categories"`
	MinOrderValue         float64   `json:"min_order_value"`
	ValidTimeWindow       string    `json:"valid_time_window"`
	TermsAndConditions    string    `json:"terms_and_conditions"`
	DiscountType          string    `json:"discount_type"` // e.g., "percentage", "fixed"
	DiscountValue         float64   `json:"discount_value"`
	MaxUsagePerUser       int       `json:"max_usage_per_user"`
	IsActive              bool      `json:"is_active"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// Item represents a cart item in the system
type Item struct {
	ID       string  `json:"id"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type Discount struct {
	ItemsDiscount   float64 `json:"items_discount"`
	ChargesDiscount float64 `json:"charges_discount"`
	TotalDiscount   float64 `json:"total_discount"`
}
