package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func PasswordHash(target string) ([]byte, error) {

	password_byte := []byte(target)

	salt := 10

	hash_password, err := bcrypt.GenerateFromPassword(password_byte, salt)
	if err != nil {
		return nil, err
	}

	return hash_password, nil
}
