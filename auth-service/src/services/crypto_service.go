package services

import "golang.org/x/crypto/bcrypt"

func hashPassword(value string) (p string) {

	hashed, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {

	} else {
		p = string(hashed)
	}
	return
}

func verifyPassword(hashed_password string, password string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(password))
	if err == nil {
		return true
	} else {
		return false
	}
}
