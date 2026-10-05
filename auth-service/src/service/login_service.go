package service

import "src/model"

func (u *UserService) Login(email string, password string) (token string) {

	user := u.GetUserByEmail(email)

	if user != nil {

		is_valid := verifyPassword(user.Password, password)

		if is_valid {
			token, _ = CreateToken(
				model.JWTClaims{
					UserId: user.Id,
					Name:   user.Name}, 60)
		}
	}
	return
}
