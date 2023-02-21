package model

type IUser interface {
	SetName(name string)
	SetLastName(lastName string)
	SetIDNumber(numberID string)
	SetEmail(email string)
	SetUsername(username string)
	SetPassword(password string)
	SetEnabled(enabled bool)

	GetUserID() uint64
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
	ID        uint64 `db:"id" json:"id" gorm:"primaryKey,autoIncrement"`
	Name      string `db:"name" json:"name" validate:"required"`
	LastName  string `db:"last_name" json:"last_name" validate:"required"`
	IDNumber  string `db:"id_number" json:"id_number" validate:"required"`
	Username  string `db:"username" json:"username,omitempty" gorm:"unique" validate:"required"`
	Password  string `db:"password" json:"password,omitempty" validate:"required"`
	Email     string `db:"email" json:"email,omitempty" gorm:"unique" validate:"required"`
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
