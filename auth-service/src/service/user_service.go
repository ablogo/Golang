package service

import (
	"src/model"
	"src/postgresql"
)

type UserService struct {
	DB *postgresql.DBHandler
}

func (u *UserService) CreateUser(new_user model.SignUp) (result bool) {
	user := model.User{
		Name:     new_user.Name,
		LastName: new_user.LastName,
		Email:    new_user.Email,
		Password: hashPassword(new_user.Password),
	}

	db_result := u.DB.InitDB().Create(&user)

	if db_result.Error == nil {
		result = true
	}
	return
}

func (u *UserService) GetUser(id int, dependencies []string) (user *model.User) {
	// To avoid repeating over and over "AND "users"."id" = 1", a new variable is used,
	// because GORM mutates the query state due to the "DB" is a global variable.
	query := u.DB.InitDB()

	for _, value := range dependencies {
		query = u.DB.InitDB().Preload(value)
	}
	err := query.Find(&user, id).Error
	if err != nil {
		return nil
	}
	return
}

func (u *UserService) GetUserByEmail(email string) (user *model.User) {
	if result := u.DB.InitDB().First(&user, "email = ?", email); result.Error != nil {
		user = nil
	}
	return
}

func (u *UserService) GetUsers() (users []model.User) {
	users = []model.User{}

	result := u.DB.InitDB().Find(&users)

	if result.Error != nil {
		users = nil
	}
	return
}

func (u *UserService) DeleteUser(id int) bool {
	result := u.DB.InitDB().Delete(&model.User{}, &id)

	if result.RowsAffected == 1 {
		return true
	} else {
		return false
	}
}

func (u *UserService) UpdateUser(updated_user model.User) bool {
	result := u.DB.InitDB().Model(&updated_user).Updates(model.User{
		Name:     updated_user.Name,
		LastName: updated_user.LastName,
	})

	if result.RowsAffected == 1 {
		return true
	} else {
		return false
	}
}

func (u *UserService) UpdateUserImage(updated_user model.User, pictureId int) bool {
	result := u.DB.InitDB().Model(&updated_user).Updates(model.User{
		PictureId: &pictureId,
	})

	if result.RowsAffected == 1 {
		return true
	} else {
		return false
	}
}

func (u *UserService) ChangePassword(user *model.User, password string) (result bool) {
	if user != nil {
		r := u.DB.InitDB().Model(&user).Update("password", hashPassword(password))

		if r.Error == nil {
			result = true
		}
	}
	return
}

func (u *UserService) DisableUser(user *model.User, value bool) bool {
	result := u.DB.InitDB().Model(&user).Updates(map[string]any{"Disabled": value})

	if result.RowsAffected == 1 {
		return true
	} else {
		return false
	}
}

func (u *UserService) AddPicture(user *model.User, image []byte, contentType string, fileName string) (result bool) {
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

func (u *UserService) AddAddress(user *model.User, model model.Address) (address *model.Address, result bool) {
	model.UserId = user.Id
	//address, result = CreateAddress(model)
	return
}
