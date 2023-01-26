package model

type ObserverUser struct {
	User          User           `json:"user" gorm:"one2one,embedded,foreignKey:id"`
	Children      []Children     `json:"children" gorm:"foreignKey:ObserverUserID;"`
	ObservedUsers []ObservedUser `json:",omitempty" gorm:"many2many:ObservedUserObserverUser, foreignKey:ObserverUser;"`
}

func NewObserverUser(observer ObserverUser) IUser {
	return observer
}

func (observer ObserverUser) GetUserID() uint64 {
	return observer.User.ID
}

func (observer ObserverUser) GetName() string {
	return observer.User.Name
}

func (observer ObserverUser) SetName(name string) {
	observer.User.Name = name
}

func (observer ObserverUser) SetLastName(lastName string) {
	observer.User.LastName = lastName
}

func (observer ObserverUser) GetLastName() string {
	return observer.User.LastName
}

func (observer ObserverUser) SetIDNumber(IDNumber string) {
	observer.User.IDNumber = IDNumber
}

func (observer ObserverUser) GetIDNumber() string {
	return observer.User.IDNumber
}

func (observer ObserverUser) SetEmail(email string) {
	observer.User.Email = email
}

func (observer ObserverUser) GetEmail() string {
	return observer.User.Email
}

func (observer ObserverUser) SetUsername(username string) {
	observer.User.Username = username
}

func (observer ObserverUser) GetUsername() string {
	return observer.User.Username
}

func (observer ObserverUser) SetPassword(password string) {
	observer.User.Password = password
}

func (observer ObserverUser) GetPassword() string {
	return observer.User.Password
}

func (observer ObserverUser) GetEnabled() bool {
	return observer.User.Enabled
}

func (observer ObserverUser) SetEnabled(enabled bool) {
	observer.User.Enabled = enabled
}

func (observer ObserverUser) GetType() string {
	return observer.User.Type
}

func (observer ObserverUser) SetChilds(childs []Children) {
	observer.Children = childs
}

func (observer ObserverUser) GetChilds() []Children {
	return observer.Children
}

func (observer ObserverUser) SetUsersObservee(usersObservee []ObservedUser) {
	observer.ObservedUsers = usersObservee
}

func (observer ObserverUser) GetUsersObservee() []ObservedUser {
	return observer.ObservedUsers
}
