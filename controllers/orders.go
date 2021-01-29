package controllers

import (
	"github.com/HaiderAliHosen/sqlitedemo/dtos"
	"github.com/HaiderAliHosen/sqlitedemo/middlewares"
	"github.com/HaiderAliHosen/sqlitedemo/models"
	"github.com/HaiderAliHosen/sqlitedemo/services"
	"github.com/gin-gonic/gin"

	"net/http"
	"strconv"
)

//RegisterOrderRoutes ---
func RegisterOrderRoutes(router *gin.RouterGroup) {
	router.POST("", CreateOrder)
	router.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		router.GET("", ListOrders)
		router.GET("/:id", ShowOrder)
	}
}

//ListOrders --
func ListOrders(c *gin.Context) {
	pageSizeStr := c.Query("page_size")
	pageStr := c.Query("page")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 5
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	userID := c.MustGet("currentUserId").(uint)

	orders, totalCommentCount, err := services.FetchOrdersPage(userID, page, pageSize)

	c.JSON(http.StatusOK, dtos.CreateOrderPagedResponse(c.Request, orders, page, pageSize, totalCommentCount, false, false))
}

//ShowOrder ---
func ShowOrder(c *gin.Context) {
	orderID, err := strconv.Atoi(c.Param("id"))
	user := c.MustGet("currentUser").(models.User)
	order, err := services.FetchOrderDetails(uint(orderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dtos.CreateDetailedErrorDto("db_error", err))
		return
	}

	if order.UserID == user.ID || user.IsAdmin() {
		c.JSON(http.StatusOK, dtos.CreateOrderDetailsDto(&order))
	} else {
		c.JSON(http.StatusForbidden, dtos.CreateErrorDtoWithMessage("Permission denied, you can not view this order"))
		return
	}
}

//CreateOrder ---
func CreateOrder(c *gin.Context) {
	var orderRequest dtos.CreateOrderRequestDto
	if err := c.ShouldBind(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}

	userObj, userLoggedIn := c.Get("currentUser")
	var user models.User
	if userLoggedIn {
		user = (userObj).(models.User)
	}

	var address models.Address
	// Reuse address can only be done by authenticated users
	if orderRequest.AddressID != 0 && userLoggedIn {
		address = services.FetchAddress(orderRequest.AddressID)
		/*if err != nil || address.ID == 0 {
			c.JSON(http.StatusBadRequest, dtos.CreateDetailedErrorDto("db_error", err))
			return
		}*/
		if address.UserID != user.ID {
			c.JSON(http.StatusForbidden, dtos.CreateErrorDtoWithMessage("permission denied"))
			return
		}
	} else if orderRequest.AddressID == 0 {
		address = models.Address{
			FirstName:     orderRequest.FirstName,
			LastName:      orderRequest.LastName,
			City:          orderRequest.City,
			Country:       orderRequest.Country,
			StreetAddress: orderRequest.StreetAddress,
			ZipCode:       orderRequest.ZipCode,
		}
		if userLoggedIn {
			address.UserID = user.ID
		}
		err := services.CreateOne(&address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

	} else {
		c.JSON(http.StatusForbidden, dtos.CreateErrorDtoWithMessage("Operation not supported, what are you trying to do?"))
		return
	}

	order := models.Order{
		TrackingNumber: randomString(16),
		OrderStatus:    0,
		Address:        address,
		AddressID:      address.ID,
	}

	if userLoggedIn {
		order.UserID = user.ID
		order.User = user
	}

	var productIds = make([]uint, len(orderRequest.CartItems))
	for i := 0; i < len(orderRequest.CartItems); i++ {
		productIds[i] = orderRequest.CartItems[i].ID
	}

	products, err := services.FetchProductsIDNameAndPrice(productIds)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, dtos.CreateDetailedErrorDto("db_error", err))
		return
	}

	if len(products) != len(orderRequest.CartItems) {
		c.JSON(http.StatusUnprocessableEntity, dtos.CreateErrorDtoWithMessage("make sure all products are still available"))
		return
	}
	orderItems := make([]models.OrderItem, len(products))

	for i := 0; i < len(products); i++ {
		// I am assuming product ids returned are in the same order as the cart_items, TODO: implement a more robust code to ensure
		orderItems[i] = models.OrderItem{
			ProductID:   products[i].ID,
			ProductName: products[i].Name,
			Slug:        products[i].Slug,
			Quantity:    orderRequest.CartItems[i].Quantity,
		}
	}

	order.OrderItems = orderItems
	err = services.CreateOne(&order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, dtos.CreateOrderCreatedDto(&order))

}
