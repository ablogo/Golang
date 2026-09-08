package api

import (
	"io"
	"net/http"
	"strconv"

	"src/models"
	"src/services"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {

	user_id, ok := c.MustGet("user_id").(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	user := services.GetUser(user_id, []string{"Address"})
	if user != nil {
		c.JSON(http.StatusOK, user)
		return
	} else {
		c.JSON(http.StatusNotFound, struct{ message string }{message: "Invalid input"})
		return
	}
}

func DeleteUser(c *gin.Context) {

	user_id, ok := c.MustGet("user_id").(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	result := services.DeleteUser(user_id)
	if result {
		c.JSON(http.StatusOK, http.NoBody)
		return
	} else {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "Invalid input"})
		return
	}
}

func UpdateUser(c *gin.Context) {

	var model models.User

	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "Invalid input"})
		return
	}

	result := services.UpdateUser(model)
	if result {
		c.JSON(http.StatusOK, http.NoBody)
		return
	} else {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "Invalid input"})
		return
	}
}

func ChangePassword(c *gin.Context) {

	user_id, ok := c.MustGet("user_id").(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	password := c.Query("password")
	if password == "" {
		c.JSON(http.StatusBadRequest, struct{ message string }{message: "Invalid input"})
	}

	user := services.GetUser(user_id, nil)
	if user == nil {
		c.JSON(http.StatusNotFound, http.NoBody)
		return
	}

	result := services.ChangePassword(user, password)
	if result {
		c.JSON(http.StatusOK, http.NoBody)
		return
	} else {
		c.JSON(http.StatusBadRequest, http.NoBody)
		return
	}

}

func AddPicture(c *gin.Context) {

	user_id, ok := c.MustGet("user_id").(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	user := services.GetUser(user_id, nil)
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

	result := services.AddPicture(user, fileBytes, file_form.Header["Content-Type"][0], file_form.Filename)
	if result {
		c.JSON(http.StatusOK, http.NoBody)
	} else {
		c.JSON(http.StatusBadRequest, http.NoBody)
	}
}

func GetPicture(c *gin.Context) {
	user_id, ok := c.MustGet("user_id").(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	picture := services.GetImageByUser(user_id)
	if picture != nil {
		c.Data(http.StatusOK, *picture.ContentType, *picture.Picture)
		return
	} else {
		c.JSON(http.StatusNotFound, http.NoBody)
		return
	}
}

func AddAddress(c *gin.Context) {

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

	user := services.GetUser(user_id, nil)
	address, result := services.AddAddress(user, model)
	if result {
		c.JSON(http.StatusOK, address)
	} else {
		c.JSON(http.StatusBadRequest, http.NoBody)
	}
}

func GetAddresses(c *gin.Context) {

	user_id, ok := c.MustGet("user_id").(int)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	addresses := services.GetAddressByUser(user_id)
	if len(*addresses) > 0 {
		c.JSON(http.StatusOK, addresses)
	} else {
		c.JSON(http.StatusBadRequest, http.NoBody)
	}
}

func DeleteAddress(c *gin.Context) {

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

	address := services.GetAddress(address_id)
	if address.UserId == user_id {

		result := services.DeleteAddress(address_id)
		if result {
			c.JSON(http.StatusOK, http.StatusOK)
		} else {
			c.JSON(http.StatusBadRequest, http.NoBody)
		}
	} else {
		c.JSON(http.StatusNotAcceptable, struct{ message string }{message: "Don't belong to the user"})
	}
}

func UpdateAddress(c *gin.Context) {
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

		result := services.UpdateAddress(model)
		if result {
			c.JSON(http.StatusOK, http.NoBody)
		} else {
			c.JSON(http.StatusBadRequest, http.NoBody)
		}
	} else {
		c.JSON(http.StatusNotAcceptable, struct{ message string }{message: "Don't belong to the user"})
	}
}
