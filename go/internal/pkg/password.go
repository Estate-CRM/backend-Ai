package pkg

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type PasswordUtils struct{}

func (pu *PasswordUtils) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
func (pu *PasswordUtils) VerifyPassword(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, errors.New("invalid password")
	}
	return true, nil
}
