package handlers

import (
	"coupon-system/internal/models"
	"coupon-system/internal/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

var CouponSvc *service.CouponService

// Handler for creating a coupon (for Gorilla Mux)
func CreateCoupon(w http.ResponseWriter, r *http.Request) {
	var coupon models.Coupon
	if err := json.NewDecoder(r.Body).Decode(&coupon); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := CouponSvc.CreateCoupon(r.Context(), &coupon); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(coupon); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

// GetApplicableCoupons retrieves applicable coupons based on cart items, order total, and timestamp
func GetApplicableCoupons(w http.ResponseWriter, r *http.Request) {
	type CouponSummary struct {
		CouponCode    string  `json:"coupon_code"`
		DiscountValue float64 `json:"discount_value"`
	}
	type ApplicableCouponsResponse struct {
		ApplicableCoupons []CouponSummary `json:"applicable_coupons"`
	}

	// Parse query params
	orderTotalStr := r.URL.Query().Get("order_total")
	timestamp := r.URL.Query().Get("timestamp")
	cartItemIDs := r.URL.Query()["cart_item_id"]
	cartItemCategories := r.URL.Query()["cart_item_category"]

	orderTotal, err := strconv.ParseFloat(orderTotalStr, 64)
	if err != nil {
		http.Error(w, "invalid order_total", http.StatusBadRequest)
		return
	}
	ts, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		http.Error(w, "invalid timestamp format", http.StatusBadRequest)
		return
	}

	// Build cart items from query params
	cartItems := make([]models.Item, 0, len(cartItemIDs))
	for i := range cartItemIDs {
		item := models.Item{
			ID:       cartItemIDs[i],
			Category: "",
		}
		if i < len(cartItemCategories) {
			item.Category = cartItemCategories[i]
		}
		cartItems = append(cartItems, item)
	}

	// Fetch all coupons
	allCoupons, _ := CouponSvc.ListAllCoupons(r.Context())

	applicable := make([]CouponSummary, 0)
	for _, coupon := range allCoupons {
		// Check expiry
		if coupon.ExpiryDate.Before(ts) {
			continue
		}
		// Check min order value
		if orderTotal < coupon.MinOrderValue {
			continue
		}
		// Check if any cart item matches medicine ID or category
		itemApplicable := false
		for _, item := range cartItems {
			for _, id := range coupon.ApplicableMedicineIDs {
				if item.ID == id {
					itemApplicable = true
					break
				}
			}
			if !itemApplicable {
				for _, cat := range coupon.ApplicableCategories {
					if item.Category == cat {
						itemApplicable = true
						break
					}
				}
			}
			if itemApplicable {
				break
			}
		}
		if !itemApplicable {
			continue
		}
		// Add to result
		applicable = append(applicable, CouponSummary{
			CouponCode:    coupon.CouponCode,
			DiscountValue: coupon.DiscountValue,
		})
	}

	resp := ApplicableCouponsResponse{ApplicableCoupons: applicable}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}

// ValidateCoupon validates a coupon for a given cart and order
func ValidateCoupon(w http.ResponseWriter, r *http.Request) {
	type ValidateCouponRequest struct {
		CouponCode string        `json:"coupon_code"`
		CartItems  []models.Item `json:"cart_items"`
		OrderTotal float64       `json:"order_total"`
		Timestamp  string        `json:"timestamp"`
	}

	type ValidateCouponResponse struct {
		IsValid  bool             `json:"is_valid"`
		Discount *models.Discount `json:"discount,omitempty"`
		Message  string           `json:"message,omitempty"`
		Reason   string           `json:"reason,omitempty"`
	}

	var req ValidateCouponRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.CouponCode == "" || req.Timestamp == "" {
		http.Error(w, "coupon_code and timestamp are required", http.StatusBadRequest)
		return
	}

	// Call the service layer for validation
	isValid, discount, err := CouponSvc.ValidateCoupon(r.Context(), req.CouponCode, req.CartItems, req.OrderTotal, req.Timestamp)
	resp := ValidateCouponResponse{IsValid: isValid}

	if isValid && err == nil {
		resp.Discount = discount
		resp.Message = "coupon applied successfully"
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("error encoding response: %v", err)
		}
		return
	}

	// Failure case
	resp.Reason = "coupon expired or not applicable"
	if err != nil {
		resp.Reason = err.Error()
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}
