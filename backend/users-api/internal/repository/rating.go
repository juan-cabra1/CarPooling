package repository

import (
	"users-api/internal/dao"

	"gorm.io/gorm"
)

// RatingRepository define las operaciones de acceso a datos para calificaciones
type RatingRepository interface {
	Create(rating *dao.RatingDAO) error
	FindByRatedUserID(userID int64, limit, offset int) ([]dao.RatingDAO, error)
	CalculateAverages(userID int64) (avgDriver, avgPassenger float64, totalDriver, totalPassenger int, err error)
	ExistsRating(raterID int64, tripID string, ratedUserID int64) (bool, error)
}

type ratingRepository struct {
	db *gorm.DB
}

// NewRatingRepository crea una nueva instancia del repositorio de calificaciones
func NewRatingRepository(db *gorm.DB) RatingRepository {
	return &ratingRepository{db: db}
}

func (r *ratingRepository) Create(rating *dao.RatingDAO) error {
	return r.db.Create(rating).Error
}

func (r *ratingRepository) FindByRatedUserID(userID int64, limit, offset int) ([]dao.RatingDAO, error) {
	var ratings []dao.RatingDAO
	err := r.db.Where("rated_user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&ratings).Error
	return ratings, err
}

// CalculateAverages calcula los promedios de calificaciones por rol
func (r *ratingRepository) CalculateAverages(userID int64) (avgDriver, avgPassenger float64, totalDriver, totalPassenger int, err error) {
	// Calcular promedio y total para conductor
	type Result struct {
		Avg   float64
		Count int
	}

	var driverResult Result
	err = r.db.Model(&dao.RatingDAO{}).
		Select("COALESCE(AVG(score), 0) as avg, COUNT(*) as count").
		Where("rated_user_id = ? AND role_rated = ?", userID, "conductor").
		Scan(&driverResult).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}

	var passengerResult Result
	err = r.db.Model(&dao.RatingDAO{}).
		Select("COALESCE(AVG(score), 0) as avg, COUNT(*) as count").
		Where("rated_user_id = ? AND role_rated = ?", userID, "pasajero").
		Scan(&passengerResult).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return driverResult.Avg, passengerResult.Avg, driverResult.Count, passengerResult.Count, nil
}

func (r *ratingRepository) ExistsRating(raterID int64, tripID string, ratedUserID int64) (bool, error) {
	var count int64
	err := r.db.Model(&dao.RatingDAO{}).
		Where("rater_id = ? AND trip_id = ? AND rated_user_id = ?", raterID, tripID, ratedUserID).
		Count(&count).Error
	return count > 0, err
}
