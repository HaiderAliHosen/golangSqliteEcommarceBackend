package services

import (
	"github.com/HaiderAliHosen/sqlitedemo/infrastructure"
	"github.com/HaiderAliHosen/sqlitedemo/models"
)

//FetchOrdersPage ---
func FetchOrdersPage(userID uint, page, pageSize int) (orders []models.Order, totalOrdersCount int, err error) {
	database := infrastructure.GetDb()

	totalOrdersCount = 0

	query := database.Model(&models.Order{}).Where(&models.Order{UserID: userID})
	query.Count(&totalOrdersCount)

	err = query.Offset((page - 1) * pageSize).Limit(pageSize).
		// TODO: Why Preload("Address") does not work?, perhaps OrderItems does
		// Preload("OrderItems").Preload("Address").
		Find(&orders).Error
	if err != nil {
		return
	}

	var orderIds = make([]uint, len(orders))
	for i := 0; i < len(orders); i++ {
		orderIds[i] = orders[i].ID
	}

	var orderItems []models.OrderItem
	if len(orders) > 0 {
		//
		database.Select("id, order_id").Where("order_id in (?)", orderIds).Find(&orderItems)

		for i := 0; i < len(orderItems); i++ {
			oi := orderItems[i]
			for j := 0; j < len(orders); j++ {
				if oi.OrderID == orders[j].ID {
					orders[j].OrderItemsCount = orders[j].OrderItemsCount + 1
				}
			}
		}
	}
	return orders, totalOrdersCount, err
}

//FetchOrderDetails ---
func FetchOrderDetails(orderID uint) (order models.Order, err error) {
	database := infrastructure.GetDb()
	err = database.Model(models.Order{}).Preload("OrderItems").First(&order, orderID).Error
	var address models.Address
	database.Model(&order).Related(&address)
	order.Address = address
	return order, err
}
