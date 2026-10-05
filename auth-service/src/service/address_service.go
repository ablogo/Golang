package service

import "src/model"

func (u *UserService) CreateAddress(model model.Address) (*model.Address, bool) {
	db_result := u.DB.InitDB().Create(&model)
	if db_result.Error != nil {
		return nil, false
	}

	return &model, true
}

func (u *UserService) GetAddress(id int) (address *model.Address) {
	u.DB.InitDB().First(&address, id)
	return
}

func (u *UserService) GetAddressByUser(user_id int) (address *[]model.Address) {
	u.DB.InitDB().Where("user_id = ?", user_id).Find(&address)
	return
}

func (u *UserService) UpdateAddress(model model.Address) (result bool) {
	r := u.DB.InitDB().Save(&model)
	if r.Error == nil {
		result = true
	}
	return
}

func (u *UserService) DeleteAddress(id int) bool {
	r := u.DB.InitDB().Delete(model.Address{}, id)

	if r.Error != nil {
		return false
	}

	return true
}
