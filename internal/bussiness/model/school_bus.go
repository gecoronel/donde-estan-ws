package model

type SchoolBus struct {
	ID           string `json:"id"`
	LicensePlate string `json:"license_plate"`
	Model        string `json:"model"`
	Brand        string `json:"brand"`
	License      string `json:"license"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
