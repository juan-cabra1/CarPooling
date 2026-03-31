package repository

import (
	"payments-api/internal/dao"

	"gorm.io/gorm"
)

// WalletRepository handles wallet database operations
type WalletRepository interface {
	Create(wallet *dao.Wallet) error
	Update(wallet *dao.Wallet) error
	FindByUserID(userID string) (*dao.Wallet, error)
	FindByUUID(uuid string) (*dao.Wallet, error)
	GetOrCreate(userID string) (*dao.Wallet, error)

	// Transactions
	CreateTransaction(tx *dao.WalletTransaction) error
	FindTransactionsByUserID(userID string, limit, offset int) ([]dao.WalletTransaction, int64, error)
	FindTransactionsByWalletID(walletID uint, limit, offset int) ([]dao.WalletTransaction, int64, error)

	// With transaction support
	UpdateWithTx(tx *gorm.DB, wallet *dao.Wallet) error
	CreateTransactionWithTx(tx *gorm.DB, transaction *dao.WalletTransaction) error
}

type walletRepository struct {
	db *gorm.DB
}

// NewWalletRepository creates a new wallet repository
func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

// Create creates a new wallet
func (r *walletRepository) Create(wallet *dao.Wallet) error {
	return r.db.Create(wallet).Error
}

// Update updates an existing wallet
func (r *walletRepository) Update(wallet *dao.Wallet) error {
	return r.db.Save(wallet).Error
}

// FindByUserID finds a wallet by user ID
func (r *walletRepository) FindByUserID(userID string) (*dao.Wallet, error) {
	var wallet dao.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// FindByUUID finds a wallet by its UUID
func (r *walletRepository) FindByUUID(uuid string) (*dao.Wallet, error) {
	var wallet dao.Wallet
	err := r.db.Where("wallet_uuid = ?", uuid).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

// GetOrCreate gets an existing wallet or creates a new one
func (r *walletRepository) GetOrCreate(userID string) (*dao.Wallet, error) {
	wallet, err := r.FindByUserID(userID)
	if err == nil {
		return wallet, nil
	}

	if err == gorm.ErrRecordNotFound {
		// Create new wallet
		newWallet := &dao.Wallet{
			UserID:   userID,
			Currency: "ARS",
		}
		if err := r.Create(newWallet); err != nil {
			return nil, err
		}
		return newWallet, nil
	}

	return nil, err
}

// CreateTransaction creates a new wallet transaction
func (r *walletRepository) CreateTransaction(tx *dao.WalletTransaction) error {
	return r.db.Create(tx).Error
}

// FindTransactionsByUserID finds transactions for a user with pagination
func (r *walletRepository) FindTransactionsByUserID(userID string, limit, offset int) ([]dao.WalletTransaction, int64, error) {
	var transactions []dao.WalletTransaction
	var total int64

	query := r.db.Model(&dao.WalletTransaction{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&transactions).Error
	return transactions, total, err
}

// FindTransactionsByWalletID finds transactions for a wallet with pagination
func (r *walletRepository) FindTransactionsByWalletID(walletID uint, limit, offset int) ([]dao.WalletTransaction, int64, error) {
	var transactions []dao.WalletTransaction
	var total int64

	query := r.db.Model(&dao.WalletTransaction{}).Where("wallet_id = ?", walletID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&transactions).Error
	return transactions, total, err
}

// UpdateWithTx updates a wallet within a transaction
func (r *walletRepository) UpdateWithTx(tx *gorm.DB, wallet *dao.Wallet) error {
	return tx.Save(wallet).Error
}

// CreateTransactionWithTx creates a transaction within a database transaction
func (r *walletRepository) CreateTransactionWithTx(tx *gorm.DB, transaction *dao.WalletTransaction) error {
	return tx.Create(transaction).Error
}
