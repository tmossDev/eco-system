package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 14)
}

func ComparePassword(actualHashed []byte, comparingUnhashed string) bool {
	return bcrypt.CompareHashAndPassword(actualHashed, []byte(comparingUnhashed)) == nil
}
