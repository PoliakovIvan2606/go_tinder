package models

import (
	"github.com/asaskevich/govalidator"
)

type Preferences struct {
    User_id   int `json:"user_id"`
	Gender string `json:"gender" valid:"in(Мужской|Женский)`
	Age_from int `json:"age_from" valid:"required,range(18|65)"`
	Age_to int `json:"age_to" valid:"required,range(18|65)"`
	Radius int `json:"radius"`
}

func (p *Preferences) Validate() error {
	_, err := govalidator.ValidateStruct(p)
	if err != nil {
		return err
	}
	return nil
}