package model

type IUser interface {
	SetName(name string)
	SetLastName(lastName string)
	SetIDNumber(numberID string)
	SetEmail(email string)
	SetUsername(username string)
	SetPassword(password string)
	SetEnabled(enabled bool)

	GetUserID() uint
	GetName() string
	GetLastName() string
	GetIDNumber() string
	GetEmail() string
	GetUsername() string
	GetPassword() string
	GetEnabled() bool
	GetType() string
}

type User struct {
	ID        uint   `db:"id" json:"id" gorm:"primaryKey,autoIncrement"`
	Name      string `db:"name" json:"name"`
	LastName  string `db:"last_name" json:"last_name"`
	IDNumber  string `db:"id_number" json:"id_number"`
	Username  string `db:"username" json:"username,omitempty" gorm:"unique"`
	Password  string `db:"password" json:"password,omitempty"`
	Email     string `db:"email" json:"email,omitempty" gorm:"unique"`
	Enabled   bool   `db:"enabled" json:"enabled,omitempty"`
	Type      string `db:"type" json:"type,omitempty"`
	CreatedAt string `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt string `db:"updated_at" json:"updated_at,omitempty"`
}

func (u User) SetName(name string) {
	u.Name = name
}

func (u User) GetName() string {
	return u.Name
}

func (u User) SetLastName(lastName string) {
	u.LastName = lastName
}

func (u User) GetLastName() string {
	return u.LastName
}

func (u User) SetIDNumber(idNumber string) {
	u.IDNumber = idNumber
}

func (u User) GetIDNumber() string {
	return u.IDNumber
}

func (u User) SetEmail(email string) {
	u.Email = email
}

func (u User) GetEmail() string {
	return u.Email
}

func (u User) SetUsername(username string) {
	u.Username = username
}

func (u User) GetUsername() string {
	return u.Username
}

func (u User) SetPassword(password string) {
	u.Password = password
}

func (u User) GetPassword() string {
	return u.Password
}
