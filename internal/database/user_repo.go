package database

import (
	"SOKOLOV_BOTS/internal/models"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = DB.Exec(
		"INSERT INTO users (email, password_hash) VALUES (?, ?)",
		email, string(hash),
	)
	return err
}

func GetUserByEmail(email string) (models.User, error) {
	var u models.User
	err := DB.QueryRow(
		"SELECT id, email, password_hash, created_at FROM users WHERE email = ?",
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	return u, err
}

func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}