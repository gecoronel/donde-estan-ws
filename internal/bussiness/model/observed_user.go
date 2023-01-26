package model

type ObservedUser struct {
	User          User           `gorm:"one2one,embedded,foreignKey:id"`
	PrivacyKey    string         `db:"privacy_key" json:"privacy_key" gorm:"unique"`
	CompanyName   string         `db:"company_name" json:"company_name"`
	SchoolBus     SchoolBus      `db:"school_bus" json:"school_bus"`
	ObserverUsers []ObserverUser `json:",omitempty" gorm:"foreignKey:ObservedUsers"`
}

func NewObservedUser(observed ObservedUser) IUser {
	return observed
}

func (u ObservedUser) GetUserID() uint64 {
	return u.User.ID
}

func (u ObservedUser) GetName() string {
	return u.User.Name
}

func (u ObservedUser) SetName(name string) {
	u.User.Name = name
}

func (u ObservedUser) GetLastName() string {
	return u.User.LastName
}

func (u ObservedUser) SetLastName(lastName string) {
	u.User.LastName = lastName
}

func (u ObservedUser) GetIDNumber() string {
	return u.User.IDNumber
}

func (u ObservedUser) SetIDNumber(idNumber string) {
	u.User.IDNumber = idNumber
}

func (u ObservedUser) GetEmail() string {
	return u.User.Email
}

func (u ObservedUser) SetEmail(email string) {
	u.User.Email = email
}

func (u ObservedUser) GetUsername() string {
	return u.User.Username
}

func (u ObservedUser) SetUsername(username string) {
	u.User.Username = username
}

func (u ObservedUser) GetPassword() string {
	return u.User.Password
}

func (u ObservedUser) SetPassword(password string) {
	u.User.Password = password
}

func (u ObservedUser) GetEnabled() bool {
	return u.User.Enabled
}

func (u ObservedUser) SetEnabled(enabled bool) {
	u.User.Enabled = enabled
}

func (u ObservedUser) GetType() string {
	return u.User.Type
}

func (u ObservedUser) GetPrivacyKey() string {
	return u.PrivacyKey
}

func (u ObservedUser) SetPrivacyKey(privacyKey string) {
	u.PrivacyKey = privacyKey
}

func (u ObservedUser) GetCompanyName() string {
	return u.CompanyName
}

func (u ObservedUser) SetCompanyName(companyName string) {
	u.CompanyName = companyName
}

func (u ObservedUser) GetLicensePlate() string {
	return u.SchoolBus.LicensePlate
}

func (u ObservedUser) SetLicensePlate(licensePlate string) {
	u.SchoolBus.LicensePlate = licensePlate
}

func (u ObservedUser) GetSchoolBusLicense() string {
	return u.SchoolBus.SchoolBusLicense
}

func (u ObservedUser) SetSchoolBusLicense(schoolBusLicense string) {
	u.SchoolBus.SchoolBusLicense = schoolBusLicense
}

func (u ObservedUser) SetObserverUsers(observerUsers []ObserverUser) {
	u.ObserverUsers = observerUsers
}

func (u ObservedUser) GetObserverUsers() []ObserverUser {
	return u.ObserverUsers
}
