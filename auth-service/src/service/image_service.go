package service

import (
	"fmt"

	"src/model"
)

func (u *UserService) GetImage(id int) (picture *model.Picture) {
	u.DB.InitDB().First(&picture, id)
	return
}

func (u *UserService) GetImageByUser(user_id int) (picture *model.Picture) {
	if result := u.DB.InitDB().Where("id = (?)", u.DB.InitDB().Select("pictureId").Where("id = ?", user_id).Table("users")).Find(&picture); result.Error != nil {
		picture = nil
		fmt.Println(result)
	}
	return
}

func (u *UserService) SaveImage(image []byte, contentType string, fileName string) (pictureId int, result bool) {
	picture := model.Picture{
		Picture:     &image,
		ContentType: &contentType,
		FileName:    &fileName,
	}

	DB_result := u.DB.InitDB().Create(&picture)
	if DB_result.Error == nil {
		pictureId = picture.Id
		result = true
	}
	return
}

func (u *UserService) SaveImageURL(image_url string) (pictureId int, result bool) {
	picture := model.Picture{
		PictureUrl: &image_url,
	}

	DB_result := u.DB.InitDB().Create(picture)
	if DB_result.Error == nil {
		pictureId = picture.Id
		result = true
	}
	return
}

func (u *UserService) UpdateImage(id int, image []byte, contentType string, fileName string) (result bool) {
	DB_result := u.DB.InitDB().Save(&model.Picture{
		Id:          id,
		Picture:     &image,
		ContentType: &contentType,
		FileName:    &fileName,
	})
	if DB_result.Error == nil {
		result = true
	}
	return
}
