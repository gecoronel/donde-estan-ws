package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var u = model.User{
	ID:        1,
	Name:      "user",
	LastName:  "user",
	IDNumber:  "11100011",
	Username:  "user",
	Password:  "user",
	Email:     "user@user.com",
	Enabled:   true,
	Type:      "observed",
	CreatedAt: "2022-12-10 17:49:30",
	UpdatedAt: "2022-12-10 17:49:30",
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestGetUser(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
		t.Fail()
	}

	ur := NewUserRepository(gdb, context.Background())

	query := "SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1"

	t.Run("GetUser successful", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "email", "username"}).
			AddRow(u.ID, u.Name, u.Email, u.Username)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(u.ID).WillReturnRows(rows)

		user, err := ur.Get(u.ID)
		assert.NotNil(t, user)
		assert.NoError(t, err)
		assert.Equal(t, u.ID, user.ID)
	})

	t.Run("GetUser error", func(t *testing.T) {
		query = "SELECT * FROM `invalid_table`"
		rows := sqlmock.NewRows([]string{"id", "name", "email", "username"}).AddRow(u.ID, u.Name, u.Email, u.Username)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(u.ID).WillReturnRows(rows)

		user, err := ur.Get(u.ID)
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("GetUser not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(u.ID).WillReturnError(gorm.ErrRecordNotFound)

		user, err := ur.Get(u.ID)
		assert.Error(t, err)
		assert.Nil(t, user)
	})

}

func TestFindByUsername(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{}) // open gorm db
	if err != nil {
		log.Error("error opening database connection")
	}

	ur := NewUserRepository(gdb, context.Background())

	query := "SELECT * FROM Users WHERE username = ?"

	t.Run("FindByUsername successful", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "lastname", "id_number", "username", "password", "email", "enabled", "type", "created_at", "updated_at"}).
			AddRow(u.ID, u.Name, u.LastName, u.IDNumber, u.Username, u.Password, u.Email, u.Enabled, u.Type, u.CreatedAt, u.UpdatedAt)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(u.Username).WillReturnRows(rows)

		user, err := ur.FindByUsername(u.Username)
		assert.NotNil(t, user)
		assert.NoError(t, err)
		assert.Equal(t, u.Username, user.Username)
	})

	t.Run("FindByUsername scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "lastname", "id_number", "username", "password"}).
			AddRow(u.ID, u.Name, u.LastName, u.IDNumber, u.Username, u.Password)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(u.Username).WillReturnRows(rows)

		user, err := ur.FindByUsername(u.Username)
		assert.Nil(t, user)
		assert.Error(t, err)
	})

	t.Run("FindByUsername rows error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(u.Username).WillReturnError(web.ErrInternalServerError)

		user, err := ur.FindByUsername(u.Username)
		assert.Nil(t, user)
		assert.Error(t, err)
	})

	t.Run("FindByUsername nil rows error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(u.Username).WillReturnError(web.ErrNoRows)

		user, err := ur.FindByUsername(u.Username)
		assert.Nil(t, user)
		assert.Nil(t, err)
	})

}

func TestFindByEmail(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
	}

	ur := NewUserRepository(gdb, context.Background())

	query := "SELECT * FROM Users WHERE email = ?"

	t.Run("FindByEmail successful", func(t *testing.T) {
		rows := sqlmock.NewRows(
			[]string{"id", "name", "lastname", "id_number", "username", "password", "email", "enabled", "type",
				"created_at", "updated_at"},
		).
			AddRow(
				u.ID, u.Name, u.LastName, u.IDNumber, u.Username, u.Password, u.Email, u.Enabled, u.Type, u.CreatedAt,
				u.UpdatedAt,
			)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(u.Email).WillReturnRows(rows)

		user, err := ur.FindByEmail(u.Email)
		assert.NotNil(t, user)
		assert.NoError(t, err)
		assert.Equal(t, u.Email, user.Email)
	})

	t.Run("FindByEmail scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "lastname", "id_number", "username", "password"}).
			AddRow(u.ID, u.Name, u.LastName, u.IDNumber, u.Username, u.Password)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(u.Email).WillReturnRows(rows)

		user, err := ur.FindByEmail(u.Email)
		assert.Nil(t, user)
		assert.Error(t, err)
	})

	t.Run("FindByEmail rows error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(u.Email).WillReturnError(web.ErrInternalServerError)

		user, err := ur.FindByEmail(u.Email)
		assert.Nil(t, user)
		assert.Error(t, err)
	})

	t.Run("FindByEmail nil rows error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "lastname", "numberID", "username", "password", "email"}).
			AddRow(u.ID, u.Name, u.LastName, u.IDNumber, "invalid username", u.Password, u.Email)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("not found").WillReturnRows(rows)

		user, err := ur.FindByEmail(u.Email)
		assert.Nil(t, user)
		assert.Error(t, err)
	})

}

