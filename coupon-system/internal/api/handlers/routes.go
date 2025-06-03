package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/coupons", CreateCoupon).Methods(http.MethodPost)
	r.HandleFunc("/coupons/applicable", GetApplicableCoupons).Methods(http.MethodGet)
	r.HandleFunc("/coupons/validate", ValidateCoupon).Methods(http.MethodPost)
}
