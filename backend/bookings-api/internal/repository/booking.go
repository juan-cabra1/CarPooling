package repository

import (
	"errors"

	"github.com/carpooling-ucc/bookings-api/internal/dao"
	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(booking *dao.BookingDAO) error
	FindByID(id string) (*dao.BookingDAO, error)
	FindByPassengerID(passengerID int64, page, limit int, status string) ([]*dao.BookingDAO, int64, error)
	FindByDriverID(driverID int64, page, limit int) ([]*dao.BookingDAO, error)
	FindByTripID(tripID string) ([]*dao.BookingDAO, error)
	Update(booking *dao.BookingDAO) error
	Delete(id string) error
	CheckScheduleConflict(passengerID int64, tripID string) (bool, error)
	CountByPassengerID(passengerID int64) (int64, error)
	CheckDuplicateBooking(tripID string, passengerID int64) (bool, error)
}

type bookingRepositoryImpl struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepositoryImpl{db: db}
}

func (r *bookingRepositoryImpl) Create(booking *dao.BookingDAO) error {
	return r.db.Create(booking).Error
}

func (r *bookingRepositoryImpl) FindByID(id string) (*dao.BookingDAO, error) {
	var booking dao.BookingDAO
	err := r.db.Where("id = ?", id).First(&booking).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepositoryImpl) FindByPassengerID(passengerID int64, page, limit int, status string) ([]*dao.BookingDAO, int64, error) {
	var bookings []*dao.BookingDAO
	var total int64

	query := r.db.Model(&dao.BookingDAO{}).Where("passenger_id = ?", passengerID)
	
	// Filter by status if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&bookings).Error
	if err != nil {
		return nil, 0, err
	}

	return bookings, total, nil
}

func (r *bookingRepositoryImpl) FindByDriverID(driverID int64, page, limit int) ([]*dao.BookingDAO, error) {
	var bookings []*dao.BookingDAO
	offset := (page - 1) * limit
	
	err := r.db.Where("driver_id = ?", driverID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&bookings).Error
	
	return bookings, err
}

func (r *bookingRepositoryImpl) FindByTripID(tripID string) ([]*dao.BookingDAO, error) {
	var bookings []*dao.BookingDAO
	err := r.db.Where("trip_id = ?", tripID).Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepositoryImpl) Update(booking *dao.BookingDAO) error {
	return r.db.Save(booking).Error
}

func (r *bookingRepositoryImpl) Delete(id string) error {
	return r.db.Delete(&dao.BookingDAO{}, "id = ?", id).Error
}

// CheckScheduleConflict checks if passenger has conflicting bookings
func (r *bookingRepositoryImpl) CheckScheduleConflict(passengerID int64, tripID string) (bool, error) {
	// For now, we'll implement a simple check
	// In a real scenario, we would need to get trip dates and check for overlaps
	var count int64
	err := r.db.Model(&dao.BookingDAO{}).
		Where("passenger_id = ? AND trip_id = ? AND status NOT IN (?)", 
			passengerID, tripID, []string{"cancelled", "completed"}).
		Count(&count).Error
	
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

func (r *bookingRepositoryImpl) CountByPassengerID(passengerID int64) (int64, error) {
	var count int64
	err := r.db.Model(&dao.BookingDAO{}).Where("passenger_id = ?", passengerID).Count(&count).Error
	return count, err
}

// CheckDuplicateBooking checks if passenger already has a booking for this trip
func (r *bookingRepositoryImpl) CheckDuplicateBooking(tripID string, passengerID int64) (bool, error) {
	var count int64
	err := r.db.Model(&dao.BookingDAO{}).
		Where("trip_id = ? AND passenger_id = ? AND status NOT IN (?)", 
			tripID, passengerID, []string{"cancelled"}).
		Count(&count).Error
	
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}
