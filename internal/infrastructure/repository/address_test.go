package repository

import (
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var a = model.Address{
	ID:             1,
	Name:           "Casa",
	Street:         "25 de Mayo",
	Number:         "1010",
	Floor:          "1",
	Apartment:      "A",
	ZipCode:        "3000",
	City:           "Santa Fe",
	State:          "Santa Fe",
	Country:        "Argentina",
	Latitude:       "60.0000121",
	Longitude:      "-19.23423",
	CreatedAt:      "2023-02-18 17:09:33",
	UpdatedAt:      "2023-02-18 17:09:33",
	ObserverUserID: uint64(10),
}

func TestGetAddress(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
		t.Fail()
	}

	ar := NewAddressRepository(gdb, context.Background())

	t.Run("successful get address", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "street", "number", "floor", "apartament", "zipCode", "city", "state",
			"country", "latitude", "longitude", "created_at", "updated_at"}).
			AddRow(a.ID, a.Name, a.Street, a.Number, a.Floor, a.Apartment, a.ZipCode, a.City, a.State, a.Country,
				a.Latitude, a.Longitude, a.CreatedAt, a.UpdatedAt)
		mock.ExpectQuery(regexp.QuoteMeta(querySelectAddressByID)).WithArgs(a.ID).WillReturnRows(rows)

		address, err := ar.Get(a.ID)
		assert.NotNil(t, address)
		assert.NoError(t, err)
		assert.Equal(t, a.ID, address.ID)
	})

	t.Run("error getting address", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(querySelectAddressByID)).
			WithArgs(a.ID).
			WillReturnError(web.ErrInternalServerError)

		user, err := ar.Get(a.ID)
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("not found error getting address", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(querySelectAddressByID)).
			WithArgs(a.ID).
			WillReturnError(web.ErrNoRows)

		user, err := ar.Get(a.ID)
		assert.NoError(t, err)
		assert.Nil(t, user)
	})

}

func TestSaveAddress(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
		t.Fail()
	}

	ar := NewAddressRepository(gdb, context.Background())

	t.Run("successful save address", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)

		mock.ExpectExec(regexp.QuoteMeta(querySaveAddress)).
			WithArgs(a.Name, a.Street, a.Number, a.Floor, a.Apartment, a.ZipCode, a.City, a.State, a.Country,
				a.Latitude, a.Longitude, a.ObserverUserID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		rows := sqlmock.NewRows([]string{"id", "name", "street", "number", "floor", "apartament", "zipCode", "city", "state",
			"country", "latitude", "longitude", "created_at", "updated_at", "observer_user_id"}).
			AddRow(a.ID, a.Name, a.Street, a.Number, a.Floor, a.Apartment, a.ZipCode, a.City, a.State, a.Country,
				a.Latitude, a.Longitude, a.CreatedAt, a.UpdatedAt, a.ObserverUserID)
		mock.ExpectQuery(regexp.QuoteMeta(querySelectAddressByID)).WithArgs(a.ID).WillReturnRows(rows)

		mock.ExpectCommit()

		address, err := ar.Save(a)
		assert.NotNil(t, address)
		assert.NoError(t, err)
		assert.Equal(t, a.ID, address.ID)
	})

	t.Run("error saving address", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(querySaveAddress)).
			WithArgs(a.Name, a.Street, a.Number, a.Floor, a.Apartment, a.ZipCode, a.City, a.State, a.Country,
				a.Latitude, a.Longitude, a.ObserverUserID).
			WillReturnError(web.ErrInternalServerError)

		Address, err := ar.Save(a)
		assert.Error(t, err)
		assert.Nil(t, Address)
	})

	t.Run("error selecting address", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)

		mock.ExpectExec(regexp.QuoteMeta(querySaveAddress)).
			WithArgs(a.Name, a.Street, a.Number, a.Floor, a.Apartment, a.ZipCode, a.City, a.State, a.Country,
				a.Latitude, a.Longitude, a.ObserverUserID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		mock.ExpectQuery(regexp.QuoteMeta(querySelectAddressByID)).
			WithArgs(a.ID).
			WillReturnError(web.ErrInternalServerError)

		mock.ExpectCommit()

		Address, err := ar.Save(a)
		assert.Error(t, err)
		assert.Nil(t, Address)
	})

}

func TestUpdateAddress(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
		t.Fail()
	}

	ar := NewAddressRepository(gdb, context.Background())

	t.Run("successful update address", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)

		mock.ExpectExec(regexp.QuoteMeta(queryUpdateAddress)).
			WithArgs(a.Name, a.Street, a.Number, a.Floor, a.Apartment, a.ZipCode, a.City, a.State, a.Country,
				a.Latitude, a.Longitude, a.UpdatedAt, a.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		rows := sqlmock.NewRows([]string{"id", "name", "street", "number", "floor", "apartament", "zipCode", "city", "state",
			"country", "latitude", "longitude", "created_at", "updated_at", "observer_user_id"}).
			AddRow(a.ID, a.Name, a.Street, a.Number, a.Floor, a.Apartment, a.ZipCode, a.City, a.State, a.Country,
				a.Latitude, a.Longitude, a.CreatedAt, a.UpdatedAt, a.ObserverUserID)
		mock.ExpectQuery(regexp.QuoteMeta(querySelectAddressByID)).WithArgs(a.ID).WillReturnRows(rows)

		mock.ExpectCommit()

		address, err := ar.Update(a)
		assert.NotNil(t, address)
		assert.NoError(t, err)
		assert.Equal(t, a.ID, address.ID)
	})

	t.Run("error updating address", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(queryUpdateAddress)).
			WithArgs(a.Name, a.Street, a.Number, a.Floor, a.Apartment, a.ZipCode, a.City, a.State, a.Country,
				a.Latitude, a.Longitude, a.UpdatedAt, a.ID).
			WillReturnError(web.ErrInternalServerError)

		Address, err := ar.Update(a)
		assert.Error(t, err)
		assert.Nil(t, Address)
	})

	t.Run("error selecting address", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)

		mock.ExpectExec(regexp.QuoteMeta(queryUpdateAddress)).
			WithArgs(a.Name, a.Street, a.Number, a.Floor, a.Apartment, a.ZipCode, a.City, a.State, a.Country,
				a.Latitude, a.Longitude, a.UpdatedAt, a.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		mock.ExpectQuery(regexp.QuoteMeta(querySelectAddressByID)).
			WithArgs(a.ID).
			WillReturnError(web.ErrInternalServerError)

		mock.ExpectCommit()

		Address, err := ar.Update(a)
		assert.Error(t, err)
		assert.Nil(t, Address)
	})

}

func TestDeleteAddress(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
		t.Fail()
	}

	ar := NewAddressRepository(gdb, context.Background())

	t.Run("successful delete address", func(t *testing.T) {

		mock.ExpectExec(regexp.QuoteMeta(queryDeleteAddress)).
			WithArgs(a.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		err := ar.Delete(a.ID)
		assert.NoError(t, err)
	})

	t.Run("error deleting address", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(queryDeleteAddress)).
			WithArgs(a.ID).
			WillReturnError(web.ErrInternalServerError)

		err := ar.Delete(a.ID)
		assert.Error(t, err)
	})

}
