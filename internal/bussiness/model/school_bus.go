package model

type SchoolBus struct {
	ID               uint   `json:"id"`
	LicensePlate     string `json:"license_plate"`
	Model            string `json:"model"`
	Brand            string `json:"brand"`
	SchoolBusLicense string `json:"school_bus_license"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}
