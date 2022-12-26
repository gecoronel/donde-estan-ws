package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gcoron/donde-estan-ws/internal/bussiness/model"
	"github.com/gcoron/donde-estan-ws/internal/bussiness/model/web"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"testing"
	"time"
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

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{}) // open gorm db
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

func TestSaveUser(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{}) // open gorm db
	if err != nil {
		log.Error("error opening database connection")
	}

	ur := NewUserRepository(gdb, context.Background())

	query := "INSERT INTO `users` (`name`,`last_name`,`id_number`,`username`,`password`,`email`,`enabled`,`type`,`created_at`,`updated_at`,`id`) VALUES (?,?,?,?,?,?,?,?,?,?,?)" //"INSERT INTO `users` VALUES `users`\\.`id` = \\?"

	t.Run("SaveUser successful", func(t *testing.T) {
		mock.ExpectBegin()
		mock.
			ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(u.Name, u.LastName, u.IDNumber, u.Username, u.Password, u.Email, u.Enabled, u.Type, u.CreatedAt, u.UpdatedAt, u.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))
		mock.ExpectCommit()

		user, err := ur.Save(u)
		assert.NotNil(t, user)
		assert.NoError(t, err)
		assert.Equal(t, u.ID, user.ID)
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

	query := "SELECT * FROM USERS WHERE username = ?"

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
		rows := sqlmock.NewRows([]string{"id", "name", "lastname", "numberID", "username", "password", "email"}).
			AddRow(u.ID, u.Name, u.LastName, u.IDNumber, "invalid username", u.Password, u.Email)
		mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("not found").WillReturnRows(rows)

		user, err := ur.FindByUsername(u.Username)
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

func TestGetObserverUser(t *testing.T) {
	var (
		err      error
		expected = model.ObserverUser{
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
			Children: []model.Children{
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
						ID:               1,
						LicensePlate:     "11AAA222",
						Model:            "Master",
						Brand:            "Renault",
						SchoolBusLicense: "11222",
						CreatedAt:        "2022-12-10 17:49:30",
						UpdatedAt:        "2022-12-10 17:49:30",
					},
				},
			},
		}
		expectedObserverUser  = model.NewObserverUser(expected)
		observerUser          *model.IUser
		statementChildren     = "SELECT c.id, c.name, c.last_name, c.school_name, c.school_start_time, c.school_end_time, c.observer_user_id, c.created_at, c.updated_at FROM ObserverUsers AS oru INNER JOIN Children AS c ON  oru.user_id = c.observer_user_id;"
		statementObservedUser = "SELECT u.id, u.name, u.last_name, u.id_number, odu.company_name, odu.privacy_key, sb.id AS school_bus_id, sb.license_plate, sb.model, sb.brand, sb.school_bus_license, sb.created_at, sb.updated_at FROM ObserverUsers AS oru INNER JOIN ObservedUsers AS odu INNER JOIN ObservedUsersObserverUsers AS oduoru INNER JOIN Users AS u INNER JOIN SchoolBuses AS sb ON odu.user_id = oduoru.observed_user_id AND oru.user_id = oduoru.observer_user_id AND u.id = odu.user_id AND odu.school_bus_id = sb.id;"
		user                  = model.ObserverUser{User: expected.User}
	)
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{}) // open gorm db
	if err != nil {
		log.Error("error opening database connection")
		t.Fail()
	}

	ur := NewUserRepository(gdb, context.Background())

	t.Run("GetObserverUser children scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "last_name", "school_name", "school_start_time", "school_end_time", "observer_user_id", "created_at", "updated_at"}).
			AddRow(expected.Children[0].ID, expected.Children[0].Name, expected.Children[0].LastName, expected.Children[0].SchoolName, expected.Children[0].SchoolStartTime, expected.Children[0].SchoolEndTime, expected.Children[0].ObserverUserID, expected.Children[0].CreatedAt, expected.Children[0].UpdatedAt)
		mock.ExpectQuery(statementChildren).WillReturnError(web.ErrInternalServerError)

		rows = sqlmock.NewRows([]string{"id", "name", "last_name", "id_number", "company_name", "privacy_key", "school_bus_id", "license_plate", "model", "brand", "school_bus_license", "created_at", "updated_at"}).
			AddRow(expected.ObservedUsers[0].GetUserID(), expected.ObservedUsers[0].GetName(), expected.ObservedUsers[0].GetLastName(), expected.ObservedUsers[0].GetIDNumber(), expected.ObservedUsers[0].CompanyName, expected.ObservedUsers[0].PrivacyKey, expected.ObservedUsers[0].SchoolBus.ID, expected.ObservedUsers[0].SchoolBus.LicensePlate, expected.ObservedUsers[0].SchoolBus.Model, expected.ObservedUsers[0].SchoolBus.Brand, expected.ObservedUsers[0].SchoolBus.SchoolBusLicense, expected.ObservedUsers[0].SchoolBus.CreatedAt, expected.ObservedUsers[0].SchoolBus.UpdatedAt)
		mock.ExpectQuery(statementObservedUser).WillReturnRows(rows)

		observerUser, err = ur.GetObserverUser(&user)

		assert.NotNil(t, err)
	})

	t.Run("GetObserverUser observed user scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "last_name", "school_name", "school_start_time", "school_end_time", "observer_user_id", "created_at", "updated_at"}).
			AddRow(expected.Children[0].ID, expected.Children[0].Name, expected.Children[0].LastName, expected.Children[0].SchoolName, expected.Children[0].SchoolStartTime, expected.Children[0].SchoolEndTime, expected.Children[0].ObserverUserID, expected.Children[0].CreatedAt, expected.Children[0].UpdatedAt)
		mock.ExpectQuery(statementChildren).WillReturnRows(rows)

		rows = sqlmock.NewRows([]string{"id", "name", "last_name", "id_number", "company_name", "privacy_key", "school_bus_id", "license_plate", "model", "brand", "school_bus_license", "created_at", "updated_at"}).
			AddRow(expected.ObservedUsers[0].GetUserID(), expected.ObservedUsers[0].GetName(), expected.ObservedUsers[0].GetLastName(), expected.ObservedUsers[0].GetIDNumber(), expected.ObservedUsers[0].CompanyName, expected.ObservedUsers[0].PrivacyKey, expected.ObservedUsers[0].SchoolBus.ID, expected.ObservedUsers[0].SchoolBus.LicensePlate, expected.ObservedUsers[0].SchoolBus.Model, expected.ObservedUsers[0].SchoolBus.Brand, expected.ObservedUsers[0].SchoolBus.SchoolBusLicense, expected.ObservedUsers[0].SchoolBus.CreatedAt, expected.ObservedUsers[0].SchoolBus.UpdatedAt)
		mock.ExpectQuery(statementObservedUser).WillReturnError(web.ErrInternalServerError)

		observerUser, err = ur.GetObserverUser(&user)

		assert.NotNil(t, err)
	})

	t.Run("GetObserverUser successful", func(t *testing.T) {
		dbWithGoRoutine, mockWithGoRoutine := NewMock()
		defer dbWithGoRoutine.Close()

		gdbWithGoRoutine, errWithGoRoutine := gorm.Open(mysql.New(mysql.Config{Conn: dbWithGoRoutine, SkipInitializeWithVersion: true}), &gorm.Config{}) // open gorm db
		if errWithGoRoutine != nil {
			log.Error("error opening database connection")
			t.Fail()
		}

		urWithGoRoutine := NewUserRepository(gdbWithGoRoutine, context.Background())

		// note this line is important for unordered expectation matching
		mockWithGoRoutine.MatchExpectationsInOrder(false)

		rows := sqlmock.NewRows([]string{"id", "name", "last_name", "school_name", "school_start_time", "school_end_time", "observer_user_id", "created_at", "updated_at"}).
			AddRow(expected.Children[0].ID, expected.Children[0].Name, expected.Children[0].LastName, expected.Children[0].SchoolName, expected.Children[0].SchoolStartTime, expected.Children[0].SchoolEndTime, expected.Children[0].ObserverUserID, expected.Children[0].CreatedAt, expected.Children[0].UpdatedAt)
		mockWithGoRoutine.ExpectQuery(statementChildren).WillReturnRows(rows)

		rows = sqlmock.NewRows([]string{"id", "name", "last_name", "id_number", "company_name", "privacy_key", "school_bus_id", "license_plate", "model", "brand", "school_bus_license", "created_at", "updated_at"}).
			AddRow(expected.ObservedUsers[0].GetUserID(), expected.ObservedUsers[0].GetName(), expected.ObservedUsers[0].GetLastName(), expected.ObservedUsers[0].GetIDNumber(), expected.ObservedUsers[0].CompanyName, expected.ObservedUsers[0].PrivacyKey, expected.ObservedUsers[0].SchoolBus.ID, expected.ObservedUsers[0].SchoolBus.LicensePlate, expected.ObservedUsers[0].SchoolBus.Model, expected.ObservedUsers[0].SchoolBus.Brand, expected.ObservedUsers[0].SchoolBus.SchoolBusLicense, expected.ObservedUsers[0].SchoolBus.CreatedAt, expected.ObservedUsers[0].SchoolBus.UpdatedAt)
		mockWithGoRoutine.ExpectQuery(statementObservedUser).WillReturnRows(rows)

		observerUser, err = urWithGoRoutine.GetObserverUser(&user)
		time.Sleep(500 * time.Millisecond)

		assert.Nil(t, err)
		assert.Equal(t, expectedObserverUser, *observerUser)
	})
}

