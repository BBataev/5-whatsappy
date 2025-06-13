package postgres

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func AddNewUser(username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3)
	`
	_, err = db.Exec(query, username, email, string(hashedPassword))
	return err
}
