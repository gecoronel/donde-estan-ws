package model

import "github.com/go-playground/validator/v10"

type ObservedUser struct {
	User          User           `gorm:"one2one,embedded,foreignKey:id" validate:"required"`
	PrivacyKey    string         `db:"privacy_key" json:"privacy_key" gorm:"unique" validate:"required"`
	CompanyName   string         `db:"company_name" json:"company_name" validate:"required"`
	SchoolBus     SchoolBus      `db:"school_bus" json:"school_bus" validate:"required"`
	ObserverUsers []ObserverUser `json:",omitempty" gorm:"foreignKey:ObservedUsers"`
}

func NewObservedUser(observed ObservedUser) IUser {
	return observed
}

func (observed ObservedUser) GetUserID() uint64 {
	return observed.User.ID
}

func (observed ObservedUser) GetName() string {
	return observed.User.Name
}

func (observed ObservedUser) SetName(name string) {
	observed.User.Name = name
}

func (observed ObservedUser) GetLastName() string {
	return observed.User.LastName
}

func (observed ObservedUser) SetLastName(lastName string) {
	observed.User.LastName = lastName
}

func (observed ObservedUser) GetIDNumber() string {
	return observed.User.IDNumber
}

func (observed ObservedUser) SetIDNumber(idNumber string) {
	observed.User.IDNumber = idNumber
}

func (observed ObservedUser) GetEmail() string {
	return observed.User.Email
}

func (observed ObservedUser) SetEmail(email string) {
	observed.User.Email = email
}

func (observed ObservedUser) GetUsername() string {
	return observed.User.Username
}

func (observed ObservedUser) SetUsername(username string) {
	observed.User.Username = username
}

func (observed ObservedUser) GetPassword() string {
	return observed.User.Password
}

func (observed ObservedUser) SetPassword(password string) {
	observed.User.Password = password
}

func (observed ObservedUser) GetEnabled() bool {
	return observed.User.Enabled
}

func (observed ObservedUser) SetEnabled(enabled bool) {
	observed.User.Enabled = enabled
}

func (observed ObservedUser) GetType() string {
	return observed.User.Type
}

func (observed ObservedUser) GetPrivacyKey() string {
	return observed.PrivacyKey
}

func (observed ObservedUser) SetPrivacyKey(privacyKey string) {
	observed.PrivacyKey = privacyKey
}

func (observed ObservedUser) GetCompanyName() string {
	return observed.CompanyName
}

func (observed ObservedUser) SetCompanyName(companyName string) {
	observed.CompanyName = companyName
}

func (observed ObservedUser) GetLicensePlate() string {
	return observed.SchoolBus.LicensePlate
}

func (observed ObservedUser) SetLicensePlate(licensePlate string) {
	observed.SchoolBus.LicensePlate = licensePlate
}

func (observed ObservedUser) GetSchoolBusLicense() string {
	return observed.SchoolBus.License
}

func (observed ObservedUser) SetSchoolBusLicense(schoolBusLicense string) {
	observed.SchoolBus.License = schoolBusLicense
}

func (observed ObservedUser) SetObserverUsers(observerUsers []ObserverUser) {
	observed.ObserverUsers = observerUsers
}

func (observed ObservedUser) GetObserverUsers() []ObserverUser {
	return observed.ObserverUsers
}

var validateObserved = validator.New()

func (observed ObservedUser) Validate() error {
	return validateObserved.Struct(observed)
}
