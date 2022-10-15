package repository

import (
	"context"
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gcoron/donde-estan-ws/internal/bussiness/model"
	"github.com/gcoron/donde-estan-ws/internal/bussiness/model/web"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var u = model.User{
	ID:       1,
	Name:     "user",
	LastName: "user",
	NumberID: "11100011",
	Username: "user",
	Password: "user",
	Email:    "user@user.com",
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
	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{}) // open gorm db
	ur := &UserRepository{gdb, context.Background()}

	query := "SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1"
	rows := sqlmock.NewRows([]string{"id", "name", "email", "username"}).
		AddRow(u.ID, u.Name, u.Email, u.Username)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(u.ID).WillReturnRows(rows)

	user, err := ur.Get(u.ID)
	assert.NotNil(t, user)
	assert.NoError(t, err)
	assert.Equal(t, u.ID, user.ID)
}

func TestGetUserError(t *testing.T) {
	db, mock := NewMock()
	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{}) // open gorm db
	ur := &UserRepository{gdb, context.Background()}

	query := "SELECT * FROM `invalid_table`"
	rows := sqlmock.NewRows([]string{"id", "name", "email", "username"}).AddRow(u.ID, u.Name, u.Email, u.Username)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(u.ID).WillReturnRows(rows)

	user, err := ur.Get(u.ID)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestGetUserNotFound(t *testing.T) {
	db, mock := NewMock()
	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{}) // open gorm db

	ur := &UserRepository{gdb, context.Background()}

	query := "SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(u.ID).WillReturnError(gorm.ErrRecordNotFound)

	user, err := ur.Get(u.ID)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestSaveUser(t *testing.T) {
	db, mock := NewMock()
	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{}) // open gorm db
	ur := &UserRepository{gdb, context.Background()}

	query := "INSERT INTO `users` (`name`,`last_name`,`number_id`,`username`,`password`,`email`,`id`) VALUES (?,?,?,?,?,?,?)" //"INSERT INTO `users` VALUES `users`\\.`id` = \\?"
	mock.ExpectBegin()
	mock.
		ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(u.Name, u.LastName, u.NumberID, u.Username, u.Password, u.Email, u.ID).
		WillReturnResult(sqlmock.NewResult(int64(1), 1))
	mock.ExpectCommit()

	user, err := ur.Save(u)
	assert.NotNil(t, user)
	assert.NoError(t, err)
	assert.Equal(t, u.ID, user.ID)
}

func TestFindByUsername(t *testing.T) {
	db, mock := NewMock()
	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{}) // open gorm db
	ur := &UserRepository{gdb, context.Background()}

	query := "SELECT * FROM USERS WHERE username = ?"
	rows := sqlmock.NewRows([]string{"id", "name", "lastname", "numberID", "username", "password", "email"}).
		AddRow(u.ID, u.Name, u.LastName, u.NumberID, u.Username, u.Password, u.Email)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(u.Username).WillReturnRows(rows)

	user, err := ur.FindByUsername(u.Username)
	assert.NotNil(t, user)
	assert.NoError(t, err)
	assert.Equal(t, u.Username, user.Username)
}

func TestFindByUsernameScanError(t *testing.T) {
	db, mock := NewMock()
	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{}) // open gorm db
	ur := &UserRepository{gdb, context.Background()}

	query := "SELECT * FROM USERS WHERE username = ?"
	rows := sqlmock.NewRows([]string{"id", "name", "lastname", "numberID", "username", "password"}).
		AddRow(u.ID, u.Name, u.LastName, u.NumberID, u.Username, u.Password)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(u.Username).WillReturnRows(rows)

	user, err := ur.FindByUsername(u.Username)
	assert.Nil(t, user)
	assert.Error(t, err)
}

func TestFindByUsernameRowsError(t *testing.T) {
	db, mock := NewMock()
	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{}) // open gorm db
	ur := &UserRepository{gdb, context.Background()}

	query := "SELECT * FROM USERS WHERE username = ?"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(u.Username).WillReturnError(web.ErrInternalServerError)

	user, err := ur.FindByUsername(u.Username)
	assert.Nil(t, user)
	assert.Error(t, err)
}

func TestFindByUsernameNilRows(t *testing.T) {
	db, mock := NewMock()
	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{}) // open gorm db
	ur := &UserRepository{gdb, context.Background()}

	query := "SELECT * FROM USERS WHERE username = ?"
	rows := sqlmock.NewRows([]string{"id", "name", "lastname", "numberID", "username", "password", "email"}).
		AddRow(u.ID, u.Name, u.LastName, u.NumberID, "invalid username", u.Password, u.Email)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("not found").WillReturnRows(rows)

	user, err := ur.FindByUsername(u.Username)
	assert.Nil(t, user)
	assert.Error(t, err)
}

func TestGetUsers(t *testing.T) {
	db, mock := NewMock()
	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{}) // open gorm db
	ur := &UserRepository{gdb, context.Background()}

	query := "SELECT * FROM USERS LIMIT ? OFFSET ?"
	rows := sqlmock.NewRows([]string{"id", "name", "lastname", "numberID", "username", "password", "email"}).
		AddRow(u.ID, u.Name, u.LastName, u.NumberID, u.Username, u.Password, u.Email).
		AddRow(u.ID, u.Name, u.LastName, u.NumberID, u.Username, u.Password, u.Email).
		AddRow(u.ID, u.Name, u.LastName, u.NumberID, u.Username, u.Password, u.Email).
		AddRow(u.ID, u.Name, u.LastName, u.NumberID, u.Username, u.Password, u.Email)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("1", "4").WillReturnRows(rows)

	users, err := ur.GetUsers("1", "4")
	assert.NotNil(t, users)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(*users))
}

func TestGetUsersScanError(t *testing.T) {
	db, mock := NewMock()
	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{}) // open gorm db
	ur := &UserRepository{gdb, context.Background()}

	query := "SELECT * FROM USERS LIMIT ? OFFSET ?"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("1", "4").WillReturnError(web.ErrInternalServerError)

	users, err := ur.GetUsers("1", "4")
	assert.Nil(t, users)
	assert.Error(t, err)
}

func TestGetUsersRowsError(t *testing.T) {
	db, mock := NewMock()
	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{}) // open gorm db
	ur := &UserRepository{gdb, context.Background()}

	query := "SELECT * FROM USERS LIMIT ? OFFSET ?"
	rows := sqlmock.NewRows([]string{"id", "name", "lastname", "numberID", "username", "password"}).
		AddRow(u.ID, u.Name, u.LastName, u.NumberID, u.Username, u.Password).
		AddRow(u.ID, u.Name, u.LastName, u.NumberID, u.Username, u.Password).
		AddRow(u.ID, u.Name, u.LastName, u.NumberID, u.Username, u.Password).
		AddRow(u.ID, u.Name, u.LastName, u.NumberID, u.Username, u.Password)
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("1", "4").WillReturnRows(rows)

	users, err := ur.GetUsers("1", "4")
	assert.Nil(t, *users)
	assert.Error(t, err)
}
