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

var c = model.Child{
	ID:              1,
	Name:            "Pilar",
	LastName:        "Dominguez",
	SchoolName:      "La Salle",
	SchoolStartTime: "8:00",
	SchoolEndTime:   "12:00",
	CreatedAt:       "2023-02-18 17:09:33",
	UpdatedAt:       "2023-02-18 17:09:33",
	ObserverUserID:  uint64(10),
}

func TestGetChild(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
		t.Fail()
	}

	ar := NewChildRepository(gdb, context.Background())

	t.Run("successful get child", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "last_name", "school_name", "school_start_time",
			"school_end_time", "observer_user_id", "created_at", "updated_at"}).
			AddRow(c.ID, c.Name, c.LastName, c.SchoolName, c.SchoolStartTime, c.SchoolEndTime, c.ObserverUserID,
				c.CreatedAt, c.UpdatedAt)
		mock.ExpectQuery(regexp.QuoteMeta(querySelectChildByID)).WithArgs(c.ID).WillReturnRows(rows)

		child, err := ar.Get(c.ID)
		assert.NotNil(t, child)
		assert.NoError(t, err)
		assert.Equal(t, c.ID, child.ID)
	})

	t.Run("error getting child", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(querySelectChildByID)).
			WithArgs(c.ID).
			WillReturnError(web.ErrInternalServerError)

		user, err := ar.Get(c.ID)
		assert.Error(t, err)
		assert.Nil(t, user)
	})

	t.Run("not found error getting child", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(querySelectChildByID)).
			WithArgs(c.ID).
			WillReturnError(web.ErrNoRows)

		user, err := ar.Get(c.ID)
		assert.NoError(t, err)
		assert.Nil(t, user)
	})

}

func TestSaveChild(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
		t.Fail()
	}

	ar := NewChildRepository(gdb, context.Background())

	t.Run("successful save child", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)

		mock.ExpectExec(regexp.QuoteMeta(querySaveChild)).
			WithArgs(c.Name, c.LastName, c.SchoolName, c.SchoolStartTime, c.SchoolEndTime, c.ObserverUserID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		rows := sqlmock.NewRows([]string{"id"}).AddRow(c.ID)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT LAST_INSERT_ID();`)).WillReturnRows(rows)

		rows = sqlmock.NewRows([]string{"id", "name", "last_name", "school_name", "school_start_time",
			"school_end_time", "observer_user_id", "created_at", "updated_at"}).
			AddRow(c.ID, c.Name, c.LastName, c.SchoolName, c.SchoolStartTime, c.SchoolEndTime, c.ObserverUserID,
				c.CreatedAt, c.UpdatedAt)
		mock.ExpectQuery(regexp.QuoteMeta(querySelectChildByID)).WithArgs(c.ID).WillReturnRows(rows)

		mock.ExpectCommit()

		child, err := ar.Save(c)
		assert.NotNil(t, child)
		assert.NoError(t, err)
		assert.Equal(t, c.ID, child.ID)
	})

	t.Run("error saving child", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(querySaveChild)).
			WithArgs(c.Name, c.LastName, c.SchoolName, c.SchoolStartTime, c.SchoolEndTime, c.ObserverUserID).
			WillReturnError(web.ErrInternalServerError)

		child, err := ar.Save(c)
		assert.Error(t, err)
		assert.Nil(t, child)
	})

	t.Run("error selecting child id", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)

		mock.ExpectExec(regexp.QuoteMeta(querySaveChild)).
			WithArgs(c.Name, c.LastName, c.SchoolName, c.SchoolStartTime, c.SchoolEndTime, c.ObserverUserID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT LAST_INSERT_ID();`)).WillReturnError(web.ErrInternalServerError)

		mock.ExpectCommit()

		child, err := ar.Save(c)
		assert.Error(t, err)
		assert.Nil(t, child)
	})

	t.Run("error selecting child", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)

		mock.ExpectExec(regexp.QuoteMeta(querySaveChild)).
			WithArgs(c.Name, c.LastName, c.SchoolName, c.SchoolStartTime, c.SchoolEndTime, c.ObserverUserID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		rows := sqlmock.NewRows([]string{"id"}).AddRow(c.ID)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT LAST_INSERT_ID();`)).WillReturnRows(rows)

		mock.ExpectQuery(regexp.QuoteMeta(querySelectChildByID)).
			WithArgs(c.ID).
			WillReturnError(web.ErrInternalServerError)

		mock.ExpectCommit()

		child, err := ar.Save(c)
		assert.Error(t, err)
		assert.Nil(t, child)
	})

}

func TestUpdateChild(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
		t.Fail()
	}

	ar := NewChildRepository(gdb, context.Background())

	t.Run("successful update child", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)

		mock.ExpectExec(regexp.QuoteMeta(queryUpdateChild)).
			WithArgs(c.Name, c.LastName, c.SchoolName, c.SchoolStartTime, c.SchoolEndTime, c.UpdatedAt, c.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		rows := sqlmock.NewRows([]string{"id", "name", "last_name", "school_name", "school_start_time",
			"school_end_time", "observer_user_id", "created_at", "updated_at"}).
			AddRow(c.ID, c.Name, c.LastName, c.SchoolName, c.SchoolStartTime, c.SchoolEndTime, c.ObserverUserID,
				c.CreatedAt, c.UpdatedAt)
		mock.ExpectQuery(regexp.QuoteMeta(querySelectChildByID)).WithArgs(c.ID).WillReturnRows(rows)

		mock.ExpectCommit()

		child, err := ar.Update(c)
		assert.NotNil(t, child)
		assert.NoError(t, err)
		assert.Equal(t, c.ID, child.ID)
	})

	t.Run("error updating child", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(queryUpdateChild)).
			WithArgs(c.Name, c.LastName, c.SchoolName, c.SchoolStartTime, c.SchoolEndTime, c.ObserverUserID,
				c.UpdatedAt, c.ID).
			WillReturnError(web.ErrInternalServerError)

		Child, err := ar.Update(c)
		assert.Error(t, err)
		assert.Nil(t, Child)
	})

	t.Run("error selecting child", func(t *testing.T) {
		mock.ExpectBegin()
		// note this line is important for unordered expectation matching
		mock.MatchExpectationsInOrder(false)

		mock.ExpectExec(regexp.QuoteMeta(queryUpdateChild)).
			WithArgs(c.Name, c.LastName, c.SchoolName, c.SchoolStartTime, c.SchoolEndTime, c.UpdatedAt, c.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		mock.ExpectQuery(regexp.QuoteMeta(querySelectChildByID)).
			WithArgs(c.ID).
			WillReturnError(web.ErrInternalServerError)

		mock.ExpectCommit()

		child, err := ar.Update(c)
		assert.Error(t, err)
		assert.Nil(t, child)
	})

}

func TestDeleteChild(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	gdb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
	if err != nil {
		log.Error("error opening database connection")
		t.Fail()
	}

	ar := NewChildRepository(gdb, context.Background())

	t.Run("successful delete child", func(t *testing.T) {

		mock.ExpectExec(regexp.QuoteMeta(queryDeleteChild)).
			WithArgs(c.ID).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		err := ar.Delete(c.ID)
		assert.NoError(t, err)
	})

	t.Run("error deleting child", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(queryDeleteChild)).
			WithArgs(c.ID).
			WillReturnError(web.ErrInternalServerError)

		err := ar.Delete(c.ID)
		assert.Error(t, err)
	})

}
