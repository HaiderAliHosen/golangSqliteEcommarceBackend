package dtos

import "github.com/HaiderAliHosen/sqlitedemo/models"

//CreateHomeResponse --
func CreateHomeResponse(tags []models.Tag, categories []models.Category) map[string]interface{} {
	return CreateSuccessDto(map[string]interface{}{
		"tags":       CreateTagListDto(tags),
		"categories": CreateCategoryListDto(categories),
	})
}
