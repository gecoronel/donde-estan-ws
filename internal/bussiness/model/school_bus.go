package model

import "github.com/go-playground/validator/v10"

type SchoolBus struct {
	ID             uint64 `json:"id"`
	LicensePlate   string `json:"license_plate" validate:"required"`
	Model          string `json:"model" validate:"required"`
	Brand          string `json:"brand" validate:"required"`
	License        string `json:"license" validate:"required"`
	ObservedUserID uint64 `json:"observed_user_id"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

var schoolBusValidate = validator.New()

func (sb SchoolBus) Validate() error {
	return schoolBusValidate.Struct(sb)
}

func (sb SchoolBus) ValidateID() bool {
	return sb.ID != 0
}

func (sb SchoolBus) ValidateObservedUserID() bool {
	return sb.ObservedUserID != 0
}
