package models

type Tag struct {
    ID     		uint       `gorm:"primary_key" json:"id"`
    Name   		string	   `gorm:"size:50;not null:unique" json:"name"`
    Products 	[]Product  `gorm:"many2many:product_tags;"`
}

type ProductTags struct {
    ProductID uint
    TagID     uint
}

func CreateOrUpdateTags(tags []Tag) ([]Tag, error) {
    for i := range tags {
        var existingTag Tag
        if err := DB.Where("name = ?", tags[i].Name).First(&existingTag).Error; err != nil {
            // Tag doesn't exist, so create it
            newTag := Tag{Name: tags[i].Name}
            if err := DB.Create(&newTag).Error; err != nil {
                return nil, err
            }
            tags[i] = newTag
        } else {
            // Tag already exists, associate it
            tags[i] = existingTag
        }
    }
    return tags, nil
}