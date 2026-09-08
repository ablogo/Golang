package models

type SignIn struct {
	UserName string `form:"username" binding:"required,email"`
	Password string `form:"password" binding:"required,min=5"`
}
