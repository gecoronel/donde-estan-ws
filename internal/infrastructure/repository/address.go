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
	querySelectAddressByID = `
		SELECT id, name, street, number, floor, apartment, zipCode, city, state, country, latitude, longitude, 
		       created_at, updated_at, observer_user_id
		FROM Addresses 
		WHERE id = ?
	`

	querySaveAddress = `
	INSERT INTO Addresses (name, street, number, floor, apartment, zipCode, city, state, country, latitude, longitude, observer_user_id) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	queryUpdateAddress = `
		UPDATE Addresses SET name = ?, street = ?, number = ?, floor = ?, apartment = ?, zipCode = ?, city = ?, state = ?, country = ?, latitude = ?, longitude = ?, updated_at = ?
		WHERE id = ?;
	`
	queryDeleteAddress = `
	DELETE FROM Addresses WHERE id = ?;
	`
)

func NewAddressRepository(db *gorm.DB, ctx context.Context) gateway.AddressRepository {
	return &AddressRepository{
		DB:      db,
		context: ctx,
	}
}

// AddressRepository represents the main repository for manage user
type AddressRepository struct {
	DB      *gorm.DB
	context context.Context
}

// Get obtains an address using AddressRepository by ID
func (a AddressRepository) Get(id uint64) (*model.Address, error) {
	var address model.Address

	err := a.DB.
		Raw(querySelectAddressByID, id).
		Row().
		Scan(
			&address.ID, &address.Name, &address.Street, &address.Number, &address.Floor, &address.Apartment,
			&address.ZipCode, &address.City, &address.State, &address.Country, &address.Latitude, &address.Longitude,
			&address.CreatedAt, &address.UpdatedAt, &address.ObserverUserID,
		)

	if err != nil {
		if err.Error() == web.ErrNoRows.Error() {
			log.Error("error row scan not found selecting address")
			return nil, nil
		}
		log.Error("error row scan selecting address")
		return nil, err
	}

	return &address, nil
}

// Save persists an address using AddressRepository.
func (a AddressRepository) Save(address model.Address) (*model.Address, error) {
	tx := a.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Error("recovering for error in save address")
			tx.Rollback()
		}
	}()

	err := tx.Exec(
		querySaveAddress,
		&address.Name, &address.Street, &address.Number, &address.Floor, &address.Apartment,
		&address.ZipCode, &address.City, &address.State, &address.Country, &address.Latitude, &address.Longitude,
		&address.ObserverUserID,
	).Error
	if err != nil {
		log.Error("error row scan saving address")
		return nil, err
	}

	err = tx.Raw(`SELECT LAST_INSERT_ID();`).Row().Scan(&address.ID)
	if err != nil {
		log.Error("error row scan selecting address")
		return nil, err
	}

	err = tx.
		Raw(querySelectAddressByID, address.ID).
		Row().
		Scan(
			&address.ID, &address.Name, &address.Street, &address.Number, &address.Floor, &address.Apartment,
			&address.ZipCode, &address.City, &address.State, &address.Country, &address.Latitude, &address.Longitude,
			&address.CreatedAt, &address.UpdatedAt, &address.ObserverUserID,
		)
	if err != nil {
		log.Error("error row scan selecting address")
		return nil, err
	}

	tx.Commit()
	return &address, nil
}

// Update an address using AddressRepository by id
func (a AddressRepository) Update(address model.Address) (*model.Address, error) {
	tx := a.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Error("recovering for error in save address")
			tx.Rollback()
		}
	}()

	err := tx.Exec(
		queryUpdateAddress,
		&address.Name, &address.Street, &address.Number, &address.Floor, &address.Apartment, &address.ZipCode,
		&address.City, &address.State, &address.Country, &address.Latitude, &address.Longitude, &address.UpdatedAt,
		&address.ID,
	).Error

	if err != nil {
		log.Error("error row scan saving address")
		return nil, err
	}

	err = tx.
		Raw(querySelectAddressByID, address.ID).
		Row().
		Scan(
			&address.ID, &address.Name, &address.Street, &address.Number, &address.Floor, &address.Apartment,
			&address.ZipCode, &address.City, &address.State, &address.Country, &address.Latitude, &address.Longitude,
			&address.CreatedAt, &address.UpdatedAt, &address.ObserverUserID,
		)

	if err != nil {
		log.Error("error row scan selecting address")
		return nil, err
	}

	tx.Commit()
	return &address, nil
}

// Delete an address using AddressRepository by id
func (a AddressRepository) Delete(id uint64) error {
	err := a.DB.Exec(queryDeleteAddress, id).Error
	if err != nil {
		log.Error("error row scan saving address")
		return err
	}

	return nil
}
