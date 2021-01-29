package services

import (
	"github.com/HaiderAliHosen/sqlitedemo/infrastructure"
	"github.com/HaiderAliHosen/sqlitedemo/models"
)

//FetchAllTags ---
func FetchAllTags() ([]models.Tag, error) {
	database := infrastructure.GetDb()
	var tags []models.Tag
	err := database.Preload("Images", "tag_id IS NOT NULL").Find(&tags).Error
	return tags, err
}
