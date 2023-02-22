package repository

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

var sb = model.SchoolBus{
	ID:           "1",
	LicensePlate: "11AAA22",
	Model:        "Master",
	Brand:        "Renault",
	License:      "111",
	CreatedAt:    "2023-02-18 17:09:33",
	UpdatedAt:    "2023-02-18 17:09:33",
}

func TestGetSchoolBus(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
		t.Fail()
	}

	ur := NewSchoolBusRepository(gdb, context.Background())

	t.Run("successful get school bus", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "license_plate", "model", "brand", "license", "created_at", "updated_at"}).
			AddRow(sb.ID, sb.LicensePlate, sb.Model, sb.Brand, sb.License, sb.CreatedAt, sb.UpdatedAt)
		mock.ExpectQuery(regexp.QuoteMeta(querySelectSchoolBusByID)).WithArgs(sb.ID).WillReturnRows(rows)

		schoolBus, err := ur.Get(sb.ID)
		assert.NotNil(t, schoolBus)
		assert.NoError(t, err)
		assert.Equal(t, sb.ID, schoolBus.ID)
	})

	t.Run("error getting school bus", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(querySelectSchoolBusByID)).
			WithArgs(sb.ID).
			WillReturnError(web.ErrInternalServerError)

		user, err := ur.Get(sb.ID)
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("not found error getting school bus", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(querySelectSchoolBusByID)).
			WithArgs(sb.ID).
			WillReturnError(web.ErrNoRows)

		user, err := ur.Get(sb.ID)
		assert.NoError(t, err)
		assert.Nil(t, user)
	})

}

func TestSave(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
		t.Fail()
	}

	ur := NewSchoolBusRepository(gdb, context.Background())

	t.Run("successful save school bus", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)

		mock.ExpectExec(regexp.QuoteMeta(querySaveSchoolBus)).
			WithArgs(sb.ID, sb.LicensePlate, sb.Model, sb.Brand, sb.License).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		rows := sqlmock.NewRows([]string{"id", "license_plate", "model", "brand", "license", "created_at", "updated_at"}).
			AddRow(sb.ID, sb.LicensePlate, sb.Model, sb.Brand, sb.License, sb.CreatedAt, sb.UpdatedAt)
		mock.ExpectQuery(regexp.QuoteMeta(querySelectSchoolBusByID)).WithArgs(sb.ID).WillReturnRows(rows)

		mock.ExpectCommit()

		schoolBus, err := ur.Save(sb)
		assert.NotNil(t, schoolBus)
		assert.NoError(t, err)
		assert.Equal(t, sb.ID, schoolBus.ID)
	})

	t.Run("error saving school bus", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(querySaveSchoolBus)).
			WithArgs(sb.ID, sb.LicensePlate, sb.Model, sb.Brand, sb.License).
			WillReturnError(web.ErrInternalServerError)

		schoolBus, err := ur.Save(sb)
		assert.Error(t, err)
		assert.Nil(t, schoolBus)
	})

	t.Run("error selecting school bus", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)

		mock.ExpectExec(regexp.QuoteMeta(querySaveSchoolBus)).
			WithArgs(sb.ID, sb.LicensePlate, sb.Model, sb.Brand, sb.License).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		mock.ExpectQuery(regexp.QuoteMeta(querySelectSchoolBusByID)).
			WithArgs(sb.ID).
			WillReturnError(web.ErrInternalServerError)

		mock.ExpectCommit()

		schoolBus, err := ur.Save(sb)
		assert.Error(t, err)
		assert.Nil(t, schoolBus)
	})

}

func TestUpdate(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
		t.Fail()
	}

	ur := NewSchoolBusRepository(gdb, context.Background())

	t.Run("successful update school bus", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)

		mock.ExpectExec(regexp.QuoteMeta(queryUpdateSchoolBus)).
			WithArgs(sb.ID, sb.LicensePlate, sb.Model, sb.Brand, sb.License, sb.UpdatedAt, sb.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		rows := sqlmock.NewRows([]string{"id", "license_plate", "model", "brand", "license", "created_at", "updated_at"}).
			AddRow(sb.ID, sb.LicensePlate, sb.Model, sb.Brand, sb.License, sb.CreatedAt, sb.UpdatedAt)
		mock.ExpectQuery(regexp.QuoteMeta(querySelectSchoolBusByID)).WithArgs(sb.ID).WillReturnRows(rows)

		mock.ExpectCommit()

		schoolBus, err := ur.Update(sb)
		assert.NotNil(t, schoolBus)
		assert.NoError(t, err)
		assert.Equal(t, sb.ID, schoolBus.ID)
	})

	t.Run("error updating school bus", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(queryUpdateSchoolBus)).
			WithArgs(sb.ID, sb.LicensePlate, sb.Model, sb.Brand, sb.License, sb.UpdatedAt, sb.ID).
			WillReturnError(web.ErrInternalServerError)

		schoolBus, err := ur.Update(sb)
		assert.Error(t, err)
		assert.Nil(t, schoolBus)
	})

	t.Run("error selecting school bus", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)

		mock.ExpectExec(regexp.QuoteMeta(queryUpdateSchoolBus)).
			WithArgs(sb.ID, sb.LicensePlate, sb.Model, sb.Brand, sb.License, sb.UpdatedAt, sb.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		mock.ExpectQuery(regexp.QuoteMeta(querySelectSchoolBusByID)).
			WithArgs(sb.ID).
			WillReturnError(web.ErrInternalServerError)

		mock.ExpectCommit()

		schoolBus, err := ur.Update(sb)
		assert.Error(t, err)
		assert.Nil(t, schoolBus)
	})

}

func TestDelete(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
		t.Fail()
	}

	ur := NewSchoolBusRepository(gdb, context.Background())

	t.Run("successful delete school bus", func(t *testing.T) {

		mock.ExpectExec(regexp.QuoteMeta(queryDeleteSchoolBus)).
			WithArgs(sb.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		err := ur.Delete(sb.ID)
		assert.NoError(t, err)
	})

	t.Run("error deleting school bus", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(querySaveSchoolBus)).
			WithArgs(sb.ID).
			WillReturnError(web.ErrInternalServerError)

		err := ur.Delete(sb.ID)
		assert.Error(t, err)
	})

}
