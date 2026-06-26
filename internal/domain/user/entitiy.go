package user

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"  gorm:"type:varchar(100);not null"`
	Email    string `json:"email"  gorm:"type:varchar(255);uniqueIndex;not null"`
	Password string `json:"password"  gorm:"type:varchar(100);not null"`
	Role     string `json:"role" gorm:"type:varchar(20);not null;default:driver"`
}

func (u *User) HashPassword(password string, cost int) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
