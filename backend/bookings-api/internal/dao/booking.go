package dao

import (
	"time"

	"gorm.io/gorm"
)

type BookingDAO struct {
	ID                 string         `gorm:"column:id;type:varchar(36);primaryKey"`
	TripID             string         `gorm:"column:trip_id;type:varchar(24);not null;index:idx_trip_id"`
	PassengerID        int64          `gorm:"column:passenger_id;type:bigint;not null;index:idx_passenger_id"`
	DriverID           int64          `gorm:"column:driver_id;type:bigint;not null;index:idx_driver_id"`
	SeatsReserved      int            `gorm:"column:seats_reserved;type:int;not null"`
	PricePerSeat       float64        `gorm:"column:price_per_seat;type:decimal(10,2);not null"`
	TotalAmount        float64        `gorm:"column:total_amount;type:decimal(10,2);not null"`
	Status             string         `gorm:"column:status;type:enum('pending','confirmed','completed','cancelled');default:'pending';not null"`
	PaymentStatus      string         `gorm:"column:payment_status;type:enum('pending','paid','refunded');default:'pending';not null"`
	ArrivedSafely      bool           `gorm:"column:arrived_safely;type:boolean;default:false;not null"`
	ArrivalConfirmedAt *time.Time     `gorm:"column:arrival_confirmed_at;type:timestamp;null"`
	CreatedAt          time.Time      `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

func (BookingDAO) TableName() string {
	return "reservations"
}

// BeforeCreate hook to set defaults
func (b *BookingDAO) BeforeCreate(tx *gorm.DB) error {
	if b.Status == "" {
		b.Status = "pending"
	}
	if b.PaymentStatus == "" {
		b.PaymentStatus = "pending"
	}
	return nil
}