func TestGetObservedUser(t *testing.T) {
	var (
		err      error
		expected = model.ObservedUser{
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
		}
		expectedObservedUser  = model.NewObservedUser(expected)
		observerUser          *model.IUser
		statementObservedUser = "SELECT ou.user_id, ou.school_bus_id, ou.privacy_key, ou.company_name, sb.license_plate, sb.model, sb.brand, sb.school_bus_license, sb.created_at, sb.updated_at FROM ObservedUsers AS ou INNER JOIN SchoolBuses AS sb WHERE user_id = ?"
		user                  = model.ObservedUser{User: expected.User}
	)

	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{}) // open gorm db
	if err != nil {
		log.Error("error opening database connection")
		t.Fail()
	}

	ur := NewUserRepository(gdb, context.Background())

	t.Run("GetObservedUser successful", func(t *testing.T) {

		rows := sqlmock.NewRows([]string{"user_id", "school_bus_id", "privacy_key", "company_name", "license_plate", "model", "brand", "school_bus_license", "created_at", "updated_at"}).
			AddRow(expected.User.ID, expected.SchoolBus.ID, expected.PrivacyKey, expected.CompanyName, expected.SchoolBus.LicensePlate, expected.SchoolBus.Model, expected.SchoolBus.Brand, expected.SchoolBus.SchoolBusLicense, expected.SchoolBus.CreatedAt, expected.SchoolBus.UpdatedAt)
		mock.ExpectQuery(statementObservedUser).WithArgs(expected.User.ID).WillReturnRows(rows)

		observerUser, err = ur.GetObservedUser(&user)

		assert.Nil(t, err)
		assert.Equal(t, expectedObservedUser, *observerUser)
	})

	t.Run("GetObservedUser scan error", func(t *testing.T) {
		mock.ExpectQuery(statementObservedUser).WithArgs(expected.User.ID).WillReturnError(web.ErrInternalServerError)

		observerUser, err = ur.GetObservedUser(&user)

		assert.NotNil(t, err)
	})
}
