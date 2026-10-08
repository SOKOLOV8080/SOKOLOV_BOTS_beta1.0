package handlers

import (
	"net/http"

	"SOKOLOV_BOTS/internal/database"
	"SOKOLOV_BOTS/internal/models"
)

func GetCurrentUser(r *http.Request) (models.User, bool) {
	cookie, err := r.Cookie("user_email")
	if err != nil {
		return models.User{}, false
	}

	user, err := database.GetUserByEmail(cookie.Value)
	if err != nil {
		return models.User{}, false
	}

	return user, true
}
