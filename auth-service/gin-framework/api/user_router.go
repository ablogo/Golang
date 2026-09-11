package api

import (
	"io"
	"net/http"
	"strconv"

	"src/models"

	"github.com/gin-gonic/gin"
)

func (a *Router) GetUser(c *gin.Context) {

	user_id, ok := c.MustGet("user_id").(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	user := a.UserSvc.GetUser(user_id, []string{"Address"})
	if user != nil {
		c.JSON(http.StatusOK, user)
		return
	} else {
		c.JSON(http.StatusNotFound, struct{ message string }{message: "Invalid input"})
		return
	}
}

func (a *Router) DeleteUser(c *gin.Context) {

	user_id, ok := c.MustGet("user_id").(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	result := a.UserSvc.DeleteUser(user_id)
	if result {
		c.JSON(http.StatusOK, http.NoBody)
		return
	} else {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "Invalid input"})
		return
	}
}

func (a *Router) UpdateUser(c *gin.Context) {

	var model models.User

	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "Invalid input"})
		return
	}

	result := a.UserSvc.UpdateUser(model)
	if result {
		c.JSON(http.StatusOK, http.NoBody)
		return
	} else {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "Invalid input"})
		return
	}
}

func (a *Router) ChangePassword(c *gin.Context) {

	user_id, ok := c.MustGet("user_id").(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	password := c.Query("password")
	if password == "" {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "Invalid input"})
	}

	user := a.UserSvc.GetUser(user_id, nil)
	if user == nil {
		c.JSON(http.StatusNotFound, http.NoBody)
		return
	}

	result := a.UserSvc.ChangePassword(user, password)
	if result {
		c.JSON(http.StatusOK, http.NoBody)
		return
	} else {
		c.JSON(http.StatusBadRequest, http.NoBody)
		return
	}

}

func (a *Router) AddPicture(c *gin.Context) {

	user_id, ok := c.MustGet("user_id").(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	user := a.UserSvc.GetUser(user_id, nil)
	if user == nil {
		c.JSON(http.StatusNotFound, http.NoBody)
		return
	}

	file_form, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "get form err: %s", err.Error())
		return
	}

	file, _ := file_form.Open()
	defer file.Close()
	fileBytes, _ := io.ReadAll(file)

	result := a.UserSvc.AddPicture(user, fileBytes, file_form.Header["Content-Type"][0], file_form.Filename)
	if result {
		c.JSON(http.StatusOK, http.NoBody)
	} else {
		c.JSON(http.StatusBadRequest, http.NoBody)
	}
}

func (a *Router) GetPicture(c *gin.Context) {
	user_id, ok := c.MustGet("user_id").(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	picture := a.UserSvc.GetImageByUser(user_id)
	if picture != nil {
		c.Data(http.StatusOK, *picture.ContentType, *picture.Picture)
		return
	} else {
		c.JSON(http.StatusNotFound, http.NoBody)
		return
	}
}

func (a *Router) AddAddress(c *gin.Context) {

	user_id, ok := c.MustGet("user_id").(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	var model models.Address
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "Invalid input"})
		return
	}

	user := a.UserSvc.GetUser(user_id, nil)
	address, result := a.UserSvc.AddAddress(user, model)
	if result {
		c.JSON(http.StatusOK, address)
	} else {
		c.JSON(http.StatusBadRequest, http.NoBody)
	}
}

func (a *Router) GetAddresses(c *gin.Context) {

	user_id, ok := c.MustGet("user_id").(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	addresses := a.UserSvc.GetAddressByUser(user_id)
	if len(*addresses) > 0 {
		c.JSON(http.StatusOK, addresses)
	} else {
		c.JSON(http.StatusBadRequest, http.NoBody)
	}
}

func (a *Router) DeleteAddress(c *gin.Context) {

	user_id, ok := c.MustGet("user_id").(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	address_id, err := strconv.Atoi(c.Param("address_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "Invalid input"})
		return
	}

	address := a.UserSvc.GetAddress(address_id)
	if address.UserId == user_id {

		result := a.UserSvc.DeleteAddress(address_id)
		if result {
			c.JSON(http.StatusOK, http.StatusOK)
		} else {
			c.JSON(http.StatusBadRequest, http.NoBody)
		}
	} else {
		c.JSON(http.StatusNotAcceptable, struct{ message string }{message: "Don't belong to the user"})
	}
}

func (a *Router) UpdateAddress(c *gin.Context) {
	var model models.Address

	user_id, ok := c.MustGet("user_id").(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "Invalid input"})
		return
	}

	if model.UserId == user_id {

		result := a.UserSvc.UpdateAddress(model)
		if result {
			c.JSON(http.StatusOK, http.NoBody)
		} else {
			c.JSON(http.StatusBadRequest, http.NoBody)
		}
	} else {
		c.JSON(http.StatusNotAcceptable, struct{ message string }{message: "Don't belong to the user"})
	}
}
