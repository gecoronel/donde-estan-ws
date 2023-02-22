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
	querySelectSchoolBusByID = `
		SELECT id, license_plate, model, brand, license, created_at, updated_at 
		FROM SchoolBuses 
		WHERE id = ?
	`

	querySaveSchoolBus   = `INSERT INTO SchoolBuses (id, license_plate, model, brand, license) VALUES (?, ?, ?, ?, ?);`
	queryUpdateSchoolBus = `
		UPDATE SchoolBuses SET id = ?, license_plate = ?, model = ?, brand = ?, license = ?, updated_at = ? 
		WHERE id = ?;
	`
	queryDeleteSchoolBus = `
	DELETE FROM SchoolBuses WHERE id = ?;
	`
)

func NewSchoolBusRepository(db *gorm.DB, ctx context.Context) gateway.SchoolBusRepository {
	return &SchoolBusRepository{
		DB:      db,
		context: ctx,
	}
}

// SchoolBusRepository represents the main repository for manage user
type SchoolBusRepository struct {
	DB      *gorm.DB
	context context.Context
}

// Get obtains a user using UserRepository by ID
func (sb SchoolBusRepository) Get(id string) (*model.SchoolBus, error) {
	var schoolBus model.SchoolBus

	err := sb.DB.
		Raw(querySelectSchoolBusByID, id).
		Row().
		Scan(
			&schoolBus.ID, &schoolBus.Model, &schoolBus.Brand, &schoolBus.LicensePlate, &schoolBus.License,
			&schoolBus.CreatedAt, &schoolBus.UpdatedAt,
		)

	if err != nil {
		if err.Error() == web.ErrNoRows.Error() {
			log.Error("error row scan not found selecting school bus")
			return nil, nil
		}
		log.Error("error row scan selecting school bus")
		return nil, err
	}

	return &schoolBus, nil
}

// Save persists a user using SchoolBusRepository.
func (sb SchoolBusRepository) Save(schoolBus model.SchoolBus) (*model.SchoolBus, error) {
	tx := sb.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Error("recovering for error in save observed user")
			tx.Rollback()
		}
	}()

	err := tx.Exec(
		querySaveSchoolBus,
		schoolBus.ID,
		schoolBus.LicensePlate,
		schoolBus.Model,
		schoolBus.Brand,
		schoolBus.License,
	).Error

	if err != nil {
		log.Error("error row scan saving school bus")
		return nil, err
	}

	err = tx.
		Raw(querySelectSchoolBusByID, schoolBus.ID).
		Row().
		Scan(
			&schoolBus.ID, &schoolBus.LicensePlate, &schoolBus.Model, &schoolBus.Brand, &schoolBus.License,
			&schoolBus.CreatedAt, &schoolBus.UpdatedAt,
		)

	if err != nil {
		log.Error("error row scan selecting school bus")
		return nil, err
	}

	tx.Commit()
	return &schoolBus, nil
}

// Update a school bus using SchoolBusRepository by id
func (sb SchoolBusRepository) Update(schoolBus model.SchoolBus) (*model.SchoolBus, error) {
	tx := sb.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Error("recovering for error in save observed user")
			tx.Rollback()
		}
	}()

	err := tx.Exec(
		queryUpdateSchoolBus,
		schoolBus.ID,
		schoolBus.LicensePlate,
		schoolBus.Model,
		schoolBus.Brand,
		schoolBus.License,
		schoolBus.UpdatedAt,
		schoolBus.ID,
	).Error

	if err != nil {
		log.Error("error row scan saving school bus")
		return nil, err
	}

	err = tx.
		Raw(querySelectSchoolBusByID, schoolBus.ID).
		Row().
		Scan(
			&schoolBus.ID, &schoolBus.LicensePlate, &schoolBus.Model, &schoolBus.Brand, &schoolBus.License,
			&schoolBus.CreatedAt, &schoolBus.UpdatedAt,
		)

	if err != nil {
		log.Error("error row scan selecting school bus")
		return nil, err
	}

	tx.Commit()
	return &schoolBus, nil
}

// Delete a school bus using SchoolBusRepository by id
func (sb SchoolBusRepository) Delete(id string) error {
	err := sb.DB.Exec(queryDeleteSchoolBus, id).Error
	if err != nil {
		log.Error("error row scan saving school bus")
		return err
	}

	return nil
}
