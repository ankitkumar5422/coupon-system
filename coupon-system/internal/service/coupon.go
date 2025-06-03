package service

import (
	"context"
	"coupon-system/internal/models"
	"coupon-system/internal/repository"
	"errors"
	"time"
)

type CouponService struct {
	repo repository.CouponRepository
}

func NewCouponService(repo repository.CouponRepository) *CouponService {
	return &CouponService{repo: repo}
}

func (s *CouponService) GetCouponByCode(ctx context.Context, code string) (*models.Coupon, error) {
	return s.repo.GetCouponByCode(ctx, code)
}

func (s *CouponService) ListAllCoupons(ctx context.Context) ([]*models.Coupon, error) {
	return s.repo.ListAllCoupons(ctx)
}

func (s *CouponService) CreateCoupon(ctx context.Context, coupon *models.Coupon) error {
	coupon.CreatedAt = time.Now()
	coupon.UpdatedAt = time.Now()
	coupon.IsActive = true

	if err := validateCouponData(coupon); err != nil {
		return err
	}
	return s.repo.CreateCoupon(ctx, coupon)
}

func (s *CouponService) ValidateCoupon(ctx context.Context, couponCode string, items []models.Item, orderTotal float64, timestamp string) (bool, *models.Discount, error) {
	coupon, err := s.repo.GetCouponByCode(ctx, couponCode)
	if err != nil {
		return false, nil, err
	}
	ts, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return false, nil, err
	}
	if coupon.ExpiryDate.Before(ts) {
		return false, nil, errors.New("coupon has expired")
	}
	discount := calculateDiscount(coupon, items, orderTotal)
	return true, discount, nil
}

func validateCouponData(coupon *models.Coupon) error {
	if coupon.CouponCode == "" {
		return errors.New("coupon code is required")
	}
	if coupon.ExpiryDate.Before(time.Now()) {
		return errors.New("expiry date must be in the future")
	}
	if coupon.DiscountValue <= 0 {
		return errors.New("discount value must be greater than 0")
	}
	switch coupon.UsageType {
	case "one_time", "multi_use", "time_based":
	default:
		return errors.New("invalid usage type")
	}
	return nil
}

func calculateDiscount(coupon *models.Coupon, items []models.Item, orderTotal float64) *models.Discount {
	var itemsDiscount float64
	var chargesDiscount float64

	for _, item := range items {
		isApplicable := false
		for _, id := range coupon.ApplicableMedicineIDs {
			if item.ID == id {
				isApplicable = true
				break
			}
		}
		if !isApplicable {
			for _, cat := range coupon.ApplicableCategories {
				if item.Category == cat {
					isApplicable = true
					break
				}
			}
		}
		if isApplicable {
			if coupon.DiscountType == "percentage" {
				itemsDiscount += item.Price * float64(item.Quantity) * (coupon.DiscountValue / 100)
			} else {
				itemsDiscount += coupon.DiscountValue
			}
		}
	}
	if coupon.DiscountType == "charges" {
		chargesDiscount = coupon.DiscountValue
	}
	return &models.Discount{
		ItemsDiscount:   itemsDiscount,
		ChargesDiscount: chargesDiscount,
		TotalDiscount:   itemsDiscount + chargesDiscount,
	}
}
