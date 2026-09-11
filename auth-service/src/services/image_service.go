package services

import (
	"fmt"

	"src/models"
)

func (u *UserService) GetImage(id int) (picture *models.Picture) {
	u.Postgres.InitDB().First(&picture, id)
	return
}

func (u *UserService) GetImageByUser(user_id int) (picture *models.Picture) {
	if result := u.Postgres.InitDB().Where("id = (?)", u.Postgres.InitDB().Select("pictureId").Where("id = ?", user_id).Table("users")).Find(&picture); result.Error != nil {
		picture = nil
		fmt.Println(result)
	}
	return
}

func (u *UserService) SaveImage(image []byte, contentType string, fileName string) (pictureId int, result bool) {
	picture := models.Picture{
		Picture:     &image,
		ContentType: &contentType,
		FileName:    &fileName,
	}

	db_result := u.Postgres.InitDB().Create(&picture)
	if db_result.Error == nil {
		pictureId = picture.Id
		result = true
	}
	return
}

func (u *UserService) SaveImageURL(image_url string) (pictureId int, result bool) {
	picture := models.Picture{
		PictureUrl: &image_url,
	}

	db_result := u.Postgres.InitDB().Create(picture)
	if db_result.Error == nil {
		pictureId = picture.Id
		result = true
	}
	return
}

func (u *UserService) UpdateImage(id int, image []byte, contentType string, fileName string) (result bool) {
	db_result := u.Postgres.InitDB().Save(&models.Picture{
		Id:          id,
		Picture:     &image,
		ContentType: &contentType,
		FileName:    &fileName,
	})
	if db_result.Error == nil {
		result = true
	}
	return
}
