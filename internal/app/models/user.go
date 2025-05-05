package models

import (
	"errors"
	"regexp"

	"github.com/asaskevich/govalidator"
)

type User struct {
    ID  int `json:"id"`
    UserCreate
    CreatedAt   string `json:"created_at"`
    Photos []string `json:"photos"`
}

type UserCreate struct {
    Name        string `json:"name" valid:"required"`
    Age         int    `json:"age" valid:"required,range(18|65)"`
    Description string `json:"description"`
    City        string `json:"city" valid:"required"`
    Coordinates string `json:"coordinates" valid:"required"`
}

func (u *UserCreate) Validate() error {
    _, err := govalidator.ValidateStruct(u)
    if err != nil {
        if !isValidPoint(u.Coordinates) {
            return errors.New(err.Error() + " 'coordinate': not valid (expected format: 'latitude,longitude')")
        }
        return err
    }
    if !isValidPoint(u.Coordinates) {
        return errors.New("coordinate not valid (expected format: 'latitude,longitude')")
    }

    return nil
}


var pointRegex = regexp.MustCompile(`^\(\s*[-+]?\d+(\.\d+)?\s*,\s*[-+]?\d+(\.\d+)?\s*\)$`)

func isValidPoint(s string) bool {
    return pointRegex.MatchString(s)
}