package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"ecommerce-backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserContextKey contextKey = "currentUser"
)

type UserClaims struct {
	UserID uint        `json:"userId"`
	Email  string      `json:"email"`
	Role   models.Role `json:"role"`
	jwt.RegisteredClaims
}

// ExtractClaims parses and validates token from Authorization header
func ExtractClaims(r *http.Request, jwtSecret string) (*UserClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header is missing")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, errors.New("invalid authorization header format")
	}

	tokenStr := parts[1]
	claims := &UserClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}

// RequireAuth ensures request has a valid token
func RequireAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := ExtractClaims(r, jwtSecret)
			if err != nil {
				http.Error(w, `{"error":"Sila log masuk untuk meneruskan"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAdmin ensures user is an authenticated Admin
func RequireAdmin(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := ExtractClaims(r, jwtSecret)
			if err != nil {
				http.Error(w, `{"error":"Akses tidak dibenarkan"}`, http.StatusUnauthorized)
				return
			}

			if claims.Role != models.RoleAdmin {
				http.Error(w, `{"error":"Hanya Administrator dibenarkan mengakses bahagian ini"}`, http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuth attaches user claims if token is present, but allows guest traffic if not
func OptionalAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := ExtractClaims(r, jwtSecret)
			if err == nil && claims != nil {
				ctx := context.WithValue(r.Context(), UserContextKey, claims)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// GetUserClaims retrieves claims stored in request context
func GetUserClaims(ctx context.Context) *UserClaims {
	claims, ok := ctx.Value(UserContextKey).(*UserClaims)
	if !ok {
		return nil
	}
	return claims
}
