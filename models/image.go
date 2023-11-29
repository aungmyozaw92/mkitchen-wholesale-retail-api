package models

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/myanmarmarathon/mkitchen-distribution-backend/helper"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/utils"
)

type Image struct {
	ID         uint   
    ImageUrl   string `json:"image_url"`
    OwnerType  string `json:"owner_type"`
    OwnerID    uint  `json:"owner_id"`
}

func transformImageURLs(images []Image) []Image {
    urlAddress := "https://" + os.Getenv("SP_BUCKET") + "." + os.Getenv("SP_URL") 

    for i := range images {
		url := urlAddress + "/" + images[i].ImageUrl
        images[i].ImageUrl = url
    }

	return images
}

func (input *Image) UploadImage() (*Image, error) {

	if input.OwnerType != "products" && input.OwnerType != "product_variations"{
		return &Image{}, errors.New("invalid image owner type")
	}

	if input.OwnerType == "products" {
		isValidId := helper.IsRecordValidByID(input.OwnerID, &Product{}, DB)

		if !isValidId {
			return &Image{}, errors.New("invalid product id")
		}
	}

	if input.OwnerType == "product_variations" {
		isValidId := helper.IsRecordValidByID(input.OwnerID, &ProductVariation{}, DB)

		if !isValidId {
			return &Image{}, errors.New("invalid product variation id")
		}
	}
	

    storagePath := "products/"
	uniqueFilename := helper.GenerateUniqueFilename()
	imageFilePath := filepath.Join(storagePath, uniqueFilename)
	imageObjectUrl := storagePath + uniqueFilename

    err := utils.SaveImageToSpaces(imageObjectUrl,  input.ImageUrl)

	if err != nil {
		return &Image{}, errors.New("failed upload to cloud space")
	}
    input.ImageUrl = imageFilePath

    err = DB.Create(&input).Error
	
	if err != nil {
		return &Image{}, err
	}
	return input, nil
}

func (input *Image) DeleteImage(id uint64) (*Image, error) {

	err := DB.Model(&Image{}).Where("id = ?", id).First(&input).Error
	if  err != nil {
        return nil, helper.ErrorRecordNotFound
    }

    if err := utils.DeleteImageFromSpaces(input.ImageUrl); err != nil {
        return &Image{}, err
    }

	err = DB.Delete(&input).Error
	if err != nil {
		return &Image{}, err
	}
	return input, nil
}