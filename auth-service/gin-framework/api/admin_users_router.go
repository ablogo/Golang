package api

import (
	"net/http"
	"strconv"

	"src/models"
	"src/services"

	"github.com/gin-gonic/gin"
)

type AdminRouter struct {
	UserSvc *services.UserService
}

func (s *AdminRouter) GetUser(c *gin.Context) {

	user_email := c.Query("email")
	if user_email == "" {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "User id is required"})
		return
	}

	user := s.UserSvc.GetUserByEmail(user_email)
	if user != nil {
		c.JSON(http.StatusOK, user)
		return
	} else {
		c.JSON(http.StatusNotFound, struct{ message string }{message: "Not Found"})
		return
	}
}

func (s *AdminRouter) GetUsers(c *gin.Context) {

	users := s.UserSvc.GetUsers()

	if users != nil {
		c.JSON(http.StatusOK, users)
		return
	} else {
		c.JSON(http.StatusNotFound, struct{ message string }{message: "Not Found"})
		return
	}
}

func (s *AdminRouter) DeleteUser(c *gin.Context) {

	user_email := c.Query("email")
	if user_email == "" {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "User id is required"})
		return
	}

	user := s.UserSvc.GetUserByEmail(user_email)
	if user == nil {
		c.JSON(http.StatusNotFound, struct{ message string }{message: "Not found"})
		return
	}

	result := s.UserSvc.DeleteUser(user.Id)
	if result {
		c.JSON(http.StatusOK, http.NoBody)
		return
	} else {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "Invalid input"})
		return
	}
}

func (s *AdminRouter) UpdateUser(c *gin.Context) {

	var model models.User

	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "Invalid input"})
		return
	}

	result := s.UserSvc.UpdateUser(model)
	if result {
		c.JSON(http.StatusOK, http.NoBody)
		return
	} else {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "Invalid input"})
		return
	}

}

func (s *AdminRouter) ChangePassword(c *gin.Context) {

	user_email := c.Query("email")
	if user_email == "" {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "User email is required"})
		return
	}

	password := c.Query("password")
	if password == "" {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "Invalid input"})
	}

	user := s.UserSvc.GetUserByEmail(user_email)
	if user != nil {
		c.JSON(http.StatusNotFound, http.NoBody)
		return
	}

	result := s.UserSvc.ChangePassword(user, password)
	if result {
		c.JSON(http.StatusOK, http.NoBody)
		return
	} else {
		c.JSON(http.StatusBadRequest, http.NoBody)
		return
	}

}

func (s *AdminRouter) AddAddress(c *gin.Context) {

	var model models.Address

	user_email := c.Query("email")
	if user_email == "" {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "User email is required"})
		return
	}

	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "Invalid input"})
		return
	}

	user := s.UserSvc.GetUserByEmail(user_email)
	if user == nil {
		c.JSON(http.StatusNotFound, struct{ message string }{message: "User not found"})
		return
	}

	address, result := s.UserSvc.AddAddress(user, model)
	if result {
		c.JSON(http.StatusOK, address)
	} else {
		c.JSON(http.StatusBadRequest, http.NoBody)
	}
}

func (s *AdminRouter) GetAddresses(c *gin.Context) {

	user_email := c.Query("email")
	if user_email == "" {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "User email is required"})
		return
	}

	user := s.UserSvc.GetUserByEmail(user_email)
	if user == nil {
		c.JSON(http.StatusNotFound, struct{ message string }{message: "Not found"})
		return
	}

	addresses := s.UserSvc.GetAddressByUser(user.Id)
	if len(*addresses) > 0 {
		c.JSON(http.StatusOK, addresses)
	} else {
		c.JSON(http.StatusNotFound, http.NoBody)
	}
}

func (s *AdminRouter) DeleteAddress(c *gin.Context) {

	user_email := c.Query("email")
	if user_email == "" {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "User email is required"})
		return
	}

	address_id, err := strconv.Atoi(c.Param("address_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "Invalid input"})
		return
	}

	user := s.UserSvc.GetUserByEmail(user_email)
	if user == nil {
		c.JSON(http.StatusNotFound, struct{ message string }{message: "Not found"})
		return
	}

	address := s.UserSvc.GetAddress(address_id)
	if address.UserId == user.Id {

		result := s.UserSvc.DeleteAddress(address_id)
		if result {
			c.JSON(http.StatusOK, http.StatusOK)
		} else {
			c.JSON(http.StatusBadRequest, http.NoBody)
		}
	} else {
		c.JSON(http.StatusNotAcceptable, struct{ message string }{message: "Don't belong to the user"})

	}
}

func (s *AdminRouter) UpdateAddress(c *gin.Context) {
	var model models.Address

	user_email := c.Query("email")
	if user_email == "" {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "User email is required"})
		return
	}

	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "Invalid input"})
		return
	}

	user := s.UserSvc.GetUserByEmail(user_email)
	if user == nil {
		c.JSON(http.StatusNotFound, struct{ message string }{message: "Not found"})
		return
	}

	if model.UserId == user.Id {

		result := s.UserSvc.UpdateAddress(model)
		if result {
			c.JSON(http.StatusOK, http.NoBody)
		} else {
			c.JSON(http.StatusBadRequest, http.NoBody)
		}
	} else {
		c.JSON(http.StatusNotAcceptable, struct{ message string }{message: "Don't belong to the user"})
	}
}
