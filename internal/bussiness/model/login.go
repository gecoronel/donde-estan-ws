package model

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

var validate = validator.New()

func (l Login) Validate() error {
	return validate.Struct(l)
}

func (l Login) String() string {
	return fmt.Sprintf("username: %s, password: %s", l.Username, l.Password)
}
