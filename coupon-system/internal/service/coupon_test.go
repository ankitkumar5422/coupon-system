package service

import (
	"context"
	"coupon-system/internal/models"
	"coupon-system/internal/repository"
	"testing"
	"time"
)

func TestCreateCoupon(t *testing.T) {
	// Setup
	repo := repository.NewCouponRepository()
	service := NewCouponService(repo)
	ctx := context.Background()

	// Test data
	coupon := &models.Coupon{
		CouponCode:    "TEST123",
		ExpiryDate:    time.Now().Add(24 * time.Hour),
		UsageType:     "one_time",
		DiscountType:  "fixed",
		DiscountValue: 10.0,
		MinOrderValue: 100.0,
	}

	// Execute
	err := service.CreateCoupon(ctx, coupon)

	// Assert
	if err != nil {
		t.Errorf("CreateCoupon() error = %v, want nil", err)
	}

	// Verify coupon was created
	retrieved, err := service.GetCouponByCode(ctx, "TEST123")
	if err != nil {
		t.Errorf("GetCouponByCode() error = %v, want nil", err)
	}

	if retrieved.CouponCode != "TEST123" {
		t.Errorf("GetCouponByCode() = %v, want TEST123", retrieved.CouponCode)
	}
}

func TestValidateCouponData(t *testing.T) {
	tests := []struct {
		name    string
		coupon  *models.Coupon
		wantErr bool
	}{
		{
			name: "valid coupon",
			coupon: &models.Coupon{
				CouponCode:    "VALID",
				ExpiryDate:    time.Now().Add(24 * time.Hour),
				UsageType:     "one_time",
				DiscountType:  "fixed",
				DiscountValue: 10.0,
			},
			wantErr: false,
		},
		{
			name: "empty coupon code",
			coupon: &models.Coupon{
				CouponCode:    "",
				ExpiryDate:    time.Now().Add(24 * time.Hour),
				UsageType:     "one_time",
				DiscountType:  "fixed",
				DiscountValue: 10.0,
			},
			wantErr: true,
		},
		{
			name: "expired coupon",
			coupon: &models.Coupon{
				CouponCode:    "EXPIRED",
				ExpiryDate:    time.Now().Add(-24 * time.Hour),
				UsageType:     "one_time",
				DiscountType:  "fixed",
				DiscountValue: 10.0,
			},
			wantErr: true,
		},
		{
			name: "invalid discount value",
			coupon: &models.Coupon{
				CouponCode:    "INVALID",
				ExpiryDate:    time.Now().Add(24 * time.Hour),
				UsageType:     "one_time",
				DiscountType:  "fixed",
				DiscountValue: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid usage type",
			coupon: &models.Coupon{
				CouponCode:    "INVALID_TYPE",
				ExpiryDate:    time.Now().Add(24 * time.Hour),
				UsageType:     "invalid_type",
				DiscountType:  "fixed",
				DiscountValue: 10.0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCouponData(tt.coupon)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCouponData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