func TestGetUsers(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{}) // open gorm db
	if err != nil {
		log.Error("error opening database connection")
		t.Fail()
	}

	ur := NewUserRepository(gdb, context.Background())

	query := "SELECT * FROM USERS LIMIT ? OFFSET ?"

	t.Run("GetUsers successful", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "lastname", "numberID", "username", "password", "email"}).
			AddRow(u.ID, u.Name, u.LastName, u.IDNumber, u.Username, u.Password, u.Email).
			AddRow(u.ID, u.Name, u.LastName, u.IDNumber, u.Username, u.Password, u.Email).
			AddRow(u.ID, u.Name, u.LastName, u.IDNumber, u.Username, u.Password, u.Email).
			AddRow(u.ID, u.Name, u.LastName, u.IDNumber, u.Username, u.Password, u.Email)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("1", "4").WillReturnRows(rows)

		users, err := ur.GetUsers("1", "4")
		assert.NotNil(t, users)
		assert.NoError(t, err)
		assert.Equal(t, 4, len(*users))
	})

	t.Run("GetUsers users scan error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("1", "4").WillReturnError(web.ErrInternalServerError)

		users, err := ur.GetUsers("1", "4")
		assert.Nil(t, users)
		assert.Error(t, err)
	})

	t.Run("GetUsers users scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "lastname", "numberID", "username", "password"}).
			AddRow(u.ID, u.Name, u.LastName, u.IDNumber, u.Username, u.Password).
			AddRow(u.ID, u.Name, u.LastName, u.IDNumber, u.Username, u.Password).
			AddRow(u.ID, u.Name, u.LastName, u.IDNumber, u.Username, u.Password).
			AddRow(u.ID, u.Name, u.LastName, u.IDNumber, u.Username, u.Password)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("1", "4").WillReturnRows(rows)

		users, err := ur.GetUsers("1", "4")
		assert.Nil(t, *users)
		assert.Error(t, err)
	})

}

var observed = model.ObservedUser{
	User: model.User{
		ID:        1,
		Name:      "Juan",
		LastName:  "Perez",
		IDNumber:  "12345678",
		Username:  "jperez",
		Password:  "jperez1234",
		Email:     "jperez@mail.com",
		Enabled:   true,
		Type:      "observed",
		CreatedAt: "2022-12-10 17:49:30",
		UpdatedAt: "2022-12-10 17:49:30",
	},
	PrivacyKey:  "juan.perez.12345678",
	CompanyName: "school bus company",
	SchoolBus: model.SchoolBus{
		ID:           1,
		LicensePlate: "11AAA22",
		Model:        "Master",
		Brand:        "Renault",
		License:      "111",
		CreatedAt:    "2023-02-18 17:09:33",
		UpdatedAt:    "2023-02-18 17:09:33",
	},
}

