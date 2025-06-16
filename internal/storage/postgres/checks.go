package postgres

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CheckUserConflict(username, email string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM users
			WHERE username = $1 OR email = $2
		)
	`
	err := db.Get(&exists, query, username, email)
	return exists, err
}

func CheckUserCredentials(username, password string) (uuid.UUID, bool, error) {
	var (
		id     uuid.UUID
		hashed string
	)
	query := `SELECT id, password FROM users WHERE username = $1`
	err := db.QueryRow(query, username).Scan(&id, &hashed)
	if err != nil {
		return uuid.Nil, false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		return uuid.Nil, false, nil
	}

	return id, true, nil
}
