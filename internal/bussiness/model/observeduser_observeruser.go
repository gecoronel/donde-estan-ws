package model

type ObservedUserObserverUser struct {
	ObservedUserID uint   `db:"observed_user_id" gorm:"foreignKey:user"`
	ObserverUserID uint   `db:"observer_user_id" gorm:"foreignKey:user"`
	CreatedAt      string `db:"created_at" json:"created_at"`
	UpdatedAt      string `db:"updated_at" json:"updated_at"`
}
