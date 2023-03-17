package model

import "github.com/go-playground/validator/v10"

type Address struct {
	ID             uint64 `json:"id"`
	Name           string `json:"name" validate:"required"`
	Street         string `json:"street" validate:"required"`
	Number         string `json:"number" validate:"required"`
	Floor          string `json:"floor,omitempty"`
	Apartment      string `json:"apartment,omitempty"`
	ZipCode        string `json:"zip_code" validate:"required"`
	City           string `json:"city" validate:"required"`
	State          string `json:"state" validate:"required"`
	Country        string `json:"country" validate:"required"`
	Latitude       string `json:"latitude" validate:"required"`
	Longitude      string `json:"longitude" validate:"required"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	ObserverUserID uint64 `json:"observer_user_id" validate:"required"`
}

var addressValidate = validator.New()

func (a Address) Validate() error {
	return addressValidate.Struct(a)
}

func (a Address) ValidateID() bool {
	return a.ID != 0
}
