package repository

import (
	"time"
	"users-api/internal/dao"

	"gorm.io/gorm"
)

// UserRepository define las operaciones de acceso a datos para usuarios
type UserRepository interface {
	Create(user *dao.UserDAO) error
	FindAllWithPagination(page, limit int, roleFilter, search string) ([]*dao.UserDAO, int64, error)
	FindByID(id int64) (*dao.UserDAO, error)
	FindByEmail(email string) (*dao.UserDAO, error)
	Update(user *dao.UserDAO) error
	UpdateRatings(userID int64, avgDriverRating, avgPassengerRating float64, totalTripsDriver, totalTripsPassenger int) error
	Delete(id int64) error
	UpdatePassword(userID int64, newPasswordHash string) error
	UpdateEmailVerified(userID int64, verified bool) error
	FindByEmailVerificationToken(token string) (*dao.UserDAO, error)
	FindByPasswordResetToken(token string) (*dao.UserDAO, error)
	SaveEmailVerificationToken(userID int64, token string) error
	SavePasswordResetToken(userID int64, token string, expiresAt time.Time) error
	ClearPasswordResetToken(userID int64) error
	UnverifyEmail(userID int64, email string) error
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository crea una nueva instancia del repositorio de usuarios
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *dao.UserDAO) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindAllWithPagination(page, limit int, roleFilter, search string) ([]*dao.UserDAO, int64, error) {
	var users []*dao.UserDAO
	var total int64

	query := r.db.Model(&dao.UserDAO{})

	// Filtro por rol
	if roleFilter != "" {
		query = query.Where("role = ?", roleFilter)
	}

	// Búsqueda por email o nombre
	if search != "" {
		query = query.Where("email LIKE ? OR name LIKE ? OR lastname LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Contar total antes de paginar
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Aplicar paginación
	offset := (page - 1) * limit
	err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) FindByID(id int64) (*dao.UserDAO, error) {
	var user dao.UserDAO
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*dao.UserDAO, error) {
	var user dao.UserDAO
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *dao.UserDAO) error {
	return r.db.Save(user).Error
}

func (r *userRepository) UpdateRatings(userID int64, avgDriverRating, avgPassengerRating float64, totalTripsDriver, totalTripsPassenger int) error {
	return r.db.Model(&dao.UserDAO{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"avg_driver_rating":     avgDriverRating,
			"avg_passenger_rating":  avgPassengerRating,
			"total_trips_driver":    totalTripsDriver,
			"total_trips_passenger": totalTripsPassenger,
		}).Error
}

func (r *userRepository) Delete(id int64) error {
	return r.db.Delete(&dao.UserDAO{}, id).Error
}

func (r *userRepository) UpdatePassword(userID int64, newPasswordHash string) error {
	return r.db.Model(&dao.UserDAO{}).
		Where("id = ?", userID).
		Update("password_hash", newPasswordHash).Error
}

func (r *userRepository) UpdateEmailVerified(userID int64, verified bool) error {
	return r.db.Model(&dao.UserDAO{}).
		Where("id = ?", userID).
		Update("email_verified", verified).Error
}

func (r *userRepository) FindByEmailVerificationToken(token string) (*dao.UserDAO, error) {
	var user dao.UserDAO
	err := r.db.Where("email_verification_token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByPasswordResetToken(token string) (*dao.UserDAO, error) {
	var user dao.UserDAO
	err := r.db.Where("password_reset_token = ? AND password_reset_expires > ?", token, time.Now()).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) SaveEmailVerificationToken(userID int64, token string) error {
	return r.db.Model(&dao.UserDAO{}).
		Where("id = ?", userID).
		Update("email_verification_token", token).Error
}

func (r *userRepository) SavePasswordResetToken(userID int64, token string, expiresAt time.Time) error {
	return r.db.Model(&dao.UserDAO{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"password_reset_token":   token,
			"password_reset_expires": expiresAt,
		}).Error
}

func (r *userRepository) ClearPasswordResetToken(userID int64) error {
	return r.db.Model(&dao.UserDAO{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"password_reset_token":   nil,
			"password_reset_expires": nil,
		}).Error
}
func (r *userRepository) UnverifyEmail(userID int64, email string) error {
	return r.db.Model(&dao.UserDAO{}).
		Where("id = ?", userID).
		Update("email_verified", false).Error
}