func TestSaveObservedUser(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
	}

	ur := NewUserRepository(gdb, context.Background())

	t.Run("SaveObservedUser successful", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectExec(regexp.QuoteMeta(querySaveUser)).
			WithArgs(observed.User.Name, observed.User.LastName, observed.User.IDNumber, observed.User.Username,
				observed.User.Password, observed.User.Email, observed.User.Type).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(observed.User.ID)
		mock.ExpectQuery(regexp.QuoteMeta(querySelectUserIDByUsername)).
			WithArgs(observed.User.Username).
			WillReturnRows(rows)

		mock.
			ExpectExec(regexp.QuoteMeta(querySaveSchoolBus)).
			WithArgs(observed.SchoolBus.ID, observed.SchoolBus.LicensePlate, observed.SchoolBus.Model,
				observed.SchoolBus.Brand, observed.SchoolBus.License).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		rows = sqlmock.NewRows([]string{"id", "license_plate", "model", "brand", "school_bus_license", "created_at",
			"updated_at"}).
			AddRow(observed.User.ID, observed.SchoolBus.LicensePlate, observed.SchoolBus.Model,
				observed.SchoolBus.Brand, observed.SchoolBus.License, observed.SchoolBus.CreatedAt,
				observed.SchoolBus.UpdatedAt)
		mock.ExpectQuery(regexp.QuoteMeta(querySelectSchoolBusByID)).
			WithArgs(observed.SchoolBus.ID).
			WillReturnRows(rows)

		mock.
			ExpectExec(regexp.QuoteMeta(querySaveObservedUser)).
			WithArgs(observed.User.ID, observed.PrivacyKey, observed.CompanyName).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		rows = sqlmock.NewRows([]string{"user_id", "privacy_key", "company_name"}).
			AddRow(observed.User.ID, observed.PrivacyKey, observed.CompanyName)
		mock.ExpectQuery(regexp.QuoteMeta(queryGetObservedUserByID)).
			WithArgs(observed.User.ID).
			WillReturnRows(rows)

		mock.ExpectCommit()

		user, err := ur.SaveObservedUser(observed)
		assert.NotNil(t, user)
		assert.NoError(t, err)
		assert.Equal(t, observed.User.ID, user.User.ID)
	})

	t.Run("SaveObservedUser error saving user", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectExec(regexp.QuoteMeta(querySaveUser)).
			WithArgs(observed.User.Name, observed.User.LastName, observed.User.IDNumber, observed.User.Username,
				observed.User.Password, observed.User.Email, observed.User.Type).
			WillReturnError(errors.New("some error"))

		mock.ExpectCommit()

		user, err := ur.SaveObservedUser(observed)
		assert.Nil(t, user)
		assert.NotNil(t, err)
	})

	t.Run("SaveObservedUser error selecting user", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectExec(regexp.QuoteMeta(querySaveUser)).
			WithArgs(observed.User.Name, observed.User.LastName, observed.User.IDNumber, observed.User.Username,
				observed.User.Password, observed.User.Email, observed.User.Type).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		mock.ExpectQuery(regexp.QuoteMeta(querySelectUserIDByUsername)).
			WithArgs(observed.User.Username).
			WillReturnError(errors.New("some error"))

		mock.ExpectCommit()

		user, err := ur.SaveObservedUser(observed)
		assert.Nil(t, user)
		assert.NotNil(t, err)
	})

	t.Run("SaveObservedUser error saving school bus", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectExec(regexp.QuoteMeta(querySaveUser)).
			WithArgs(observed.User.Name, observed.User.LastName, observed.User.IDNumber, observed.User.Username,
				observed.User.Password, observed.User.Email, observed.User.Type).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(observed.User.ID)
		mock.ExpectQuery(regexp.QuoteMeta(querySelectUserIDByUsername)).
			WithArgs(observed.User.Username).
			WillReturnRows(rows)

		mock.
			ExpectExec(regexp.QuoteMeta(querySaveSchoolBus)).
			WithArgs(observed.SchoolBus.ID, observed.SchoolBus.LicensePlate, observed.SchoolBus.Model,
				observed.SchoolBus.Brand, observed.SchoolBus.License).
			WillReturnError(errors.New("some error"))

		mock.ExpectCommit()

		user, err := ur.SaveObservedUser(observed)
		assert.Nil(t, user)
		assert.NotNil(t, err)
	})

	t.Run("SaveObservedUser error selecting school bus", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectExec(regexp.QuoteMeta(querySaveUser)).
			WithArgs(observed.User.Name, observed.User.LastName, observed.User.IDNumber, observed.User.Username,
				observed.User.Password, observed.User.Email, observed.User.Type).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(observed.User.ID)
		mock.ExpectQuery(regexp.QuoteMeta(querySelectUserIDByUsername)).
			WithArgs(observed.User.Username).
			WillReturnRows(rows)

		mock.
			ExpectExec(regexp.QuoteMeta(querySaveSchoolBus)).
			WithArgs(observed.SchoolBus.ID, observed.SchoolBus.LicensePlate, observed.SchoolBus.Model,
				observed.SchoolBus.Brand, observed.SchoolBus.License).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		mock.ExpectQuery(regexp.QuoteMeta(querySelectSchoolBusByID)).
			WithArgs(observed.SchoolBus.ID).
			WillReturnError(errors.New("some error"))

		mock.ExpectCommit()

		user, err := ur.SaveObservedUser(observed)
		assert.Nil(t, user)
		assert.NotNil(t, err)
	})

	t.Run("SaveObservedUser error saving observed user", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectExec(regexp.QuoteMeta(querySaveUser)).
			WithArgs(observed.User.Name, observed.User.LastName, observed.User.IDNumber, observed.User.Username,
				observed.User.Password, observed.User.Email, observed.User.Type).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(observed.User.ID)
		mock.ExpectQuery(regexp.QuoteMeta(querySelectUserIDByUsername)).
			WithArgs(observed.User.Username).
			WillReturnRows(rows)

		mock.
			ExpectExec(regexp.QuoteMeta(querySaveSchoolBus)).
			WithArgs(observed.SchoolBus.ID, observed.SchoolBus.LicensePlate, observed.SchoolBus.Model,
				observed.SchoolBus.Brand, observed.SchoolBus.License).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		rows = sqlmock.NewRows([]string{"id", "license_plate", "model", "brand", "school_bus_license", "created_at",
			"updated_at"}).
			AddRow(observed.User.ID, observed.SchoolBus.LicensePlate, observed.SchoolBus.Model,
				observed.SchoolBus.Brand, observed.SchoolBus.License, observed.SchoolBus.CreatedAt,
				observed.SchoolBus.UpdatedAt)
		mock.ExpectQuery(regexp.QuoteMeta(querySelectSchoolBusByID)).
			WithArgs(observed.SchoolBus.ID).
			WillReturnRows(rows)

		mock.
			ExpectExec(regexp.QuoteMeta(querySaveObservedUser)).
			WithArgs(observed.User.ID, observed.PrivacyKey, observed.CompanyName).
			WillReturnError(errors.New("some error"))

		mock.ExpectCommit()

		user, err := ur.SaveObservedUser(observed)
		assert.Nil(t, user)
		assert.NotNil(t, err)
	})

	t.Run("SaveObservedUser error selecting observed user", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectExec(regexp.QuoteMeta(querySaveUser)).
			WithArgs(observed.User.Name, observed.User.LastName, observed.User.IDNumber, observed.User.Username,
				observed.User.Password, observed.User.Email, observed.User.Type).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		rows := sqlmock.NewRows([]string{"id"}).
			AddRow(observed.User.ID)
		mock.ExpectQuery(regexp.QuoteMeta(querySelectUserIDByUsername)).
			WithArgs(observed.User.Username).
			WillReturnRows(rows)

		mock.
			ExpectExec(regexp.QuoteMeta(querySaveSchoolBus)).
			WithArgs(observed.SchoolBus.ID, observed.SchoolBus.LicensePlate, observed.SchoolBus.Model,
				observed.SchoolBus.Brand, observed.SchoolBus.License).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		rows = sqlmock.NewRows([]string{"id", "license_plate", "model", "brand", "school_bus_license", "created_at",
			"updated_at"}).
			AddRow(observed.User.ID, observed.SchoolBus.LicensePlate, observed.SchoolBus.Model,
				observed.SchoolBus.Brand, observed.SchoolBus.License, observed.SchoolBus.CreatedAt,
				observed.SchoolBus.UpdatedAt)
		mock.ExpectQuery(regexp.QuoteMeta(querySelectSchoolBusByID)).
			WithArgs(observed.SchoolBus.ID).
			WillReturnRows(rows)

		mock.
			ExpectExec(regexp.QuoteMeta(querySaveObservedUser)).
			WithArgs(observed.User.ID, observed.PrivacyKey, observed.CompanyName).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		mock.ExpectQuery(regexp.QuoteMeta(queryGetObservedUserByID)).
			WithArgs(observed.User.ID).
			WillReturnError(errors.New("some error"))

		mock.ExpectCommit()

		user, err := ur.SaveObservedUser(observed)
		assert.Nil(t, user)
		assert.NotNil(t, err)
	})
}

