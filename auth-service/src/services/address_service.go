package services

import (
	"src/models"
)

func CreateAddress(model models.Address) (*models.Address, bool) {
	db_result := db.Create(&model)
	if db_result.Error != nil {
		return nil, false
	}

	return &model, true
}

func GetAddress(id int) (address *models.Address) {
	db.First(&address, id)
	return
}

func GetAddressByUser(user_id int) (address *[]models.Address) {
	db.Where("user_id = ?", user_id).Find(&address)
	return
}

func UpdateAddress(model models.Address) (result bool) {
	r := db.Save(&model)
	if r.Error == nil {
		result = true
	}
	return
}

func DeleteAddress(id int) bool {
	r := db.Delete(models.Address{}, id)

	if r.Error != nil {
		return false
	}

	return true
}
