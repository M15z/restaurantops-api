package handlers

import (
	"net/http"

	"github.com/m15z/restaurantops-api/internal/middleware"
)

func Me(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(int64)
	role, _ := r.Context().Value(middleware.RoleKey).(string)

	writeJSON(w, http.StatusOK, map[string]any{
		"user_id": userID,
		"role":    role,
	})
}
