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
		SELECT id, name, last_name, school_name, school_start_time, school_end_time, observer_user_id, created_at, 
		       updated_at,
		FROM Children 
		WHERE id = ?
	`

	querySaveChild = `
	INSERT INTO Children (name, last_name, school_name, school_start_time, school_end_time, observer_user_id) 
	VALUES (?, ?, ?, ?, ?, ?);
	`
	queryUpdateChild = `
		UPDATE Children SET name = ?, last_name = ?, school_name = ?, school_start_time = ?, school_end_time = ?, observer_user_id = ?, updated_at = ?
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

// Get obtains an address using ChildRepository by ID
func (a ChildRepository) Get(id uint64) (*model.Child, error) {
	var Child model.Child

	err := a.DB.
		Raw(querySelectChildByID, id).
		Row().
		Scan(
			&Child.ID, &Child.Name, &Child.LastName, &Child.SchoolName, &Child.SchoolStartTime, &Child.SchoolEndTime,
			&Child.ObserverUserID, &Child.CreatedAt, &Child.UpdatedAt,
		)

	if err != nil {
		if err.Error() == web.ErrNoRows.Error() {
			log.Error("error row scan not found selecting address")
			return nil, nil
		}
		log.Error("error row scan selecting address")
		return nil, err
	}

	return &Child, nil
}

// Save persists an address using ChildRepository.
func (a ChildRepository) Save(Child model.Child) (*model.Child, error) {
	tx := a.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Error("recovering for error in save address")
			tx.Rollback()
		}
	}()

	err := tx.Exec(
		querySaveChild,
		&Child.Name, &Child.LastName, &Child.SchoolName, &Child.SchoolStartTime, &Child.SchoolEndTime,
		&Child.ObserverUserID,
	).Error

	if err != nil {
		log.Error("error row scan saving address")
		return nil, err
	}

	err = tx.
		Raw(querySelectChildByID, Child.ID).
		Row().
		Scan(
			&Child.ID, &Child.Name, &Child.LastName, &Child.SchoolName, &Child.SchoolStartTime, &Child.SchoolEndTime,
			&Child.ObserverUserID, &Child.CreatedAt, &Child.UpdatedAt,
		)

	if err != nil {
		log.Error("error row scan selecting address")
		return nil, err
	}

	tx.Commit()
	return &Child, nil
}

// Update an address using ChildRepository by id
func (a ChildRepository) Update(Child model.Child) (*model.Child, error) {
	tx := a.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Error("recovering for error in save address")
			tx.Rollback()
		}
	}()

	err := tx.Exec(
		queryUpdateChild,
		&Child.Name, &Child.LastName, &Child.SchoolName, &Child.SchoolStartTime, &Child.SchoolEndTime,
		&Child.ObserverUserID, &Child.UpdatedAt, &Child.ID,
	).Error

	if err != nil {
		log.Error("error row scan saving address")
		return nil, err
	}

	err = tx.
		Raw(querySelectChildByID, Child.ID).
		Row().
		Scan(
			&Child.ID, &Child.Name, &Child.LastName, &Child.SchoolName, &Child.SchoolStartTime, &Child.SchoolEndTime,
			&Child.ObserverUserID, &Child.CreatedAt, &Child.UpdatedAt,
		)

	if err != nil {
		log.Error("error row scan selecting address")
		return nil, err
	}

	tx.Commit()
	return &Child, nil
}

// Delete an address using ChildRepository by id
func (a ChildRepository) Delete(id uint64) error {
	err := a.DB.Exec(queryDeleteChild, id).Error
	if err != nil {
		log.Error("error row scan saving address")
		return err
	}

	return nil
}
