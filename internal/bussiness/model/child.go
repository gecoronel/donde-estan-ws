package model

import "github.com/go-playground/validator/v10"

type Child struct {
	ID              uint64 `json:"id" gorm:"primaryKey,autoIncrement"`
	ObserverUserID  uint64 `json:"observer_user_id" validate:"required"`
	Name            string `json:"name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	SchoolName      string `json:"school_name" validate:"required"`
	SchoolStartTime string `json:"school_start_time" validate:"required"`
	SchoolEndTime   string `json:"school_end_time" validate:"required"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

var childValidate = validator.New()

func (c Child) Validate() error {
	return childValidate.Struct(c)
}