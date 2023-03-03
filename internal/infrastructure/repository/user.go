package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"sync"

	"github.com/gecoronel/donde-estan-ws/internal/bussiness/gateway"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
	log "github.com/sirupsen/logrus"
)

const (
	querySaveUser = `
		INSERT INTO Users (name, last_name, id_number, username, password, email, type) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	querySelectUserIDByID       = `SELECT u.id FROM Users AS u WHERE u.id = @id`
	querySelectUserIDByUsername = `SELECT id FROM Users WHERE username = ?`
	querySelectUserIDByEmail    = `SELECT u.id FROM Users AS u WHERE email = @email`

	queryGetObservedUser = `
		SELECT u.id, u.name, u.last_name, u.id_number, u.username, u.password, u.email, u.type, u.enabled, 
		       ou.privacy_key, ou.company_name, sb.id, sb.license_plate, sb.model, sb.brand, sb.license, 
		       sb.created_at, sb.updated_at 
		FROM ObservedUsers AS ou 
    	INNER JOIN Users AS u ON u.id = ou.user_id 
		INNER JOIN SchoolBuses AS sb ON sb.id = ou.school_bus_id 
		WHERE u.id = ?;
	`
	queryGetObservedUserByPrivacyKey = `
		SELECT u.id, u.name, u.last_name, u.id_number, u.username, u.password, u.email, u.type, u.enabled, 
		       ou.privacy_key, ou.company_name, sb.id AS school_bus_id, sb.license_plate, sb.model, sb.brand, sb.license, 
		       sb.created_at, sb.updated_at 
		FROM ObservedUsers AS ou 
    	INNER JOIN Users AS u ON u.id = ou.user_id 
		INNER JOIN SchoolBuses AS sb ON sb.id = ou.school_bus_id 
		WHERE ou.privacy_key = ?;
	`
	querySaveObservedUser = `
		INSERT INTO ObservedUsers (user_id, privacy_key, company_name, school_bus_id) VALUES (?, ?, ?, ?);
	`
	queryGetObservedUserByID = `SELECT user_id, privacy_key, company_name FROM ObservedUsers WHERE user_id = ?`

	queryGetUser = `
		SELECT u.id, u.name, u.last_name, u.id_number, u.username, u.password, u.email, u.type, u.enabled 
		FROM Users AS u 
		WHERE u.id = %d;
	`
	queryGetChildren = `
		SELECT c.id, c.name, c.last_name, c.school_name, c.school_start_time, c.school_end_time, c.observer_user_id, 
		       c.created_at, c.updated_at 
		FROM ObserverUsers AS oru 
		INNER JOIN Children AS c ON oru.user_id = c.observer_user_id 
		WHERE user_id = %d;
	`
	queryGetObservedUserOfObserverUser = `
		SELECT u.id, u.name, u.last_name, u.id_number, odu.company_name, odu.privacy_key, sb.id AS school_bus_id, 
		       sb.license_plate, sb.model, sb.brand, sb.license, sb.created_at, sb.updated_at 
		FROM ObserverUsers AS oru 
		INNER JOIN ObservedUsers AS odu 
		INNER JOIN ObservedUsersObserverUsers AS oduoru 
		INNER JOIN Users AS u 
		INNER JOIN SchoolBuses AS sb ON odu.user_id = oduoru.observed_user_id 
			AND oru.user_id = oduoru.observer_user_id 
			AND u.id = odu.user_id 
			AND odu.school_bus_id = sb.id 
		WHERE oru.user_id = %d;
	`
	querySaveObservedUserInObserverUser = `
		INSERT INTO ObservedUsersObserverUsers (observed_user_id, observer_user_id)
		VALUES (1, 3);
	`
	queryDeleteObservedUserInObserverUser = `
		DELETE FROM ObservedUsersObserverUsers
		WHERE observed_user_id = ? AND observer_user_id = ?;
	`
	querySelectObservedUserObserverUser = `
		SELECT * FROM ObservedUsersObserverUsers 
		WHERE observed_user_id = ? 
		  AND observer_user_id = ?;`
)

func NewUserRepository(db *gorm.DB, ctx context.Context) gateway.UserRepository {
	return &UserRepository{
		DB:      db,
		context: ctx,
	}
}

// UserRepository represents the main repository for manage user
type UserRepository struct {
	DB      *gorm.DB
	context context.Context
}

// Get obtains a user using UserRepository by ID
func (r UserRepository) Get(id uint64) (*model.User, error) {
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

// FindByUsername obtains a user using UserRepository by username
func (r UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User

	err := r.DB.
		Raw("SELECT * FROM Users WHERE username = @username", sql.Named("username", username)).
		Row().
		Scan(
			&user.ID, &user.Name, &user.LastName, &user.IDNumber, &user.Username, &user.Password, &user.Email,
			&user.Enabled, &user.Type, &user.CreatedAt, &user.UpdatedAt,
		)

	if err != nil {
		if err.Error() == web.ErrNoRows.Error() {
			log.Error("error row scan not found")
			return nil, nil
		}
		log.Error("error row scan")
		return nil, err
	}

	return &user, nil
}

// FindByEmail obtains a user using UserRepository by email
func (r UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User

	err := r.DB.
		Raw("SELECT * FROM Users WHERE email = @email", sql.Named("email", email)).
		Row().
		Scan(
			&user.ID, &user.Name, &user.LastName, &user.IDNumber, &user.Username, &user.Password, &user.Email,
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
		err := rows.Scan(
			&user.ID, &user.Name, &user.LastName, &user.IDNumber, &user.Username, &user.Password, &user.Email,
		)
		if err != nil {
			log.Error("error rows scan")
			return &users, err
		}

		users = append(users, user)
	}

	return &users, nil
}

// GetObservedUser obtains a observedUser using UserRepository by user_id
func (r UserRepository) GetObservedUser(id uint64) (*model.ObservedUser, error) {
	var observed model.ObservedUser

	err := r.DB.
		Raw(queryGetObservedUser, id).
		Row().
		Scan(
			&observed.User.ID, &observed.User.Name, &observed.User.LastName, &observed.User.IDNumber,
			&observed.User.Username, &observed.User.Password, &observed.User.Email, &observed.User.Type,
			&observed.User.Enabled, &observed.SchoolBus.ID, &observed.PrivacyKey, &observed.CompanyName,
			&observed.SchoolBus.LicensePlate, &observed.SchoolBus.Model, &observed.SchoolBus.Brand,
			&observed.SchoolBus.License, &observed.SchoolBus.CreatedAt, &observed.SchoolBus.UpdatedAt,
		)

	if err != nil {
		log.Error("error row scan")
		return nil, err
	}

	return &observed, nil
}

// SaveObservedUser create a observedUser using UserRepository.
func (r UserRepository) SaveObservedUser(user model.ObservedUser) (*model.ObservedUser, error) {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Error("recovering for error in save observed user")
			tx.Rollback()
		}
	}()

	err := saveUserExec(tx, &user.User)
	if err != nil {
		log.Error("error creating user")
		tx.Rollback()
		return nil, err
	}

	err = selectIDByUsernameExec(tx, &user.User)
	if err != nil {
		log.Error("error creating user")
		tx.Rollback()
		return nil, err
	}

	err = saveSchoolBusExec(tx, &user)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = selectSchoolBusByIDExec(tx, &user)
	if err != nil {
		log.Error("error creating school bus")
		tx.Rollback()
		return nil, err
	}

	err = saveObservedUserExec(tx, &user)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = selectObservedUserByID(tx, &user)
	if err != nil {
		log.Error("error creating observed user")
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &user, nil
}

// GetObserverUser obtains a observerUser using UserRepository by user_id
func (r UserRepository) GetObserverUser(id uint64) (*model.ObserverUser, error) {
	var (
		errUser               error
		errChildren           error
		errObservedUser       error
		statementUser         = fmt.Sprintf(queryGetUser, id)
		statementChildren     = fmt.Sprintf(queryGetChildren, id)
		statementObservedUser = fmt.Sprintf(queryGetObservedUserOfObserverUser, id)
		children              []model.Child
		child                 model.Child
		observedUsers         []odUser
		observedUser          odUser
		observer              model.ObserverUser
		wg                    = &sync.WaitGroup{}
	)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer handleGoRoutinePanic(wg)
		errUser = r.DB.
			Raw(statementUser).
			Row().
			Scan(
				&observer.User.ID, &observer.User.Name, &observer.User.LastName, &observer.User.IDNumber, &observer.User.Username, &observer.User.Password,
				&observer.User.Email, &observer.User.Type, &observer.User.Enabled,
			)
	}(wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer handleGoRoutinePanic(wg)
		children, errChildren = scanRows(r.DB, statementChildren, children, child)
	}(wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer handleGoRoutinePanic(wg)
		observedUsers, errObservedUser = scanRows(r.DB, statementObservedUser, observedUsers, observedUser)
	}(wg)

	wg.Wait()

	if errUser != nil {
		log.Errorf("error user rows: %s", errUser.Error())
		return nil, errUser
	}

	if errChildren != nil {
		log.Errorf("error children rows: %s", errChildren.Error())
		return nil, errChildren
	}

	if errObservedUser != nil {
		log.Errorf("error observed user rows: %s", errObservedUser.Error())
		return nil, errObservedUser
	}

	observer.Children = children
	observer.ObservedUsers = mapToObservedUser(observedUsers)
	return &observer, nil
}

// SaveObserverUser create a observerUser using UserRepository
func (r UserRepository) SaveObserverUser(u model.ObserverUser) (*model.ObserverUser, error) {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Error("recovering for error in save observed user")
			tx.Rollback()
		}
	}()

	err := saveUserExec(tx, &u.User)
	if err != nil {
		log.Error("error creating user")
		tx.Rollback()
		return nil, err
	}

	err = selectIDByUsernameExec(tx, &u.User)
	if err != nil {
		log.Error("error creating user")
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &u, nil
}

// FindObservedUserByPrivacyKey obtains a observedUser using UserRepository by privacy_key
func (r UserRepository) FindObservedUserByPrivacyKey(privacyKey string) (*model.ObservedUser, error) {
	var observed model.ObservedUser

	err := r.DB.
		Raw(queryGetObservedUserByPrivacyKey, privacyKey).
		Row().
		Scan(
			&observed.User.ID, &observed.User.Name, &observed.User.LastName, &observed.User.IDNumber,
			&observed.User.Username, &observed.User.Password, &observed.User.Email, &observed.User.Type,
			&observed.User.Enabled, &observed.PrivacyKey, &observed.CompanyName, &observed.SchoolBus.ID,
			&observed.SchoolBus.LicensePlate, &observed.SchoolBus.Model, &observed.SchoolBus.Brand,
			&observed.SchoolBus.License, &observed.SchoolBus.CreatedAt, &observed.SchoolBus.UpdatedAt,
		)

	if err != nil {
		if err.Error() == web.ErrNoRows.Error() {
			log.Error("error row scan not found")
			return nil, nil
		}
		log.Error("error row scan")
		return nil, err
	}

	return &observed, nil
}

// SaveObservedUserInObserverUser add a observedUser into of the drivers list of an observerUser using UserRepository.
func (r UserRepository) SaveObservedUserInObserverUser(observedUserID uint64, observerUserID uint64) error {
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Error("recovering for error saving observed user")
			tx.Rollback()
		}
	}()

	err := tx.Exec(querySaveObservedUserInObserverUser, observedUserID, observerUserID).Error
	if err != nil {
		log.Error("error creating user")
		tx.Rollback()
		return err
	}

	err = tx.
		Raw(querySelectObservedUserObserverUser, observedUserID, observerUserID).
		Row().
		Scan(&observedUserID, &observerUserID)
	if err != nil {
		log.Error("error selecting observed user and observer user")
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// DeleteObservedUserInObserverUser delete a observedUser into of the drivers list of an observerUser using UserRepository.
func (r UserRepository) DeleteObservedUserInObserverUser(observedUserID uint64, observerUserID uint64) error {
	err := r.DB.Exec(queryDeleteObservedUserInObserverUser, observedUserID, observerUserID).Error
	if err != nil {
		log.Error("error deleting user")
		return err
	}

	return nil
}

func selectObservedUserByID(tx *gorm.DB, user *model.ObservedUser) error {
	return tx.
		Raw(queryGetObservedUserByID, user.User.ID).
		Row().
		Scan(&user.User.ID, &user.PrivacyKey, &user.CompanyName)
}

func saveObservedUserExec(tx *gorm.DB, user *model.ObservedUser) error {
	return tx.Exec(
		querySaveObservedUser,
		user.User.ID,
		user.PrivacyKey,
		user.CompanyName,
		user.SchoolBus.ID,
	).Error
}

func selectSchoolBusByIDExec(tx *gorm.DB, user *model.ObservedUser) error {
	return tx.
		Raw(querySelectSchoolBusByID, user.SchoolBus.ID).
		Row().
		Scan(&user.SchoolBus.ID, &user.SchoolBus.LicensePlate, &user.SchoolBus.Model, &user.SchoolBus.Brand,
			&user.SchoolBus.License, &user.SchoolBus.CreatedAt, &user.SchoolBus.UpdatedAt)
}

func saveSchoolBusExec(tx *gorm.DB, user *model.ObservedUser) error {
	return tx.Exec(
		querySaveSchoolBus,
		user.SchoolBus.ID,
		user.SchoolBus.LicensePlate,
		user.SchoolBus.Model,
		user.SchoolBus.Brand,
		user.SchoolBus.License,
	).Error
}

func selectIDByUsernameExec(tx *gorm.DB, user *model.User) error {
	return tx.
		Raw(querySelectUserIDByUsername, user.Username).
		Row().
		Scan(&user.ID)
}

func saveUserExec(tx *gorm.DB, user *model.User) error {
	return tx.Exec(
		querySaveUser,
		user.Name,
		user.LastName,
		user.IDNumber,
		user.Username,
		user.Password,
		user.Email,
		user.Type,
	).Error
}

func scanRows[T allowScan](db *gorm.DB, statement string, list []T, object T) ([]T, error) {
	rows, err := db.
		Raw(statement).
		Rows()

	if err != nil {
		log.Errorf("error scaning rows: %s", err.Error())
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = db.ScanRows(rows, &object)
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
				ID:           user.SchoolBusID,
				LicensePlate: user.LicensePlate,
				Model:        user.Model,
				Brand:        user.Brand,
				License:      user.SchoolBusLicense,
				CreatedAt:    user.CreatedAt,
				UpdatedAt:    user.UpdatedAt,
			},
		}
		observedUsers = append(observedUsers, observedUser)
	}

	return observedUsers
}

type (
	allowScan interface {
		model.User | model.ObserverUser | model.ObservedUser | model.Child | model.SchoolBus | odUser
	}
	odUser struct {
		ID               uint64 `json:"id"`
		Name             string `json:"name"`
		LastName         string `json:"last_name"`
		IDNumber         string `json:"id_number"`
		PrivacyKey       string `json:"privacy_key"`
		CompanyName      string `json:"company_name"`
		SchoolBusID      string `json:"school_bus_id"`
		LicensePlate     string `json:"license_plate"`
		Model            string `json:"model"`
		Brand            string `json:"brand"`
		SchoolBusLicense string `json:"school_bus_license"`
		CreatedAt        string `json:"created_at"`
		UpdatedAt        string `json:"updated_at"`
	}
)
