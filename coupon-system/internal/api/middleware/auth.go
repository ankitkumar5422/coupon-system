package middleware

import (
    "net/http"
    "strings"
)

// AuthMiddleware is a middleware function that checks for a valid authorization token.
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if !isValidToken(token) {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// isValidToken checks if the provided token is valid.
func isValidToken(token string) bool {
    // Here you can implement your token validation logic.
    // For now, we will just check if the token is not empty and starts with "Bearer ".
    return strings.HasPrefix(token, "Bearer ") && len(strings.TrimSpace(token)) > len("Bearer ")
}