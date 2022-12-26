package model

type Children struct {
	ID              uint   `json:"id" gorm:"primaryKey,autoIncrement"`
	ObserverUserID  uint   `json:"observer_user_id"`
	Name            string `json:"name"`
	LastName        string `json:"last_name"`
	SchoolName      string `json:"school_name"`
	SchoolStartTime string `json:"school_start_time"`
	SchoolEndTime   string `json:"school_end_time"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}
