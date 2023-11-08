package models

import "os"

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