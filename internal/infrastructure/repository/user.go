package repository

import (
	"context"
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"sync"

	"github.com/gcoron/donde-estan-ws/internal/bussiness/gateway"
	"github.com/gcoron/donde-estan-ws/internal/bussiness/model"
	"github.com/gcoron/donde-estan-ws/internal/bussiness/model/web"
	log "github.com/sirupsen/logrus"
)

func NewUserRepository(db *gorm.DB, ctx context.Context) gateway.UserRepository {
	return &UserRepository{
		DB:      db,
		context: ctx,
	}
}

// UserRepository represents the main repository for manage user.
type UserRepository struct {
	DB      *gorm.DB
	context context.Context
}

// Get obtains a user using UserRepository by ID.
func (r UserRepository) Get(id uint) (*model.User, error) {
	var user model.User
	result := r.DB.First(&user, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, web.NewError(http.StatusNotFound, "user not found")
		}
		return nil, result.Error
	}

	return &user, nil
}

// Save persists a user using UserRepository.
func (r UserRepository) Save(user model.User) (*model.User, error) {
	result := r.DB.Model(&user).Create(&user)

	return &user, result.Error
}

// FindByUsername obtains a user using UserRepository by username.
func (r UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User

	err := r.DB.
		Raw("SELECT * FROM USERS WHERE username = @username", sql.Named("username", username)).
		Row().
		Scan(&user.ID, &user.Name, &user.LastName, &user.IDNumber, &user.Username, &user.Password, &user.Email,
			&user.Enabled, &user.Type, &user.CreatedAt, &user.UpdatedAt,
		)

	if err != nil {
		log.Error("error row scan")
		return nil, err
	}

	return &user, nil
}

// GetUsers obtains users using UserRepository .
func (r UserRepository) GetUsers(limit string, offset string) (*[]model.User, error) {
	var user model.User
	var users []model.User

	rows, err := r.DB.Raw(
		"SELECT * FROM USERS LIMIT @limit OFFSET @offset",
		sql.Named("limit", limit),
		sql.Named("offset", offset),
	).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Name, &user.LastName, &user.IDNumber, &user.Username, &user.Password, &user.Email)
		if err != nil {
			log.Error("error rows scan")
			return &users, err
		}

		users = append(users, user)
	}

	return &users, nil
}

// GetObservedUser obtains a observedUser using UserRepository by user_id.
func (r UserRepository) GetObservedUser(user *model.ObservedUser) (*model.IUser, error) {
	err := r.DB.
		Raw(
			"SELECT ou.user_id, ou.school_bus_id, ou.privacy_key, ou.company_name, sb.license_plate, sb.model, sb.brand, sb.school_bus_license, sb.created_at, sb.updated_at FROM ObservedUsers AS ou INNER JOIN SchoolBuses AS sb WHERE user_id = @user_id",
			sql.Named("user_id", user.User.ID),
		).
		Row().
		Scan(
			&user.User.ID,
			&user.SchoolBus.ID,
			&user.PrivacyKey,
			&user.CompanyName,
			&user.SchoolBus.LicensePlate,
			&user.SchoolBus.Model,
			&user.SchoolBus.Brand,
			&user.SchoolBus.SchoolBusLicense,
			&user.SchoolBus.CreatedAt,
			&user.SchoolBus.UpdatedAt,
		)

	if err != nil {
		log.Error("error row scan")
		return nil, err
	}

	U := model.NewObservedUser(*user)
	return &U, nil
}

