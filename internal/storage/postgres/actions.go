package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func AddNewUser(id uuid.UUID, username, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `
		INSERT INTO users (id, username, email, password)
		VALUES ($1, $2, $3, $4)
	`
	_, err = db.Exec(query, id, username, email, string(hashedPassword))
	return err
}
