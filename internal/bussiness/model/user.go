package model

type IUser interface {
	SetUserId(id int)
	SetName(name string)
	SetLastName(lastName string)
	SetNumberId(numberID string)
	SetEmail(email string)
	SetUsername(username string)
	SetPassword(password string)
	GetUserId() int
	GetName() string
	GetLastName() string
	GetNumberId() string
	GetEmail() string
	GetUsername() string
	GetPassword() string
}

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey,autoIncrement"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	NumberID string `json:"number_id"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Email    string `json:"email" gorm:"unique"`
}

func (u *User) SetName(name string) {
	u.Name = name
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) SetLastName(lastName string) {
	u.LastName = lastName
}

func (u *User) GetLastName() string {
	return u.LastName
}

func (u *User) SetNumberID(numberID string) {
	u.NumberID = numberID
}

func (u *User) GetNumberID() string {
	return u.NumberID
}

func (u *User) SetEmail(email string) {
	u.Email = email
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) SetUsername(username string) {
	u.Username = username
}

func (u *User) GetUsername() string {
	return u.Username
}

func (u *User) SetPassword(password string) {
	u.Password = password
}

func (u *User) GetPassword() string {
	return u.Password
}
