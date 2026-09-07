package api

import (
	"net/http"

	"src/models"
	"src/services"

	"github.com/gin-gonic/gin"
)

func SignUp(c *gin.Context) {
	var model models.SignUp

	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "Invalid input"})
		return
	}

	result := services.CreateUser(model)

	if result {
		c.JSON(http.StatusOK, http.NoBody)
		return
	} else {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
}

func SignIn(c *gin.Context) {
	var model models.SignIn

	if err := c.ShouldBind(&model); err != nil {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: err.Error()})
		return
	}

	token := services.Login(model.UserName, model.Password)

	if token != "" {
		c.JSON(http.StatusOK, models.Token{AccessToken: token, TokenType: "bearer"})
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
