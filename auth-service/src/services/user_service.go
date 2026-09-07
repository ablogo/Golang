package services

import (
	"src/db"
	"src/models"
)

func CreateUser(new_user models.SignUp) (result bool) {
	db := db.InitDB()
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
	db := db.InitDB()

	for _, value := range dependencies {
		db = db.Preload(value)
	}
	db.Find(&user, id)
	return
}

func GetUserByEmail(email string) (user *models.User) {
	db := db.InitDB()

	if result := db.First(&user, "email = ?", email); result.Error != nil {
		user = nil
	}
	return
}

func GetUsers() (users []models.User) {
	db := db.InitDB()
	users = []models.User{}

	result := db.Find(&users)

	if result.Error != nil {
		users = nil
	}
	return
}

func DeleteUser(id int) bool {
	db := db.InitDB()
	result := db.Delete(&models.User{}, &id)

	if result.RowsAffected == 1 {
		return true
	} else {
		return false
	}
}

func UpdateUser(updated_user models.User) bool {
	db := db.InitDB()

	result := db.Model(&updated_user).Updates(models.User{
		Name:      updated_user.Name,
		LastName:  updated_user.LastName,
		PictureId: updated_user.PictureId,
	})

	if result.RowsAffected == 1 {
		return true
	} else {
		return false
	}
}

func ChangePassword(user *models.User, password string) (result bool) {

	db := db.InitDB()

	if user != nil {
		r := db.Model(&user).Update("password", hashPassword(password))

		if r.Error == nil {
			result = true
		}
	}
	return
}

func DisableUser(user *models.User, value bool) bool {
	db := db.InitDB()

	result := db.Model(&user).Updates(map[string]any{"Disabled": value})

	if result.RowsAffected == 1 {
		return true
	} else {
		return false
	}
}

func AddPicture(user *models.User, image []byte, content_type string, file_name string) (result bool) {

	if user.PictureId == nil {
		picture_id, result := SaveImage(image, content_type, file_name)
		if result {
			user.PictureId = &picture_id
			result = UpdateUser(*user)
		}
	} else {
		result = UpdateImage(*user.PictureId, image, content_type, file_name)
	}
	return
}

func AddAddress(user *models.User, model models.Address) (address *models.Address, result bool) {
	model.UserId = user.Id
	address, result = CreateAddress(model)
	return
}
