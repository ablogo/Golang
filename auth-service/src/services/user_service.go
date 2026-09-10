package services

import (
	database "src/db"
	"src/models"
	"src/utils"
)

var settings = &utils.Settings{}
var dbHandler = &database.DBHandler{Settings: settings.GetInstance()}
var db = dbHandler.InitDB()

func CreateUser(new_user models.SignUp) (result bool) {
	user := models.User{
		Name:     new_user.Name,
		LastName: new_user.LastName,
		Email:    new_user.Email,
		Password: hashPassword(new_user.Password),
	}

	db_result := db.Create(&user)

	if db_result.Error == nil {
		result = true
	}
	return
}

func GetUser(id int, dependencies []string) (user *models.User) {
	// To avoid repeating over and over "AND "users"."id" = 1", a new variable is used,
	// because GORM mutates the query state due to the "db" is a global variable.
	query := db

	for _, value := range dependencies {
		query = db.Preload(value)
	}
	err := query.Find(&user, id).Error
	if err != nil {
		return nil
	}
	return
}

func GetUserByEmail(email string) (user *models.User) {
	if result := db.First(&user, "email = ?", email); result.Error != nil {
		user = nil
	}
	return
}

func GetUsers() (users []models.User) {
	users = []models.User{}

	result := db.Find(&users)

	if result.Error != nil {
		users = nil
	}
	return
}

func DeleteUser(id int) bool {
	result := db.Delete(&models.User{}, &id)

	if result.RowsAffected == 1 {
		return true
	} else {
		return false
	}
}

func UpdateUser(updated_user models.User) bool {
	result := db.Model(&updated_user).Updates(models.User{
		Name:     updated_user.Name,
		LastName: updated_user.LastName,
	})

	if result.RowsAffected == 1 {
		return true
	} else {
		return false
	}
}

func UpdateUserImage(updated_user models.User, pictureId int) bool {
	result := db.Model(&updated_user).Updates(models.User{
		PictureId: &pictureId,
	})

	if result.RowsAffected == 1 {
		return true
	} else {
		return false
	}
}

func ChangePassword(user *models.User, password string) (result bool) {
	if user != nil {
		r := db.Model(&user).Update("password", hashPassword(password))

		if r.Error == nil {
			result = true
		}
	}
	return
}

func DisableUser(user *models.User, value bool) bool {
	result := db.Model(&user).Updates(map[string]any{"Disabled": value})

	if result.RowsAffected == 1 {
		return true
	} else {
		return false
	}
}

func AddPicture(user *models.User, image []byte, contentType string, fileName string) (result bool) {
	if user.PictureId == nil {
		picture_id, result := SaveImage(image, contentType, fileName)
		if result {
			result = UpdateUserImage(*user, picture_id)
		}
	} else {
		result = UpdateImage(*user.PictureId, image, contentType, fileName)
	}
	return
}

func AddAddress(user *models.User, model models.Address) (address *models.Address, result bool) {
	model.UserId = user.Id
	address, result = CreateAddress(model)
	return
}
