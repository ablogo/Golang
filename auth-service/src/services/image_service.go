package services

import (
	"fmt"

	"src/models"
)

func GetImage(id int) (picture *models.Picture) {
	db.First(&picture, id)
	return
}

func GetImageByUser(user_id int) (picture *models.Picture) {
	if result := db.Where("id = (?)", db.Select("pictureId").Where("id = ?", user_id).Table("users")).Find(&picture); result.Error != nil {
		picture = nil
		fmt.Println(result)
	}
	return
}

func SaveImage(image []byte, contentType string, fileName string) (pictureId int, result bool) {
	picture := models.Picture{
		Picture:     &image,
		ContentType: &contentType,
		FileName:    &fileName,
	}

	db_result := db.Create(&picture)
	if db_result.Error == nil {
		pictureId = picture.Id
		result = true
	}
	return
}

func SaveImageURL(image_url string) (pictureId int, result bool) {
	picture := models.Picture{
		PictureUrl: &image_url,
	}

	db_result := db.Create(picture)
	if db_result.Error == nil {
		pictureId = picture.Id
		result = true
	}
	return
}

func UpdateImage(id int, image []byte, contentType string, fileName string) (result bool) {
	db_result := db.Save(&models.Picture{
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
