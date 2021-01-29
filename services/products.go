package services

import (
	"github.com/HaiderAliHosen/sqlitedemo/infrastructure"
	"github.com/HaiderAliHosen/sqlitedemo/models"
)

//FetchProductsPage ---
func FetchProductsPage(page int, pageSize int) ([]models.Product, int, []int, error) {
	database := infrastructure.GetDb()
	var products []models.Product
	var count int
	tx := database.Begin()
	database.Model(&products).Count(&count)
	database.Offset((page - 1) * pageSize).Limit(pageSize).Find(&products)
	tx.Model(&products).
		Preload("Tags").Preload("Categories").Preload("Images").
		Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&products)
	commentsCount := make([]int, len(products))

	for index, product := range products {
		commentsCount[index] = tx.Model(&product).Association("Comments").Count()
	}
	err := tx.Commit().Error
	return products, count, commentsCount, err
}

//FetchProductDetails ---
func FetchProductDetails(condition interface{}, optional ...bool) models.Product {
	database := infrastructure.GetDb()
	var product models.Product

	query := database.Where(condition).
		Preload("Tags").Preload("Categories").Preload("Images").Preload("Comments")
	// Unfortunately .Preload("Comments.User") does not work as the doc states ...
	query.First(&product)
	includeUserComment := false

	if len(optional) > 0 {
		includeUserComment = optional[0]
	}

	if includeUserComment {

		for i := 0; i < len(product.Comments); i++ {
			database.Model(&product.Comments[i]).Related(&product.Comments[i].User, "UserId")
		}

		var userIds = make([]uint, len(product.Comments))
		var users []models.User
		for i := 0; i < len(product.Comments); i++ {
			userIds[i] = product.Comments[i].UserID
		}
		// WHERE users.id IN userIds; This will also work: Select([]string{"id", "username"})
		database.Select("id, username").Where(userIds).Find(&users)

		for i := 0; i < len(product.Comments); i++ {
			user := users[i]
			comment := product.Comments[i]
			if comment.UserID == user.ID {
				product.Comments[i].User = users[i]
			}
		}
	}

	return product
}

//FetchProductID ---
func FetchProductID(slug string) (uint, error) {
	productID := -1
	database := infrastructure.GetDb()
	err := database.Model(&models.Product{}).Where(&models.Product{Slug: slug}).Select("id").Row().Scan(&productID)
	return uint(productID), err
}

//SetTags ---
func SetTags(product *models.Product, tags []string) error {
	database := infrastructure.GetDb()
	var tagList []models.Tag
	for _, tag := range tags {
		var tagModel models.Tag
		err := database.FirstOrCreate(&tagModel, models.Tag{Name: tag}).Error
		if err != nil {
			return err
		}
		tagList = append(tagList, tagModel)
	}
	product.Tags = tagList
	return nil
}

//Update ---
func Update(product *models.Product, data interface{}) error {
	database := infrastructure.GetDb()
	err := database.Model(product).Update(data).Error
	return err
}

//DeleteProduct ---
func DeleteProduct(condition interface{}) error {
	db := infrastructure.GetDb()
	err := db.Where(condition).Delete(models.Product{}).Error
	return err
}

//FetchProductsIDNameAndPrice ---
func FetchProductsIDNameAndPrice(productIds []uint) (products []models.Product, err error) {
	database := infrastructure.GetDb()
	err = database.Select([]string{"id", "name", "slug", "price"}).Find(&products, productIds).Error
	return products, err
}
