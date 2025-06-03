package validator

import (
	"errors"
	"time"
)

type Coupon struct {
	CouponCode          string
	ExpiryDate          time.Time
	UsageType           string
	ApplicableMedicineIDs []string
	ApplicableCategories []string
	MinOrderValue       float64
	ValidTimeWindow     time.Duration
	TermsAndConditions  string
	DiscountType        string
	DiscountValue       float64
	MaxUsagePerUser     int
}

func ValidateCoupon(coupon Coupon, orderTotal float64, currentTime time.Time) error {
	if coupon.CouponCode == "" {
		return errors.New("coupon code is required")
	}
	if currentTime.After(coupon.ExpiryDate) {
		return errors.New("coupon has expired")
	}
	if orderTotal < coupon.MinOrderValue {
		return errors.New("order total does not meet minimum order value")
	}
	// Additional validation logic can be added here
	return nil
}