package model

type ObserverUser struct {
	User          User           `gorm:"one2one,embedded"`
	Childs        []string       `json:"childs"`
	UsersObservee []UserObservee `json:"users_observee" gorm:"many2many:observee_observer;"`
}

func (ou ObserverUser) SetChilds(childs []string) {
	ou.Childs = childs
}

func (ou ObserverUser) GetChilds() []string {
	return ou.Childs
}

func (ou ObserverUser) SetUsersObservee(usersObservee []UserObservee) {
	ou.UsersObservee = usersObservee
}

func (ou ObserverUser) GetUsersObservee() []UserObservee {
	return ou.UsersObservee
}
