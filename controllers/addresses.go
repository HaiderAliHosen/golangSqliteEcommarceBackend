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

//RegisterAddressesRoutes ---
func RegisterAddressesRoutes(router *gin.RouterGroup) {

	router.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		router.GET("/addresses", ListAddresses)
		router.POST("/addresses", CreateAddress)
	}

}

//ListAddresses ---
func ListAddresses(c *gin.Context) {

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

	// userId:= c.Keys["currentUserId"].(uint) // or
	userID := c.MustGet("currentUserID").(uint)
	includeUser := false
	addresses, totalCommentCount := services.FetchAddressesPage(userID, page, pageSize, includeUser)

	c.JSON(http.StatusOK, dtos.CreateAddressPagedResponse(c.Request, addresses, page, pageSize, totalCommentCount, includeUser))
}

//CreateAddress ---
func CreateAddress(c *gin.Context) {

	user := c.MustGet("currentUser").(models.User)

	var json dtos.CreateAddress
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}
	firstName := json.FirstName
	lastName := json.LastName
	if firstName == "" {
		firstName = user.FirstName
	}
	if lastName == "" {
		lastName = user.LastName
	}
	address := models.Address{
		FirstName:     firstName,
		LastName:      lastName,
		Country:       json.Country,
		City:          json.City,
		StreetAddress: json.StreetAddress,
		ZipCode:       json.ZipCode,
		User:          user,
		UserID:        user.ID,
	}

	if err := services.SaveOne(&address); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dtos.CreateDetailedErrorDto("database_error", err))
		return
	}

	c.JSON(http.StatusOK, dtos.GetAddressCreatedDto(&address, false))
}
