package repository

import (
	"payments-api/internal/dao"

	"gorm.io/gorm"
)

// SellerRepository handles seller account database operations
type SellerRepository interface {
	Create(seller *dao.SellerAccount) error
	Update(seller *dao.SellerAccount) error
	FindByUserID(userID string) (*dao.SellerAccount, error)
	FindByUUID(uuid string) (*dao.SellerAccount, error)
	FindByMPUserID(mpUserID string) (*dao.SellerAccount, error)
	FindActiveByUserID(userID string) (*dao.SellerAccount, error)
	FindNeedingTokenRefresh() ([]dao.SellerAccount, error)
}

type sellerRepository struct {
	db *gorm.DB
}

// NewSellerRepository creates a new seller repository
func NewSellerRepository(db *gorm.DB) SellerRepository {
	return &sellerRepository{db: db}
}

// Create creates a new seller account
func (r *sellerRepository) Create(seller *dao.SellerAccount) error {
	return r.db.Create(seller).Error
}

// Update updates an existing seller account
func (r *sellerRepository) Update(seller *dao.SellerAccount) error {
	return r.db.Save(seller).Error
}

// FindByUserID finds a seller account by user ID
func (r *sellerRepository) FindByUserID(userID string) (*dao.SellerAccount, error) {
	var seller dao.SellerAccount
	err := r.db.Where("user_id = ?", userID).First(&seller).Error
	if err != nil {
		return nil, err
	}
	return &seller, nil
}

// FindByUUID finds a seller account by its UUID
func (r *sellerRepository) FindByUUID(uuid string) (*dao.SellerAccount, error) {
	var seller dao.SellerAccount
	err := r.db.Where("seller_uuid = ?", uuid).First(&seller).Error
	if err != nil {
		return nil, err
	}
	return &seller, nil
}

// FindByMPUserID finds a seller account by Mercado Pago user ID
func (r *sellerRepository) FindByMPUserID(mpUserID string) (*dao.SellerAccount, error) {
	var seller dao.SellerAccount
	err := r.db.Where("mp_user_id = ?", mpUserID).First(&seller).Error
	if err != nil {
		return nil, err
	}
	return &seller, nil
}

// FindActiveByUserID finds an active seller account by user ID
func (r *sellerRepository) FindActiveByUserID(userID string) (*dao.SellerAccount, error) {
	var seller dao.SellerAccount
	err := r.db.Where("user_id = ? AND status = ?", userID, dao.SellerStatusActive).First(&seller).Error
	if err != nil {
		return nil, err
	}
	return &seller, nil
}

// FindNeedingTokenRefresh finds seller accounts with expired or soon-to-expire tokens
func (r *sellerRepository) FindNeedingTokenRefresh() ([]dao.SellerAccount, error) {
	var sellers []dao.SellerAccount
	// Find active sellers whose token expires in the next 10 minutes
	err := r.db.Where(
		"status = ? AND token_expires_at IS NOT NULL AND token_expires_at < DATE_ADD(NOW(), INTERVAL 10 MINUTE)",
		dao.SellerStatusActive,
	).Find(&sellers).Error
	return sellers, err
}
