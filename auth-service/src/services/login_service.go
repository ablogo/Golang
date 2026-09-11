package services

import "src/models"

func (u *UserService) Login(email string, password string) (token string) {

	user := u.GetUserByEmail(email)

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