func TestUpdateObservedUser(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
	}

	ur := NewUserRepository(gdb, context.Background())

	t.Run("UpdateObservedUser successful", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectExec(regexp.QuoteMeta(queryUpdateUser)).
			WithArgs(observed.User.Name, observed.User.LastName, observed.User.IDNumber, observed.User.Username,
				observed.User.Password, observed.User.Email, observed.User.Type, observed.User.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		mock.
			ExpectExec(regexp.QuoteMeta(queryUpdateObservedUser)).
			WithArgs(observed.PrivacyKey, observed.CompanyName, observed.SchoolBus.ID, observed.User.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		mock.ExpectCommit()

		user, err := ur.UpdateObservedUser(observed)
		assert.NotNil(t, user)
		assert.NoError(t, err)
		assert.Equal(t, observed.User.ID, user.User.ID)
	})

	t.Run("UpdateObservedUser error saving user", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectExec(regexp.QuoteMeta(queryUpdateUser)).
			WithArgs(observed.User.Name, observed.User.LastName, observed.User.IDNumber, observed.User.Username,
				observed.User.Password, observed.User.Email, observed.User.Type, observed.User.ID).
			WillReturnError(errors.New("some error"))

		mock.ExpectCommit()

		user, err := ur.UpdateObservedUser(observed)
		assert.Nil(t, user)
		assert.NotNil(t, err)
	})

	t.Run("UpdateObservedUser error updating observed user", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectExec(regexp.QuoteMeta(queryUpdateUser)).
			WithArgs(observed.User.Name, observed.User.LastName, observed.User.IDNumber, observed.User.Username,
				observed.User.Password, observed.User.Email, observed.User.Type, observed.User.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		mock.
			ExpectExec(regexp.QuoteMeta(queryUpdateObservedUser)).
			WithArgs(observed.PrivacyKey, observed.CompanyName, observed.SchoolBus.ID, observed.User.ID).
			WillReturnError(errors.New("some error"))

		mock.ExpectCommit()

		user, err := ur.UpdateObservedUser(observed)
		assert.Nil(t, user)
		assert.NotNil(t, err)
	})
}

func TestGetObservedUser(t *testing.T) {
	var (
		err      error
		expected = model.ObservedUser{
			User: model.User{
				ID:       1,
				Name:     "Juan",
				LastName: "Perez",
				IDNumber: "12345678",
				Username: "jperez",
				Password: "jperez1234",
				Email:    "jperez@mail.com",
				Enabled:  true,
				Type:     "observed",
			},
			PrivacyKey:  "juan.perez.12345678",
			CompanyName: "school bus company",
		}
		expectedObservedUser  = model.NewObservedUser(&expected)
		observedUser          model.IUser
		user                  = model.ObservedUser{User: expected.User}
		statementObservedUser = `
			SELECT u.id, u.name, u.last_name, u.id_number, u.username, u.password, u.email, u.type, u.enabled, 
			       ou.privacy_key, ou.company_name, sb.id, sb.license_plate, sb.model, sb.brand, sb.license, 
			       sb.created_at, sb.updated_at 
			FROM ObservedUsers AS ou 
			    INNER JOIN Users AS u ON u.id = ou.user_id 
			    INNER JOIN SchoolBuses AS sb ON sb.id = ou.school_bus_id 
			WHERE u.id = ?
			`
	)

	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{}) // open gorm db
	if err != nil {
		log.Error("error opening database connection")
		t.Fail()
	}

	ur := NewUserRepository(gdb, context.Background())

	t.Run("get observed user successful", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"user_id", "name", "last_name", "id_number", "username", "password", "email",
			"type", "enabled", "school_bus_id", "privacy_key", "company_name", "license_plate", "model", "brand",
			"license", "created_at", "updated_at"}).
			AddRow(expected.User.ID, expected.User.Name, expected.User.LastName, expected.User.IDNumber,
				expected.User.Username, expected.User.Password, expected.User.Email, expected.User.Type,
				expected.User.Enabled, expected.SchoolBus.ID, expected.PrivacyKey, expected.CompanyName,
				expected.SchoolBus.LicensePlate, expected.SchoolBus.Model, expected.SchoolBus.Brand,
				expected.SchoolBus.License, expected.SchoolBus.CreatedAt, expected.SchoolBus.UpdatedAt)
		mock.ExpectQuery(statementObservedUser).WithArgs(expected.User.ID).WillReturnRows(rows)

		observedUser, err = ur.GetObservedUser(user.GetUserID())

		assert.Nil(t, err)
		assert.Equal(t, expectedObservedUser, observedUser)
	})

	t.Run("get observed user scan error", func(t *testing.T) {
		mock.ExpectQuery(statementObservedUser).WithArgs(expected.User.ID).WillReturnError(web.ErrInternalServerError)

		observedUser, err = ur.GetObservedUser(user.GetUserID())

		assert.NotNil(t, err)
	})
}

func TestDeleteObservedUser(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
	}

	ur := NewUserRepository(gdb, context.Background())

	t.Run("delete observed user successful", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectExec(regexp.QuoteMeta(queryDeleteObservedUser)).
			WithArgs(observed.User.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		mock.
			ExpectExec(regexp.QuoteMeta(queryDeleteUser)).
			WithArgs(observed.User.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		mock.ExpectCommit()

		err := ur.DeleteObservedUser(observed.User.ID)
		assert.NoError(t, err)
	})

	t.Run("delete user error deleting observed user", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectExec(regexp.QuoteMeta(queryDeleteObservedUser)).
			WithArgs(observed.User.ID).
			WillReturnError(errors.New("some error"))

		mock.ExpectCommit()

		err = ur.DeleteObservedUser(observed.User.ID)
		assert.NotNil(t, err)
	})

	t.Run("SaveUser error deleting user", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectExec(regexp.QuoteMeta(queryDeleteObservedUser)).
			WithArgs(observed.User.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		mock.
			ExpectExec(regexp.QuoteMeta(queryDeleteUser)).
			WithArgs(observed.User.ID).
			WillReturnError(errors.New("some error"))

		mock.ExpectCommit()

		err := ur.DeleteObservedUser(observed.User.ID)
		assert.NotNil(t, err)
	})
}

