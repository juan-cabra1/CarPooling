package repository

import (
	"payments-api/internal/dao"

	"gorm.io/gorm"
)

// WithdrawalRepository handles withdrawal database operations
type WithdrawalRepository interface {
	Create(withdrawal *dao.Withdrawal) error
	Update(withdrawal *dao.Withdrawal) error
	FindByUUID(uuid string) (*dao.Withdrawal, error)
	FindByUserID(userID string, limit, offset int) ([]dao.Withdrawal, int64, error)
	FindPending() ([]dao.Withdrawal, error)
	HasPendingWithdrawal(userID string) (bool, error)
}

type withdrawalRepository struct {
	db *gorm.DB
}

// NewWithdrawalRepository creates a new withdrawal repository
func NewWithdrawalRepository(db *gorm.DB) WithdrawalRepository {
	return &withdrawalRepository{db: db}
}

// Create creates a new withdrawal
func (r *withdrawalRepository) Create(withdrawal *dao.Withdrawal) error {
	return r.db.Create(withdrawal).Error
}

// Update updates an existing withdrawal
func (r *withdrawalRepository) Update(withdrawal *dao.Withdrawal) error {
	return r.db.Save(withdrawal).Error
}

// FindByUUID finds a withdrawal by its UUID
func (r *withdrawalRepository) FindByUUID(uuid string) (*dao.Withdrawal, error) {
	var withdrawal dao.Withdrawal
	err := r.db.Where("withdrawal_uuid = ?", uuid).First(&withdrawal).Error
	if err != nil {
		return nil, err
	}
	return &withdrawal, nil
}

// FindByUserID finds withdrawals for a user with pagination
func (r *withdrawalRepository) FindByUserID(userID string, limit, offset int) ([]dao.Withdrawal, int64, error) {
	var withdrawals []dao.Withdrawal
	var total int64

	query := r.db.Model(&dao.Withdrawal{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&withdrawals).Error
	return withdrawals, total, err
}

// FindPending finds all pending withdrawals
func (r *withdrawalRepository) FindPending() ([]dao.Withdrawal, error) {
	var withdrawals []dao.Withdrawal
	err := r.db.Where("status = ?", dao.WithdrawalStatusPending).Find(&withdrawals).Error
	return withdrawals, err
}

// HasPendingWithdrawal checks if user has a pending withdrawal
func (r *withdrawalRepository) HasPendingWithdrawal(userID string) (bool, error) {
	var count int64
	err := r.db.Model(&dao.Withdrawal{}).
		Where("user_id = ? AND status IN (?, ?)",
			userID,
			dao.WithdrawalStatusPending,
			dao.WithdrawalStatusProcessing,
		).Count(&count).Error
	return count > 0, err
}
