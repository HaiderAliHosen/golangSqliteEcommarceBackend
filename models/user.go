package models

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//User --
type User struct {
	gorm.Model
	//Id           uint    `gorm:"primary_key"`
	FirstName string `gorm:"varchar(255);not null"`
	LastName  string `gorm:"varchar(255);not null"`
	Username  string `gorm:"column:username"`
	Email     string `gorm:"column:email;unique_index"`
	Password  string `gorm:"column:password;not null"`

	Comments []Comment `gorm:"foreignkey:UserId"`

	Roles     []Role     `gorm:"many2many:users_roles;"`
	UserRoles []UserRole `gorm:"foreignkey:UserId"`
}

/*SetPassword .... What's bcrypt? https://en.wikipedia.org/wiki/Bcrypt
Golang bcrypt doc: https://godoc.org/golang.org/x/crypto/bcrypt
You can change the value in bcrypt.DefaultCost to adjust the security index.
	err := userModel.setPassword("password0")*/
func (user *User) SetPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty")
	}
	bytePassword := []byte(password)
	// Make sure the second param `bcrypt generator cost` between [4, 32)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	user.Password = string(passwordHash)
	return nil
}

/*IsValidPassword --- Database will only save the hashed string, you should check it by util function.
if err := serModel.checkPassword("password0"); err != nil { password error }*/
func (user *User) IsValidPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(user.Password)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

//BeforeSave ---ajsdkh
func (user *User) BeforeSave(db *gorm.DB) (err error) {
	if len(user.Roles) == 0 {
		// role := Role{}
		userRole := Role{}
		// db.Model(&role).Where("name = ?", "ROLE_USER").First(&userRole)
		db.Model(&Role{}).Where("name = ?", "ROLE_USER").First(&userRole)
		//db.Where(&models.Role{Name: "ROLE_USER"}).Attrs(models.Role{Description: "For standard Users"}).FirstOrCreate(&userRole)
		user.Roles = append(user.Roles, userRole)
	}
	return
}

//GenerateJwtToken Generate JWT token associated to this user
func (user *User) GenerateJwtToken() string {
	// jwt.New(jwt.GetSigningMethod("HS512"))
	jwtToken := jwt.New(jwt.SigningMethodHS512)

	var roles []string
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}

	jwtToken.Claims = jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"roles":    roles,
		"exp":      time.Now().Add(time.Hour * 24 * 90).Unix(),
	}
	// Sign and get the complete encoded token as a string
	token, _ := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return token
}

//IsAdmin ---
func (user *User) IsAdmin() bool {
	for _, role := range user.Roles {
		if role.Name == "ROLE_ADMIN" {
			return true
		}
	}
	return false
}

//IsNotAdmin ---
func (user *User) IsNotAdmin() bool {
	return !user.IsAdmin()
}