var observer = model.ObserverUser{
	User: model.User{
		ID:        2,
		Name:      "Maria",
		LastName:  "Dominguez",
		IDNumber:  "87654321",
		Username:  "mdominguez",
		Password:  "mdominguez1234",
		Email:     "mdominguez@mail.com",
		Enabled:   true,
		Type:      "observer",
		CreatedAt: "2022-12-10 17:49:30",
		UpdatedAt: "2022-12-10 17:49:30",
	},
	Children: []model.Child{
		{
			ID:              1,
			ObserverUserID:  2,
			Name:            "Pilar",
			LastName:        "Dominguez",
			SchoolName:      "La Salle",
			SchoolStartTime: "08:00:00",
			SchoolEndTime:   "12:00:00",
			CreatedAt:       "2022-12-10 17:49:32",
			UpdatedAt:       "2022-12-10 17:49:32",
		},
	},
	ObservedUsers: []model.ObservedUser{
		{
			User:        model.User{ID: 1, Name: "Juan", LastName: "Perez", IDNumber: "12345678"},
			PrivacyKey:  "juan.perez.12345678",
			CompanyName: "company school bus",
			SchoolBus: model.SchoolBus{
				ID:           1,
				LicensePlate: "11AAA222",
				Model:        "Master",
				Brand:        "Renault",
				License:      "11222",
				CreatedAt:    "2022-12-10 17:49:30",
				UpdatedAt:    "2022-12-10 17:49:30",
			},
		},
	},
}

func TestSaveObserverUser(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
	}

	ur := NewUserRepository(gdb, context.Background())

	t.Run("SaveUser successful", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectExec(regexp.QuoteMeta(querySaveUser)).
			WithArgs(
				observer.User.Name, observer.User.LastName, observer.User.IDNumber, observer.User.Username,
				observer.User.Password, observer.User.Email, observer.User.Type,
			).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		rows := sqlmock.NewRows([]string{"id"}).AddRow(observer.User.ID)
		mock.ExpectQuery(querySelectUserIDByUsername).WillReturnRows(rows)

		mock.ExpectCommit()

		user, err := ur.SaveObserverUser(observer)
		assert.NotNil(t, user)
		assert.NoError(t, err)
		assert.Equal(t, user.User.ID, observer.User.ID)
	})

	t.Run("SaveUser error saving observer user", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectExec(regexp.QuoteMeta(querySaveUser)).
			WithArgs(observer.User.Name, observer.User.LastName, observer.User.IDNumber, observer.User.Username,
				observer.User.Password, observer.User.Email, observer.User.Type).
			WillReturnError(errors.New("some error"))

		mock.ExpectCommit()

		user, err := ur.SaveObserverUser(observer)
		assert.Nil(t, user)
		assert.NotNil(t, err)
	})

	t.Run("SaveUser error selecting observer user", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectExec(regexp.QuoteMeta(querySaveUser)).
			WithArgs(observer.User.Name, observer.User.LastName, observer.User.IDNumber, observer.User.Username,
				observer.User.Password, observer.User.Email, observer.User.Type).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		mock.ExpectQuery(querySelectUserIDByUsername).WillReturnError(errors.New("some error"))

		mock.ExpectCommit()

		user, err := ur.SaveObserverUser(observer)
		assert.Nil(t, user)
		assert.NotNil(t, err)
	})

}

func TestUpdateObserverUser(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
	}

	ur := NewUserRepository(gdb, context.Background())

	t.Run("UpdateObserverUser successful", func(t *testing.T) {
		mock.
			ExpectExec(regexp.QuoteMeta(queryUpdateUser)).
			WithArgs(
				observer.User.Name, observer.User.LastName, observer.User.IDNumber, observer.User.Username,
				observer.User.Password, observer.User.Email, observer.User.Type, observer.User.ID,
			).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		user, err := ur.UpdateObserverUser(observer)
		assert.NotNil(t, user)
		assert.NoError(t, err)
		assert.Equal(t, user.User.ID, observer.User.ID)
	})

	t.Run("UpdateObserverUser error saving observer user", func(t *testing.T) {
		mock.
			ExpectExec(regexp.QuoteMeta(querySaveUser)).
			WithArgs(observer.User.Name, observer.User.LastName, observer.User.IDNumber, observer.User.Username,
				observer.User.Password, observer.User.Email, observer.User.Type, observer.User.ID).
			WillReturnError(errors.New("some error"))

		user, err := ur.UpdateObserverUser(observer)
		assert.Nil(t, user)
		assert.NotNil(t, err)
	})
}

func TestDeleteObserverUser(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
	}

	ur := NewUserRepository(gdb, context.Background())

	t.Run("delete observer user successful", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectExec(regexp.QuoteMeta(queryDeleteObserverUser)).
			WithArgs(observer.User.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		mock.
			ExpectExec(regexp.QuoteMeta(queryDeleteUser)).
			WithArgs(observer.User.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		mock.ExpectCommit()

		err = ur.DeleteObserverUser(observer.User.ID)
		assert.NoError(t, err)
	})

	t.Run("delete user error deleting observer user", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectExec(regexp.QuoteMeta(queryDeleteObserverUser)).
			WithArgs(observer.User.ID).
			WillReturnError(errors.New("some error"))

		mock.ExpectCommit()

		err = ur.DeleteObserverUser(observer.User.ID)
		assert.NotNil(t, err)
	})

	t.Run("delete user error deleting user", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectExec(regexp.QuoteMeta(queryDeleteObserverUser)).
			WithArgs(observer.User.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		mock.
			ExpectExec(regexp.QuoteMeta(queryDeleteUser)).
			WithArgs(observer.User.ID).
			WillReturnError(errors.New("some error"))

		mock.ExpectCommit()

		err = ur.DeleteObserverUser(observer.User.ID)
		assert.NotNil(t, err)
	})
}

