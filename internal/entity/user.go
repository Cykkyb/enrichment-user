package entity

import (
	"github.com/go-playground/validator/v10"
)

type User struct {
	Id          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name" validate:"required"`
	Surname     string `json:"surname" db:"surname" validate:"required"`
	Patronymic  string `json:"patronymic" db:"patronymic"`
	Age         int    `json:"age" db:"age" validate:"required"`
	Gender      string `json:"gender" db:"gender" validate:"required"`
	Nationality string `json:"nationality" db:"nationality" validate:"required"`
}

type UserFilter struct {
	Id          int    `form:"id"`
	Name        string `form:"name"`
	Surname     string `form:"surname"`
	Age         int    `form:"age"`
	MinAge      int    `form:"min_age"`
	MaxAge      int    `form:"max_age"`
	Offset      int    `form:"offset"`
	Limit       int    `form:"limit"`
	Nationality string `form:"nationality"`
}

func (u *User) Validate() error {
	validate := validator.New()

	if err := validate.Struct(u); err != nil {
		return err
	}

	return nil
}
