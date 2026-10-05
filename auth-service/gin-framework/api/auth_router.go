package api

import (
	"net/http"

	"src/model"
	"src/service"

	"github.com/gin-gonic/gin"
)

type Router struct {
	UserSvc *service.UserService
}

func (a *Router) SignUp(c *gin.Context) {
	var model model.SignUp

	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "Invalid input"})
		return
	}

	result := a.UserSvc.CreateUser(model)

	if result {
		c.JSON(http.StatusOK, http.NoBody)
		return
	} else {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
}

func (a *Router) SignIn(c *gin.Context) {
	var requestModel model.SignIn

	if err := c.ShouldBind(&requestModel); err != nil {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: err.Error()})
		return
	}

	token := a.UserSvc.Login(requestModel.UserName, requestModel.Password)

	if token != "" {
		c.JSON(http.StatusOK, model.Token{AccessToken: token, TokenType: "bearer"})
		return
	} else {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}

}

func VerifyToken(c *gin.Context) {

	_, ok := c.MustGet("user_id").(int)
	if ok {
		c.JSON(http.StatusOK, http.NoBody)

	} else {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "Invalid input"})
		return
	}
}
