package services

import "github.com/HaiderAliHosen/sqlitedemo/infrastructure"

//CreateOne ---
func CreateOne(data interface{}) error {
	database := infrastructure.GetDb()
	err := database.Create(data).Error
	return err
}

//SaveOne ---
func SaveOne(data interface{}) error {
	database := infrastructure.GetDb()
	err := database.Save(data).Error
	return err
}
