package services

import "src/models"

func Login(email string, password string) (token string) {

	user := GetUserByEmail(email)

	if user != nil {

		is_valid := verifyPassword(user.Password, password)

		if is_valid {
			token, _ = CreateToken(
				models.JWTClaims{
					UserId: user.Id,
					Name:   user.Name}, 60)
		}
	}
	return
}
