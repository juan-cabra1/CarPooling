package repository

import (
	"payments-api/internal/dao"

	"gorm.io/gorm"
)

// PaymentRepository handles payment database operations
type PaymentRepository interface {
	Create(payment *dao.Payment) error
	Update(payment *dao.Payment) error
	FindByUUID(uuid string) (*dao.Payment, error)
	FindByBookingID(bookingID string) (*dao.Payment, error)
	FindByMPPreferenceID(preferenceID string) (*dao.Payment, error)
	FindByMPPaymentID(paymentID string) (*dao.Payment, error)
	FindPendingReleases(olderThanHours int) ([]dao.Payment, error)
	FindByDriverID(driverID string, limit, offset int) ([]dao.Payment, int64, error)
	FindByPassengerID(passengerID string, limit, offset int) ([]dao.Payment, int64, error)
}

type paymentRepository struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new payment repository
func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

// Create creates a new payment
func (r *paymentRepository) Create(payment *dao.Payment) error {
	return r.db.Create(payment).Error
}

// Update updates an existing payment
func (r *paymentRepository) Update(payment *dao.Payment) error {
	return r.db.Save(payment).Error
}

// FindByUUID finds a payment by its UUID
func (r *paymentRepository) FindByUUID(uuid string) (*dao.Payment, error) {
	var payment dao.Payment
	err := r.db.Where("payment_uuid = ?", uuid).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// FindByBookingID finds a payment by booking ID
func (r *paymentRepository) FindByBookingID(bookingID string) (*dao.Payment, error) {
	var payment dao.Payment
	err := r.db.Where("booking_id = ?", bookingID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// FindByMPPreferenceID finds a payment by MP preference ID
func (r *paymentRepository) FindByMPPreferenceID(preferenceID string) (*dao.Payment, error) {
	var payment dao.Payment
	err := r.db.Where("mp_preference_id = ?", preferenceID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// FindByMPPaymentID finds a payment by MP payment ID
func (r *paymentRepository) FindByMPPaymentID(paymentID string) (*dao.Payment, error) {
	var payment dao.Payment
	err := r.db.Where("mp_payment_id = ?", paymentID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// FindPendingReleases finds payments that are approved but not released after X hours
func (r *paymentRepository) FindPendingReleases(olderThanHours int) ([]dao.Payment, error) {
	var payments []dao.Payment
	err := r.db.Where(
		"status = ? AND fund_status = ? AND approved_at < DATE_SUB(NOW(), INTERVAL ? HOUR)",
		dao.PaymentStatusApproved,
		dao.FundStatusHeld,
		olderThanHours,
	).Find(&payments).Error
	return payments, err
}

// FindByDriverID finds payments for a driver with pagination
func (r *paymentRepository) FindByDriverID(driverID string, limit, offset int) ([]dao.Payment, int64, error) {
	var payments []dao.Payment
	var total int64

	query := r.db.Model(&dao.Payment{}).Where("driver_id = ?", driverID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&payments).Error
	return payments, total, err
}

// FindByPassengerID finds payments for a passenger with pagination
func (r *paymentRepository) FindByPassengerID(passengerID string, limit, offset int) ([]dao.Payment, int64, error) {
	var payments []dao.Payment
	var total int64

	query := r.db.Model(&dao.Payment{}).Where("passenger_id = ?", passengerID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&payments).Error
	return payments, total, err
}
