package repository

import (
	"context"
	"gorm.io/gorm"

	"github.com/gecoronel/donde-estan-ws/internal/bussiness/gateway"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
	log "github.com/sirupsen/logrus"
)

const (
	querySelectChildByID = `
		SELECT id, name, last_name, school_name, school_start_time, school_end_time, observer_user_id, created_at, updated_at
		FROM Children 
		WHERE id = ?
	`

	querySaveChild = `
	INSERT INTO Children (name, last_name, school_name, school_start_time, school_end_time, observer_user_id) 
	VALUES (?, ?, ?, ?, ?, ?);
	`
	queryUpdateChild = `
		UPDATE Children SET name = ?, last_name = ?, school_name = ?, school_start_time = ?, school_end_time = ?, updated_at = ?
		WHERE id = ?;
	`
	queryDeleteChild = `
	DELETE FROM Children WHERE id = ?;
	`
)

func NewChildRepository(db *gorm.DB, ctx context.Context) gateway.ChildRepository {
	return &ChildRepository{
		DB:      db,
		context: ctx,
	}
}

// ChildRepository represents the main repository for manage user
type ChildRepository struct {
	DB      *gorm.DB
	context context.Context
}

// Get obtains an child using ChildRepository by ID
func (a ChildRepository) Get(id uint64) (*model.Child, error) {
	var child model.Child

	err := a.DB.
		Raw(querySelectChildByID, id).
		Row().
		Scan(
			&child.ID, &child.Name, &child.LastName, &child.SchoolName, &child.SchoolStartTime, &child.SchoolEndTime,
			&child.ObserverUserID, &child.CreatedAt, &child.UpdatedAt,
		)

	if err != nil {
		if err.Error() == web.ErrNoRows.Error() {
			log.Error("error row scan not found selecting address")
			return nil, nil
		}
		log.Error("error row scan selecting address")
		return nil, err
	}

	return &child, nil
}

// Save persists an child using ChildRepository.
func (a ChildRepository) Save(child model.Child) (*model.Child, error) {
	tx := a.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Error("recovering for error in save address")
			tx.Rollback()
		}
	}()

	err := tx.Exec(
		querySaveChild,
		&child.Name, &child.LastName, &child.SchoolName, &child.SchoolStartTime, &child.SchoolEndTime,
		&child.ObserverUserID,
	).Error

	if err != nil {
		log.Error("error row scan saving child")
		return nil, err
	}

	err = tx.Raw(`SELECT LAST_INSERT_ID();`).Row().Scan(&child.ID)
	if err != nil {
		log.Error("error selecting child id")
		tx.Rollback()
		return nil, err
	}

	err = tx.
		Raw(querySelectChildByID, child.ID).
		Row().
		Scan(
			&child.ID, &child.Name, &child.LastName, &child.SchoolName, &child.SchoolStartTime, &child.SchoolEndTime,
			&child.ObserverUserID, &child.CreatedAt, &child.UpdatedAt,
		)
	if err != nil {
		log.Error("error row scan selecting child")
		return nil, err
	}

	tx.Commit()
	return &child, nil
}

// Update an child using ChildRepository by id
func (a ChildRepository) Update(child model.Child) (*model.Child, error) {
	tx := a.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Error("recovering for error in save child")
			tx.Rollback()
		}
	}()

	err := tx.Exec(
		queryUpdateChild,
		&child.Name, &child.LastName, &child.SchoolName, &child.SchoolStartTime, child.SchoolEndTime,
		&child.UpdatedAt, &child.ID,
	).Error

	if err != nil {
		log.Error("error updating child in repository")
		return nil, err
	}

	err = tx.
		Raw(querySelectChildByID, child.ID).
		Row().
		Scan(
			&child.ID, &child.Name, &child.LastName, &child.SchoolName, &child.SchoolStartTime, &child.SchoolEndTime,
			&child.ObserverUserID, &child.CreatedAt, &child.UpdatedAt,
		)
	if err != nil {
		log.Error("error selecting child")
		return nil, err
	}

	tx.Commit()
	return &child, nil
}

// Delete an child using ChildRepository by id
func (a ChildRepository) Delete(id uint64) error {
	err := a.DB.Exec(queryDeleteChild, id).Error
	if err != nil {
		log.Error("error deleting child")
		return err
	}

	return nil
}
