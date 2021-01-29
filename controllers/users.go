package controllers

import (
	"errors"
	"fmt"

	"github.com/HaiderAliHosen/sqlitedemo/dtos"
	"github.com/HaiderAliHosen/sqlitedemo/models"
	"github.com/HaiderAliHosen/sqlitedemo/services"
	"github.com/gin-gonic/gin"

	"net/http"

	"golang.org/x/crypto/bcrypt"
)

//RegisterUserRoutes ---
func RegisterUserRoutes(router *gin.RouterGroup) {
	router.POST("/", UsersRegistration)
	router.POST("/login", UsersLogin)
}

//UsersRegistration ----
func UsersRegistration(c *gin.Context) {

	var json dtos.RegisterRequestDto
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)
	if err := services.CreateOne(&models.User{
		Username:  json.Username,
		Password:  string(password),
		FirstName: json.FirstName,
		LastName:  json.LastName,
		Email:     json.Email,
	}); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dtos.CreateDetailedErrorDto("database", err))
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success":       true,
		"full_messages": []string{"User created successfully"}})
}

//UsersLogin ---
func UsersLogin(c *gin.Context) {

	var json dtos.LoginRequestDto
	fmt.Println(json)
	if err := c.ShouldBindJSON(&json); err != nil {
		fmt.Println(" can't ShouldBindJSON")
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}

	user, err := services.FindOneUser(&models.User{Username: json.Username})

	if err != nil {
		c.JSON(http.StatusForbidden, dtos.CreateDetailedErrorDto("login_error", err))
		fmt.Println("login_error CreateDetailedErrorDto")
		return
	}

	if user.IsValidPassword(json.Password) != nil {
		c.JSON(http.StatusForbidden, dtos.CreateDetailedErrorDto("login", errors.New("invalid credentials")))
		fmt.Println("invalid credentials password")
		return
	}

	c.JSON(http.StatusOK, dtos.CreateLoginSuccessful(&user))

}
