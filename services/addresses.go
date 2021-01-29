package services

import (
	"github.com/HaiderAliHosen/sqlitedemo/infrastructure"
	"github.com/HaiderAliHosen/sqlitedemo/models"
)

//FetchAddressesPage ---
func FetchAddressesPage(userID uint, page, pageSize int, includeUser bool) ([]models.Address, int) {
	var addresses []models.Address
	var totalAddressesCount int
	database := infrastructure.GetDb()
	database.Model(&models.Address{}).Where(&models.Address{UserID: uint(userID)}).Count(&totalAddressesCount)
	database.Where(&models.Address{UserID: uint(userID)}).
		Offset((page - 1) * pageSize).Limit(pageSize).
		Preload("User").
		Find(&addresses)

	if includeUser {
		var userIds = make([]uint, len(addresses))
		var users []models.User
		for i := 0; i < len(addresses); i++ {
			userIds[i] = addresses[i].UserID
		}
		database.Select([]string{"id", "username"}).Where(userIds).Find(&users)

		// If the user gets deleted and the comment is still in the database we may have less users than addresses
		// Another scenario (the one I run into) is there is a problem with the Comment.User, the Comment.UserId does not get saved automatically
		for i := 0; i < len(addresses); i++ {
			address := addresses[i]
			for j := 0; j < len(users); j++ {
				user := users[j]
				if address.UserID == user.ID {
					addresses[i].User = users[j]
				}
			}
		}
	}
	return addresses, totalAddressesCount
}

//FetchAddress ---
func FetchAddress(addressID uint) (address models.Address) {
	database := infrastructure.GetDb()
	database.First(&address, addressID)
	return address
}

//FetchIdsFromAddress --
func FetchIdsFromAddress(addressID uint) (address models.Address) {
	database := infrastructure.GetDb()
	database.Select("id, user_id").First(&address, addressID)
	return
}
