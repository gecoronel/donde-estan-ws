package model

import "github.com/go-playground/validator/v10"

type SchoolBus struct {
	ID           string `json:"id" validate:"required"`
	LicensePlate string `json:"license_plate" validate:"required"`
	Model        string `json:"model" validate:"required"`
	Brand        string `json:"brand" validate:"required"`
	License      string `json:"license" validate:"required"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

var schoolBusValidate = validator.New()

func (sb SchoolBus) Validate() error {
	return schoolBusValidate.Struct(sb)
}
