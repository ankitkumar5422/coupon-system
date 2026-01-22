package repository

import (
	"context"
	"coupon-system/internal/models"
	"testing"
	"time"
)

func TestInMemoryCouponRepository_CreateCoupon(t *testing.T) {
	repo := NewCouponRepository()
	ctx := context.Background()

	coupon := &models.Coupon{
		CouponCode:    "TEST123",
		ExpiryDate:    time.Now().Add(24 * time.Hour),
		UsageType:     "one_time",
		DiscountType:  "fixed",
		DiscountValue: 10.0,
	}

	err := repo.CreateCoupon(ctx, coupon)
	if err != nil {
		t.Errorf("CreateCoupon() error = %v, want nil", err)
	}

	// Try to create duplicate
	err = repo.CreateCoupon(ctx, coupon)
	if err == nil {
		t.Error("CreateCoupon() with duplicate should return error")
	}
}

func TestInMemoryCouponRepository_GetCouponByCode(t *testing.T) {
	repo := NewCouponRepository()
	ctx := context.Background()

	coupon := &models.Coupon{
		CouponCode:    "TEST123",
		ExpiryDate:    time.Now().Add(24 * time.Hour),
		UsageType:     "one_time",
		DiscountType:  "fixed",
		DiscountValue: 10.0,
	}

	// Test getting non-existent coupon
	_, err := repo.GetCouponByCode(ctx, "NONEXISTENT")
	if err == nil {
		t.Error("GetCouponByCode() for non-existent coupon should return error")
	}

	// Create and get coupon
	_ = repo.CreateCoupon(ctx, coupon)
	retrieved, err := repo.GetCouponByCode(ctx, "TEST123")
	if err != nil {
		t.Errorf("GetCouponByCode() error = %v, want nil", err)
	}

	if retrieved.CouponCode != "TEST123" {
		t.Errorf("GetCouponByCode() = %v, want TEST123", retrieved.CouponCode)
	}
}

func TestInMemoryCouponRepository_UpdateCoupon(t *testing.T) {
	repo := NewCouponRepository()
	ctx := context.Background()

	coupon := &models.Coupon{
		CouponCode:    "TEST123",
		ExpiryDate:    time.Now().Add(24 * time.Hour),
		UsageType:     "one_time",
		DiscountType:  "fixed",
		DiscountValue: 10.0,
	}

	// Test updating non-existent coupon
	err := repo.UpdateCoupon(ctx, coupon)
	if err == nil {
		t.Error("UpdateCoupon() for non-existent coupon should return error")
	}

	// Create and update coupon
	_ = repo.CreateCoupon(ctx, coupon)
	coupon.DiscountValue = 20.0
	err = repo.UpdateCoupon(ctx, coupon)
	if err != nil {
		t.Errorf("UpdateCoupon() error = %v, want nil", err)
	}

	retrieved, _ := repo.GetCouponByCode(ctx, "TEST123")
	if retrieved.DiscountValue != 20.0 {
		t.Errorf("UpdateCoupon() discount value = %v, want 20.0", retrieved.DiscountValue)
	}
}

func TestInMemoryCouponRepository_DeleteCoupon(t *testing.T) {
	repo := NewCouponRepository()
	ctx := context.Background()

	coupon := &models.Coupon{
		CouponCode:    "TEST123",
		ExpiryDate:    time.Now().Add(24 * time.Hour),
		UsageType:     "one_time",
		DiscountType:  "fixed",
		DiscountValue: 10.0,
	}

	// Test deleting non-existent coupon
	err := repo.DeleteCoupon(ctx, "NONEXISTENT")
	if err == nil {
		t.Error("DeleteCoupon() for non-existent coupon should return error")
	}

	// Create and delete coupon
	_ = repo.CreateCoupon(ctx, coupon)
	err = repo.DeleteCoupon(ctx, "TEST123")
	if err != nil {
		t.Errorf("DeleteCoupon() error = %v, want nil", err)
	}

	_, err = repo.GetCouponByCode(ctx, "TEST123")
	if err == nil {
		t.Error("GetCouponByCode() after delete should return error")
	}
}

func TestInMemoryCouponRepository_ListAllCoupons(t *testing.T) {
	repo := NewCouponRepository()
	ctx := context.Background()

	// Test empty list
	coupons, err := repo.ListAllCoupons(ctx)
	if err != nil {
		t.Errorf("ListAllCoupons() error = %v, want nil", err)
	}
	if len(coupons) != 0 {
		t.Errorf("ListAllCoupons() count = %v, want 0", len(coupons))
	}

	// Add coupons
	_ = repo.CreateCoupon(ctx, &models.Coupon{
		CouponCode:    "TEST1",
		ExpiryDate:    time.Now().Add(24 * time.Hour),
		UsageType:     "one_time",
		DiscountType:  "fixed",
		DiscountValue: 10.0,
	})
	_ = repo.CreateCoupon(ctx, &models.Coupon{
		CouponCode:    "TEST2",
		ExpiryDate:    time.Now().Add(24 * time.Hour),
		UsageType:     "multi_use",
		DiscountType:  "percentage",
		DiscountValue: 15.0,
	})

	coupons, err = repo.ListAllCoupons(ctx)
	if err != nil {
		t.Errorf("ListAllCoupons() error = %v, want nil", err)
	}
	if len(coupons) != 2 {
		t.Errorf("ListAllCoupons() count = %v, want 2", len(coupons))
	}
}

func TestCachedCouponRepository(t *testing.T) {
	innerRepo := NewCouponRepository()
	cachedRepo := NewCachedCouponRepository(innerRepo, 1*time.Minute)
	ctx := context.Background()

	coupon := &models.Coupon{
		CouponCode:    "CACHED123",
		ExpiryDate:    time.Now().Add(24 * time.Hour),
		UsageType:     "one_time",
		DiscountType:  "fixed",
		DiscountValue: 10.0,
	}

	// Create through cached repo
	err := cachedRepo.CreateCoupon(ctx, coupon)
	if err != nil {
		t.Errorf("CreateCoupon() error = %v, want nil", err)
	}

	// Get coupon (should cache it)
	retrieved, err := cachedRepo.GetCouponByCode(ctx, "CACHED123")
	if err != nil {
		t.Errorf("GetCouponByCode() error = %v, want nil", err)
	}
	if retrieved.CouponCode != "CACHED123" {
		t.Errorf("GetCouponByCode() = %v, want CACHED123", retrieved.CouponCode)
	}

	// Get again (should come from cache)
	retrieved2, err := cachedRepo.GetCouponByCode(ctx, "CACHED123")
	if err != nil {
		t.Errorf("GetCouponByCode() (cached) error = %v, want nil", err)
	}
	if retrieved2.CouponCode != "CACHED123" {
		t.Errorf("GetCouponByCode() (cached) = %v, want CACHED123", retrieved2.CouponCode)
	}
}
