package repository

import (
	"bookings-api/internal/dao"
	"time"

	"gorm.io/gorm"
)

// BookingRepository defines the interface for booking data access operations
type BookingRepository interface {
	// Create creates a new booking in the database
	Create(booking *dao.Booking) error

	// FindByID finds a booking by its UUID
	FindByID(id string) (*dao.Booking, error)

	// FindByPassengerID finds all bookings for a passenger with pagination
	// Returns bookings slice, total count, and error
	FindByPassengerID(passengerID int64, page, limit int) ([]dao.Booking, int64, error)

	// FindByTripID finds all bookings for a specific trip
	FindByTripID(tripID string) ([]dao.Booking, error)

	// Update updates an existing booking
	Update(booking *dao.Booking) error

	// UpdateStatus updates only the status of a booking
	UpdateStatus(bookingUUID string, status string) error

	// CancelBooking cancels a booking with a reason
	CancelBooking(bookingUUID string, reason string) error
}

// bookingRepository implements BookingRepository using GORM
type bookingRepository struct {
	db *gorm.DB
}

// NewBookingRepository creates a new instance of BookingRepository
func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

// Create creates a new booking in the database
func (r *bookingRepository) Create(booking *dao.Booking) error {
	return r.db.Create(booking).Error
}

// FindByID finds a booking by its UUID
func (r *bookingRepository) FindByID(id string) (*dao.Booking, error) {
	var booking dao.Booking
	err := r.db.Where("booking_uuid = ?", id).First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

// FindByPassengerID finds all bookings for a passenger with pagination
func (r *bookingRepository) FindByPassengerID(passengerID int64, page, limit int) ([]dao.Booking, int64, error) {
	var bookings []dao.Booking
	var total int64

	// Count total bookings for this passenger
	if err := r.db.Model(&dao.Booking{}).
		Where("passenger_id = ?", passengerID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset for pagination
	offset := (page - 1) * limit

	// Query bookings with pagination
	err := r.db.Where("passenger_id = ?", passengerID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bookings).Error

	if err != nil {
		return nil, 0, err
	}

	return bookings, total, nil
}

// FindByTripID finds all bookings for a specific trip
func (r *bookingRepository) FindByTripID(tripID string) ([]dao.Booking, error) {
	var bookings []dao.Booking
	err := r.db.Where("trip_id = ?", tripID).
		Order("created_at ASC").
		Find(&bookings).Error

	if err != nil {
		return nil, err
	}

	return bookings, nil
}

// Update updates an existing booking
func (r *bookingRepository) Update(booking *dao.Booking) error {
	return r.db.Save(booking).Error
}

// UpdateStatus updates only the status of a booking
func (r *bookingRepository) UpdateStatus(bookingUUID string, status string) error {
	return r.db.Model(&dao.Booking{}).
		Where("booking_uuid = ?", bookingUUID).
		Update("status", status).Error
}

// CancelBooking cancels a booking with a reason
func (r *bookingRepository) CancelBooking(bookingUUID string, reason string) error {
	now := time.Now()
	return r.db.Model(&dao.Booking{}).
		Where("booking_uuid = ?", bookingUUID).
		Updates(map[string]interface{}{
			"status":              dao.BookingStatusCancelled,
			"cancelled_at":        &now,
			"cancellation_reason": reason,
		}).Error
}
