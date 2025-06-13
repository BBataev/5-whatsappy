package postgres

import "golang.org/x/crypto/bcrypt"

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

func CheckUserCredentials(username, password string) (bool, error) {
	var hashed string
	query := `SELECT password FROM users WHERE username = $1`
	err := db.Get(&hashed, query, username)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil, nil
}
