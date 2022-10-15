package model

type UserObservee struct {
	User            User           `gorm:"one2one,embedded,foreignKey:id"`
	PrivacyKey      string         `json:"privacy_key"`
	CompanyName     string         `json:"company_name"`
	LicensePlate    string         `json:"license_plate"`
	CarRegistration string         `json:"car_registration"`
	ObserverUsers   []ObserverUser `json:"observer_users" gorm:"many2many:observee_observers"`
}

func (u *UserObservee) SetPrivacyKey(privacyKey string) {
	u.PrivacyKey = privacyKey
}

func (u *UserObservee) GetPrivacyKey() string {
	return u.PrivacyKey
}

func (u *UserObservee) SetCompanyName(companyName string) {
	u.CompanyName = companyName
}

func (u *UserObservee) GetCompanyName() string {
	return u.CompanyName
}

func (u *UserObservee) SetLicensePlate(licensePlate string) {
	u.LicensePlate = licensePlate
}

func (u *UserObservee) GetLicensePlate() string {
	return u.LicensePlate
}

func (u *UserObservee) SetCarRegistration(carRegistration string) {
	u.CarRegistration = carRegistration
}

func (u *UserObservee) GetCarRegistration() string {
	return u.CarRegistration
}

func (u *UserObservee) SetObserverUsers(observerUsers []ObserverUser) {
	u.ObserverUsers = observerUsers
}

func (u *UserObservee) GetObserverUsers() []ObserverUser {
	return u.ObserverUsers
}
