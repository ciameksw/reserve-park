package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type RoleType string

const (
	RoleAdmin RoleType = "admin"
	RoleUser  RoleType = "user"
)

type authorizeResponse struct {
	Role   string `json:"role"`
	UserID string `json:"user_id"`
}

type contextKey string

const authorizeKey contextKey = "authorizeResponse"

func (s *Server) authorize(requiredRole RoleType, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		resp, err := s.UserService.Authorize(authHeader)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var authResp authorizeResponse
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		err = json.Unmarshal(body, &authResp)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if RoleType(authResp.Role) != requiredRole && RoleType(authResp.Role) != RoleAdmin {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), authorizeKey, authResp)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
