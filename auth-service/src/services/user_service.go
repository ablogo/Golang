package services

import (
	database "src/db"
	"src/models"
)

type UserService struct {
	Postgres *database.DBHandler
}

func (u *UserService) CreateUser(new_user models.SignUp) (result bool) {
	user := models.User{
		Name:     new_user.Name,
		LastName: new_user.LastName,
		Email:    new_user.Email,
		Password: hashPassword(new_user.Password),
	}

	db_result := u.Postgres.InitDB().Create(&user)

	if db_result.Error == nil {
		result = true
	}
	return
}

func (u *UserService) GetUser(id int, dependencies []string) (user *models.User) {
	// To avoid repeating over and over "AND "users"."id" = 1", a new variable is used,
	// because GORM mutates the query state due to the "db" is a global variable.
	query := u.Postgres.InitDB()

	for _, value := range dependencies {
		query = u.Postgres.InitDB().Preload(value)
	}
	err := query.Find(&user, id).Error
	if err != nil {
		return nil
	}
	return
}

func (u *UserService) GetUserByEmail(email string) (user *models.User) {
	if result := u.Postgres.InitDB().First(&user, "email = ?", email); result.Error != nil {
		user = nil
	}
	return
}

func (u *UserService) GetUsers() (users []models.User) {
	users = []models.User{}

	result := u.Postgres.InitDB().Find(&users)

	if result.Error != nil {
		users = nil
	}
	return
}

func (u *UserService) DeleteUser(id int) bool {
	result := u.Postgres.InitDB().Delete(&models.User{}, &id)

	if result.RowsAffected == 1 {
		return true
	} else {
		return false
	}
}

func (u *UserService) UpdateUser(updated_user models.User) bool {
	result := u.Postgres.InitDB().Model(&updated_user).Updates(models.User{
		Name:     updated_user.Name,
		LastName: updated_user.LastName,
	})

	if result.RowsAffected == 1 {
		return true
	} else {
		return false
	}
}

func (u *UserService) UpdateUserImage(updated_user models.User, pictureId int) bool {
	result := u.Postgres.InitDB().Model(&updated_user).Updates(models.User{
		PictureId: &pictureId,
	})

	if result.RowsAffected == 1 {
		return true
	} else {
		return false
	}
}

func (u *UserService) ChangePassword(user *models.User, password string) (result bool) {
	if user != nil {
		r := u.Postgres.InitDB().Model(&user).Update("password", hashPassword(password))

		if r.Error == nil {
			result = true
		}
	}
	return
}

func (u *UserService) DisableUser(user *models.User, value bool) bool {
	result := u.Postgres.InitDB().Model(&user).Updates(map[string]any{"Disabled": value})

	if result.RowsAffected == 1 {
		return true
	} else {
		return false
	}
}

func (u *UserService) AddPicture(user *models.User, image []byte, contentType string, fileName string) (result bool) {
	/*if user.PictureId == nil {
		picture_id, result := SaveImage(image, contentType, fileName)
		if result {
			result = u.UpdateUserImage(*user, picture_id)
		}
	} else {
		result = UpdateImage(*user.PictureId, image, contentType, fileName)
	}*/
	return
}

func (u *UserService) AddAddress(user *models.User, model models.Address) (address *models.Address, result bool) {
	model.UserId = user.Id
	//address, result = CreateAddress(model)
	return
}