func TestGetObserverUser(t *testing.T) {
	var (
		err      error
		expected = model.ObserverUser{
			User: model.User{
				ID:       2,
				Name:     "Maria",
				LastName: "Dominguez",
				IDNumber: "87654321",
				Username: "mdominguez",
				Password: "mdominguez1234",
				Email:    "mdominguez@mail.com",
				Enabled:  true,
				Type:     "observer",
			},
			Children: []model.Child{
				{
					ID:              1,
					ObserverUserID:  2,
					Name:            "Pilar",
					LastName:        "Dominguez",
					SchoolName:      "La Salle",
					SchoolStartTime: "08:00:00",
					SchoolEndTime:   "12:00:00",
				},
			},
			ObservedUsers: []model.ObservedUser{
				{
					User:        model.User{ID: 1, Name: "Juan", LastName: "Perez", IDNumber: "12345678"},
					PrivacyKey:  "juan.perez.12345678",
					CompanyName: "company school bus",
					SchoolBus: model.SchoolBus{
						ID:           1,
						LicensePlate: "11AAA222",
						Model:        "Master",
						Brand:        "Renault",
						License:      "11222",
					},
				},
			},
		}
		expectedObserverUser                = model.NewObserverUser(&expected)
		observerUser                        model.IUser
		statementUser                       = fmt.Sprintf(queryGetUser, expected.User.ID)
		statementChildren                   = fmt.Sprintf(queryGetChildren, expected.User.ID)
		statementObservedUserOfObserverUser = fmt.Sprintf(queryGetObservedUserOfObserverUser, expected.User.ID)
		user                                = model.ObserverUser{User: expected.User}
	)
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
		t.Fail()
	}

	ur := NewUserRepository(gdb, context.Background())

	t.Run("get observer user, error scaning user", func(t *testing.T) {
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)

		mock.ExpectQuery(statementUser).WillReturnError(web.ErrInternalServerError)

		rows := sqlmock.NewRows([]string{"id", "name", "last_name", "school_name", "school_start_time",
			"school_end_time", "observer_user_id", "created_at", "updated_at"}).
			AddRow(expected.Children[0].ID, expected.Children[0].Name, expected.Children[0].LastName,
				expected.Children[0].SchoolName, expected.Children[0].SchoolStartTime,
				expected.Children[0].SchoolEndTime, expected.Children[0].ObserverUserID, expected.Children[0].CreatedAt,
				expected.Children[0].UpdatedAt)
		mock.ExpectQuery(statementChildren).WillReturnRows(rows)

		rows = sqlmock.NewRows([]string{"id", "name", "last_name", "id_number", "company_name", "privacy_key",
			"school_bus_id", "license_plate", "model", "brand", "school_bus_license", "created_at", "updated_at"}).
			AddRow(expected.ObservedUsers[0].User.ID, expected.ObservedUsers[0].User.Name,
				expected.ObservedUsers[0].User.LastName, expected.ObservedUsers[0].User.IDNumber,
				expected.ObservedUsers[0].CompanyName, expected.ObservedUsers[0].PrivacyKey,
				expected.ObservedUsers[0].SchoolBus.ID, expected.ObservedUsers[0].SchoolBus.LicensePlate,
				expected.ObservedUsers[0].SchoolBus.Model, expected.ObservedUsers[0].SchoolBus.Brand,
				expected.ObservedUsers[0].SchoolBus.License, expected.ObservedUsers[0].SchoolBus.CreatedAt,
				expected.ObservedUsers[0].SchoolBus.UpdatedAt)
		mock.ExpectQuery(statementObservedUserOfObserverUser).WillReturnRows(rows)

		observerUser, err := ur.GetObserverUser(user.GetUserID())

		assert.NotNil(t, err)
		assert.Equal(t, web.ErrInternalServerError, err)
		assert.Nil(t, observerUser)
	})

	t.Run("get observer user, error scaning children", func(t *testing.T) {
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)

		rows := sqlmock.NewRows([]string{"id", "name", "last_name", "id_number", "username", "password", "email",
			"type", "enabled"}).
			AddRow(expected.User.ID, expected.User.Name, expected.User.LastName, expected.User.IDNumber,
				expected.User.Username, expected.User.Password, expected.User.Email, expected.User.Type,
				expected.User.Enabled)
		mock.ExpectQuery(statementUser).WillReturnRows(rows)

		mock.ExpectQuery(statementChildren).WillReturnError(web.ErrInternalServerError)

		rows = sqlmock.NewRows([]string{"id", "name", "last_name", "id_number", "company_name", "privacy_key",
			"school_bus_id", "license_plate", "model", "brand", "school_bus_license", "created_at", "updated_at"}).
			AddRow(expected.ObservedUsers[0].User.ID, expected.ObservedUsers[0].User.Name,
				expected.ObservedUsers[0].User.LastName, expected.ObservedUsers[0].User.IDNumber,
				expected.ObservedUsers[0].CompanyName, expected.ObservedUsers[0].PrivacyKey,
				expected.ObservedUsers[0].SchoolBus.ID, expected.ObservedUsers[0].SchoolBus.LicensePlate,
				expected.ObservedUsers[0].SchoolBus.Model, expected.ObservedUsers[0].SchoolBus.Brand,
				expected.ObservedUsers[0].SchoolBus.License, expected.ObservedUsers[0].SchoolBus.CreatedAt,
				expected.ObservedUsers[0].SchoolBus.UpdatedAt)
		mock.ExpectQuery(statementObservedUserOfObserverUser).WillReturnRows(rows)

		observerUser, err := ur.GetObserverUser(user.GetUserID())

		assert.NotNil(t, err)
		assert.Equal(t, web.ErrInternalServerError, err)
		assert.Nil(t, observerUser)
	})

	t.Run("get observer user, error scaning observerd user of observer user", func(t *testing.T) {
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)

		rows := sqlmock.NewRows([]string{"id", "name", "last_name", "id_number", "username", "password", "email",
			"type", "enabled"}).
			AddRow(expected.User.ID, expected.User.Name, expected.User.LastName, expected.User.IDNumber,
				expected.User.Username, expected.User.Password, expected.User.Email, expected.User.Type,
				expected.User.Enabled)
		mock.ExpectQuery(statementUser).WillReturnRows(rows)

		rows = sqlmock.NewRows([]string{"id", "name", "last_name", "school_name", "school_start_time",
			"school_end_time", "observer_user_id", "created_at", "updated_at"}).
			AddRow(expected.Children[0].ID, expected.Children[0].Name, expected.Children[0].LastName,
				expected.Children[0].SchoolName, expected.Children[0].SchoolStartTime,
				expected.Children[0].SchoolEndTime, expected.Children[0].ObserverUserID, expected.Children[0].CreatedAt,
				expected.Children[0].UpdatedAt)
		mock.ExpectQuery(statementChildren).WillReturnRows(rows)

		mock.ExpectQuery(statementObservedUserOfObserverUser).WillReturnError(web.ErrInternalServerError)

		observerUser, err := ur.GetObserverUser(user.GetUserID())

		assert.NotNil(t, err)
		assert.Equal(t, web.ErrInternalServerError, err)
		assert.Nil(t, observerUser)
	})

	t.Run("get observer user successful", func(t *testing.T) {
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)

		rows := sqlmock.NewRows([]string{"id", "name", "last_name", "id_number", "username", "password", "email",
			"type", "enabled"}).
			AddRow(expected.User.ID, expected.User.Name, expected.User.LastName, expected.User.IDNumber,
				expected.User.Username, expected.User.Password, expected.User.Email, expected.User.Type,
				expected.User.Enabled)
		mock.ExpectQuery(statementUser).WillReturnRows(rows)

		rows = sqlmock.NewRows([]string{"id", "name", "last_name", "school_name", "school_start_time",
			"school_end_time", "observer_user_id", "created_at", "updated_at"}).
			AddRow(expected.Children[0].ID, expected.Children[0].Name, expected.Children[0].LastName,
				expected.Children[0].SchoolName, expected.Children[0].SchoolStartTime,
				expected.Children[0].SchoolEndTime, expected.Children[0].ObserverUserID, expected.Children[0].CreatedAt,
				expected.Children[0].UpdatedAt)
		mock.ExpectQuery(statementChildren).WillReturnRows(rows)

		rows = sqlmock.NewRows([]string{"id", "name", "last_name", "id_number", "company_name", "privacy_key",
			"school_bus_id", "license_plate", "model", "brand", "school_bus_license", "created_at", "updated_at"}).
			AddRow(expected.ObservedUsers[0].User.ID, expected.ObservedUsers[0].User.Name,
				expected.ObservedUsers[0].User.LastName, expected.ObservedUsers[0].User.IDNumber,
				expected.ObservedUsers[0].CompanyName, expected.ObservedUsers[0].PrivacyKey,
				expected.ObservedUsers[0].SchoolBus.ID, expected.ObservedUsers[0].SchoolBus.LicensePlate,
				expected.ObservedUsers[0].SchoolBus.Model, expected.ObservedUsers[0].SchoolBus.Brand,
				expected.ObservedUsers[0].SchoolBus.License, expected.ObservedUsers[0].SchoolBus.CreatedAt,
				expected.ObservedUsers[0].SchoolBus.UpdatedAt)
		mock.ExpectQuery(statementObservedUserOfObserverUser).WillReturnRows(rows)

		observerUser, err = ur.GetObserverUser(user.GetUserID())
		time.Sleep(200 * time.Millisecond)

		assert.Nil(t, err)
		assert.Equal(t, expectedObserverUser, observerUser)
	})
}

