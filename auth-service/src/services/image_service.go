package services

import (
	"fmt"

	"src/db"
	"src/models"
)

func GetImage(id int) (picture *models.Picture) {
	db := db.InitDB()

	db.First(&picture, id)
	return
}

func GetImageByUser(user_id int) (picture *models.Picture) {
	db := db.InitDB()

	if result := db.Where("id = (?)", db.Select("picture_id").Where("id = ?", user_id).Table("users")).Find(&picture); result.Error != nil {
		picture = nil
		fmt.Println(result)
	}
	return
}

func SaveImage(image []byte, content_type string, file_name string) (picture_id int, result bool) {
	db := db.InitDB()

	picture := models.Picture{
		Picture:     &image,
		ContentType: &content_type,
		FileName:    &file_name,
	}

	db_result := db.Create(&picture)
	if db_result.Error == nil {
		picture_id = picture.Id
		result = true
	}
	return
}

func SaveImageURL(image_url string) (picture_id int, result bool) {
	db := db.InitDB()

	picture := models.Picture{
		PictureUrl: &image_url,
	}

	db_result := db.Create(picture)
	if db_result.Error == nil {
		picture_id = picture.Id
		result = true
	}
	return
}

func UpdateImage(id int, image []byte, content_type string, file_name string) (result bool) {
	db := db.InitDB()

	db_result := db.Save(&models.Picture{
		Id:          id,
		Picture:     &image,
		ContentType: &content_type,
		FileName:    &file_name,
	})
	if db_result.Error == nil {
		result = true
	}
	return
}
