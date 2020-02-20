package util

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func IsEmailValid(email string) bool {
	regex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return regex.MatchString(email)
}

func IsPasswordSame(hashedPassword string, password []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), password)

	if err != nil {
		return false
	}

	return true
}