func TestFindObservedUserByPrivacyKey(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
	}

	ur := NewUserRepository(gdb, context.Background())

	t.Run("FindObservedUserByPrivacyKey successful", func(t *testing.T) {
		rows := sqlmock.
			NewRows(
				[]string{"id", "name", "last_name", "id_number", "username", "password", "email", "type", "enabled",
					"privacy_key", "company_name", "school_bus_id", "license_plate", "model", "brand", "license",
					"created_at", "updated_at"},
			).
			AddRow(
				observed.User.ID, observed.User.Name, observed.User.LastName, observed.User.IDNumber,
				observed.User.Username, observed.User.Password, observed.User.Email, observed.User.Type,
				observed.User.Enabled, observed.PrivacyKey, observed.CompanyName, observed.SchoolBus.ID,
				observed.SchoolBus.LicensePlate, observed.SchoolBus.Model, observed.SchoolBus.Brand,
				observed.SchoolBus.License, observed.SchoolBus.CreatedAt, observed.SchoolBus.UpdatedAt,
			)
		mock.
			ExpectQuery(regexp.QuoteMeta(queryGetObservedUserByPrivacyKey)).
			WithArgs(observed.PrivacyKey).
			WillReturnRows(rows)

		user, err := ur.FindObservedUserByPrivacyKey(observed.PrivacyKey)
		assert.NotNil(t, user)
		assert.NoError(t, err)
		assert.Equal(t, observed.PrivacyKey, user.PrivacyKey)
	})

	t.Run("FindObservedUserByPrivacyKey scan error", func(t *testing.T) {
		rows := sqlmock.
			NewRows(
				[]string{"id", "name", "last_name"},
			).
			AddRow(
				observed.User.ID, observed.User.Name, observed.User.LastName,
			)
		mock.
			ExpectQuery(regexp.QuoteMeta(queryGetObservedUserByPrivacyKey)).
			WithArgs(observed.PrivacyKey).
			WillReturnRows(rows)

		user, err := ur.FindObservedUserByPrivacyKey(observed.PrivacyKey)
		assert.Nil(t, user)
		assert.Error(t, err)
	})

	t.Run("FindObservedUserByPrivacyKey rows error", func(t *testing.T) {
		mock.
			ExpectQuery(regexp.QuoteMeta(queryGetObservedUserByPrivacyKey)).
			WithArgs(observed.PrivacyKey).
			WillReturnError(web.ErrInternalServerError)

		user, err := ur.FindObservedUserByPrivacyKey(observed.PrivacyKey)
		assert.Nil(t, user)
		assert.Error(t, err)
	})

	t.Run("FindObservedUserByPrivacyKey not found error", func(t *testing.T) {
		mock.
			ExpectQuery(regexp.QuoteMeta(queryGetObservedUserByPrivacyKey)).
			WithArgs(observed.PrivacyKey).
			WillReturnError(web.ErrNoRows)

		user, err := ur.FindObservedUserByPrivacyKey(observed.PrivacyKey)
		assert.Nil(t, user)
		assert.NoError(t, err)
	})

	t.Run("FindObservedUserByPrivacyKey nil rows error", func(t *testing.T) {
		rows := sqlmock.
			NewRows(
				[]string{"id", "name", "last_name", "id_number", "username", "password", "email", "type", "enabled",
					"privacy_key", "company_name", "school_bus_id", "license_plate", "model", "brand", "license",
					"created_at", "updated_at"},
			).
			AddRow(
				observed.User.ID, observed.User.Name, observed.User.LastName, observed.User.IDNumber,
				observed.User.Username, observed.User.Password, observed.User.Email, observed.User.Type,
				observed.User.Enabled, "invalid privacy key", observed.CompanyName, observed.SchoolBus.ID,
				observed.SchoolBus.LicensePlate, observed.SchoolBus.Model, observed.SchoolBus.Brand,
				observed.SchoolBus.License, observed.SchoolBus.CreatedAt, observed.SchoolBus.UpdatedAt,
			)
		mock.
			ExpectQuery(regexp.QuoteMeta(queryGetObservedUserByPrivacyKey)).
			WithArgs("not found").
			WillReturnRows(rows)

		user, err := ur.FindObservedUserByPrivacyKey(observed.PrivacyKey)
		assert.Nil(t, user)
		assert.Error(t, err)
	})

}

