package services

import (
	"src/models"
)

func (u *UserService) CreateAddress(model models.Address) (*models.Address, bool) {
	db_result := u.Postgres.InitDB().Create(&model)
	if db_result.Error != nil {
		return nil, false
	}

	return &model, true
}

func (u *UserService) GetAddress(id int) (address *models.Address) {
	u.Postgres.InitDB().First(&address, id)
	return
}

func (u *UserService) GetAddressByUser(user_id int) (address *[]models.Address) {
	u.Postgres.InitDB().Where("user_id = ?", user_id).Find(&address)
	return
}

func (u *UserService) UpdateAddress(model models.Address) (result bool) {
	r := u.Postgres.InitDB().Save(&model)
	if r.Error == nil {
		result = true
	}
	return
}

func (u *UserService) DeleteAddress(id int) bool {
	r := u.Postgres.InitDB().Delete(models.Address{}, id)

	if r.Error != nil {
		return false
	}

	return true
}
