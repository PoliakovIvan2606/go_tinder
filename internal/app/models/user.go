package models

import (
	// "errors"
	"fmt"
	// "regexp"

	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
    ID  int `json:"id"`
    UserCreate
    CreatedAt   string `json:"created_at"`
    Photos []string `json:"photos"`
}

type UserCreate struct {
    Name        string `json:"name" valid:"required"`
    Email       string `json:"email" valid:"required,email"`
    Password    string `json:"password" valid:"required"` // сырой пароль
    PasswordHash string // захешированный пароль
    Age         int    `json:"age" valid:"required,range(18|65)"`
    Description string `json:"description"`
    City        string `json:"city" valid:"required"`
    // Coordinates string `json:"coordinates" valid:"required"`
    Latitude    float64 `json:"latitude" valid:"required"`
	Longitude   float64 `json:"longitude" valid:"required"`
}

type UserCheck struct {
    Email       string `json:"email" valid:"required,email"`
    Password    string `json:"password" valid:"required"`
}

func (u *UserCreate) Validate() error {
    _, err := govalidator.ValidateStruct(u)
    if err != nil {
        // if !isValidPoint(u.Coordinates) {
        //     return errors.New(err.Error() + " 'coordinate': not valid (expected format: 'latitude,longitude')")
        // }
        return err
    }


	// вручную проверяем диапазон координат
	if u.Latitude < -90 || u.Latitude > 90 {
		return fmt.Errorf("широта должна быть в диапазоне от -90 до 90")
	}
	if u.Longitude < -180 || u.Longitude > 180 {
		return fmt.Errorf("долгота должна быть в диапазоне от -180 до 180")
	}
    // if !isValidPoint(u.Coordinates) {
    //     return errors.New("coordinate not valid (expected format: 'latitude,longitude')")
    // }

    return nil
}

func (u *UserCreate) ToWKT() string {
	return fmt.Sprintf("SRID=4326;POINT(%f %f)", u.Longitude, u.Latitude)
}

// var pointRegex = regexp.MustCompile(`^\(\s*[-+]?\d+(\.\d+)?\s*,\s*[-+]?\d+(\.\d+)?\s*\)$`)

// // Валидация координат
// func isValidPoint(s string) bool {
//     return pointRegex.MatchString(s)
// }

// Хеширование пароля
func (u *UserCreate) HashPassword() error {
    bytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    u.PasswordHash = string(bytes)
    return err
}

// Проверка пароля
func CheckPasswordHash(passwordHash, password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
    return err == nil
}