func TestSaveObservedUserInObserverUser(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
	}

	ur := NewUserRepository(gdb, context.Background())

	t.Run("save observed user in observer user successful", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectExec(regexp.QuoteMeta(querySaveObservedUserInObserverUser)).
			WithArgs(observed.User.ID, observer.User.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		rows := sqlmock.
			NewRows([]string{"observed_user_id", "observer_user_id"}).
			AddRow(observed.User.ID, observer.User.ID)
		mock.
			ExpectQuery(
				regexp.QuoteMeta(
					`SELECT * FROM ObservedUsersObserverUsers WHERE observed_user_id = ? AND observer_user_id = ?`,
				),
			).
			WithArgs(observed.User.ID, observer.User.ID).
			WillReturnRows(rows)

		mock.ExpectCommit()

		err := ur.SaveObservedUserInObserverUser(observed.User.ID, observer.User.ID)
		assert.NoError(t, err)
	})

	t.Run("save observed user in observer user error saving observed user", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectQuery(
				regexp.QuoteMeta(
					`SELECT * FROM ObservedUsersObserverUsers WHERE observed_user_id = ? AND observer_user_id = ?`,
				),
			).
			WithArgs(observed.User.ID, observer.User.ID).
			WillReturnError(errors.New("some error"))

		mock.ExpectCommit()

		err := ur.SaveObservedUserInObserverUser(observed.User.ID, observer.User.ID)
		assert.NotNil(t, err)
	})

	t.Run("save observed user in observer user error selecting observed and observer user", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectExec(regexp.QuoteMeta(querySaveObservedUserInObserverUser)).
			WithArgs(observed.User.ID, observer.User.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		mock.
			ExpectQuery(
				regexp.QuoteMeta(
					`SELECT * FROM ObservedUsersObserverUsers WHERE observed_user_id = ? AND observer_user_id = ?`,
				),
			).
			WithArgs(observed.User.ID, observer.User.ID).
			WillReturnError(errors.New("some error"))

		mock.ExpectCommit()

		err := ur.SaveObservedUserInObserverUser(observed.User.ID, observer.User.ID)
		assert.NotNil(t, err)
	})

}

func TestDeleteObservedUserInObserverUser(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
	}

	ur := NewUserRepository(gdb, context.Background())

	t.Run("delete observed user in observer user error saving observed user", func(t *testing.T) {
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)
		mock.
			ExpectQuery(
				regexp.QuoteMeta(
					`DELETE FROM ObservedUsersObserverUsers WHERE observed_user_id = ? AND observer_user_id = ?`,
				),
			).
			WithArgs(observed.User.ID, observer.User.ID).
			WillReturnError(errors.New("some error"))

		err := ur.DeleteObservedUserInObserverUser(observed.User.ID, observer.User.ID)
		assert.NotNil(t, err)
	})

	t.Run("successful delete observed user in observer user", func(t *testing.T) {
		mock.
			ExpectExec(
				regexp.QuoteMeta(
					`DELETE FROM ObservedUsersObserverUsers WHERE observed_user_id = ? AND observer_user_id = ?`,
				),
			).
			WithArgs(observed.User.ID, observer.User.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		err := ur.DeleteObservedUserInObserverUser(observed.User.ID, observer.User.ID)
		assert.Nil(t, err)
	})

}
