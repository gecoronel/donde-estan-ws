package repository

import (
	"context"
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"net/http"

	"github.com/gecoronel/donde-estan-ws/internal/bussiness/gateway"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model"
	"github.com/gecoronel/donde-estan-ws/internal/bussiness/model/web"
	log "github.com/sirupsen/logrus"
)

const (
	querySelectSchoolBusByID = "SELECT * FROM SchoolBuses WHERE id = ?"
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
	result := sb.DB.First(&schoolBus, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, web.NewError(http.StatusNotFound, "school_bus not found")
		}
		return nil, result.Error
	}

	return &schoolBus, nil
}

// Save persists a user using SchoolBusRepository.
func (sb SchoolBusRepository) Save(schoolBus model.SchoolBus) (*model.SchoolBus, error) {
	result := sb.DB.Model(&schoolBus).Create(&schoolBus)

	return &schoolBus, result.Error
}

// FindByID obtains a user using UserRepository by username
func (sb SchoolBusRepository) FindByID(id string) (*model.SchoolBus, error) {
	var schoolBus model.SchoolBus

	err := sb.DB.
		Raw("SELECT * FROM SchoolBuses WHERE id = @id", sql.Named("id", id)).
		Row().
		Scan(
			&schoolBus.ID, &schoolBus.Model, &schoolBus.Brand, &schoolBus.LicensePlate, &schoolBus.License,
			&schoolBus.CreatedAt, &schoolBus.UpdatedAt,
		)

	if err != nil {
		log.Error("error row scan")
		return nil, err
	}

	return &schoolBus, nil
}
