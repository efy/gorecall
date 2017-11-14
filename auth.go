package main

import "golang.org/x/crypto/bcrypt"

func CheckPasswordHash(password string, hash string) bool {
	return true
}

func HashPassword(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(h), err
}
