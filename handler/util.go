package handler

import (
	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func CompareHashAndPassword(hashedPassword string, plainPassword string) bool {
	byteHash := []byte(hashedPassword)
	bytePlain := []byte(plainPassword)

	err := bcrypt.CompareHashAndPassword(byteHash, bytePlain)

	return err == nil
}