// GetObserverUser obtains a observerUser using UserRepository by user_id.
func (r UserRepository) GetObserverUser(user *model.ObserverUser) (*model.IUser, error) {
	var (
		errChildren           error
		errObservedUser       error
		err                   error
		statementChildren     = "SELECT c.id, c.name, c.last_name, c.school_name, c.school_start_time, c.school_end_time, c.observer_user_id, c.created_at, c.updated_at FROM ObserverUsers AS oru INNER JOIN Children AS c ON  oru.user_id = c.observer_user_id;"
		statementObservedUser = "SELECT u.id, u.name, u.last_name, u.id_number, odu.company_name, odu.privacy_key, sb.id AS school_bus_id, sb.license_plate, sb.model, sb.brand, sb.school_bus_license, sb.created_at, sb.updated_at FROM ObserverUsers AS oru INNER JOIN ObservedUsers AS odu INNER JOIN ObservedUsersObserverUsers AS oduoru INNER JOIN Users AS u INNER JOIN SchoolBuses AS sb ON odu.user_id = oduoru.observed_user_id AND oru.user_id = oduoru.observer_user_id AND u.id = odu.user_id AND odu.school_bus_id = sb.id;"
		children              []model.Children
		child                 model.Children
		observedUsers         []odUser
		observedUser          odUser
		u                     model.IUser
		wg                    = &sync.WaitGroup{}
		usersCompleted        = make(chan struct{}, 3)
		chanErr               = make(chan error, 1)
	)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer handleGoRoutinePanic(wg)
		children, errChildren = scanRows(r, statementChildren, children, child)
		if errChildren != nil {
			chanErr <- errChildren
			return
		}
	}(wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer handleGoRoutinePanic(wg)
		observedUsers, errObservedUser = scanRows(r, statementObservedUser, observedUsers, observedUser)
		if errObservedUser != nil {
			chanErr <- errObservedUser
			return
		}
	}(wg)

	go func() {
		wg.Wait()

		user.Children = children
		user.ObservedUsers = mapToObservedUser(observedUsers)

		u = model.NewObserverUser(*user)

		usersCompleted <- struct{}{}
	}()

	select {
	case <-usersCompleted:
		return &u, nil
	case err = <-chanErr:
		return nil, err
	}
}

func scanRows[T allowScan](r UserRepository, statement string, list []T, object T) ([]T, error) {
	rows, err := r.DB.
		Raw(statement).
		Rows()

	if err != nil {
		log.Error("error in scan rows ")
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = r.DB.ScanRows(rows, &object)
		if err != nil {
			return nil, err
		}
		list = append(list, object)
	}

	return list, err
}

func handleGoRoutinePanic(wg *sync.WaitGroup) {
	if r := recover(); r != nil {
		log.Errorf("%v", r)
	}

	wg.Done()
}

func mapToObservedUser(users []odUser) []model.ObservedUser {
	var observedUsers []model.ObservedUser

	for _, user := range users {
		observedUser := model.ObservedUser{
			User:        model.User{ID: user.ID, Name: user.Name, LastName: user.LastName, IDNumber: user.IDNumber},
			PrivacyKey:  user.PrivacyKey,
			CompanyName: user.CompanyName,
			SchoolBus: model.SchoolBus{
				ID:               user.SchoolBusID,
				LicensePlate:     user.LicensePlate,
				Model:            user.Model,
				Brand:            user.Brand,
				SchoolBusLicense: user.SchoolBusLicense,
				CreatedAt:        user.CreatedAt,
				UpdatedAt:        user.UpdatedAt,
			},
		}
		observedUsers = append(observedUsers, observedUser)
	}

	return observedUsers
}

type (
	allowScan interface {
		model.User | model.ObserverUser | model.ObservedUser | model.Children | model.SchoolBus | odUser
	}
	odUser struct {
		ID               uint   `json:"id"`
		Name             string `json:"name"`
		LastName         string `json:"last_name"`
		IDNumber         string `json:"id_number"`
		PrivacyKey       string `json:"privacy_key"`
		CompanyName      string `json:"company_name"`
		SchoolBusID      uint   `json:"school_bus_id"`
		LicensePlate     string `json:"license_plate"`
		Model            string `json:"model"`
		Brand            string `json:"brand"`
		SchoolBusLicense string `json:"school_bus_license"`
		CreatedAt        string `json:"created_at"`
		UpdatedAt        string `json:"updated_at"`
	}
)
