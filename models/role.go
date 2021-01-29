package models

import "github.com/jinzhu/gorm"

//Role ---
type Role struct {
	gorm.Model
	Name        string
	Description string
	Users       []User     `gorm:"many2many:users_roles;"`
	UserRoles   []UserRole `gorm:"foreignkey:RoleId"`
}

//UserRole ---
type UserRole struct {
	User   User `gorm:"association_foreignkey:UserId"`
	UserID uint
	Role   User `gorm:"association_foreignkey:RoleId"`
	RoleID uint
}

//TableName ---
func (UserRole) TableName() string {
	return "users_roles"
}

//Any ---
func Any(roles []Role, f func(Role) bool) bool {
	for _, role := range roles {
		if f(role) {
			return true
		}
	}
	return false
}
