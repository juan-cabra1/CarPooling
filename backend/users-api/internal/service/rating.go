package service

import (
	"errors"
	"users-api/internal/dao"
	"users-api/internal/domain"
	"users-api/internal/repository"
)

// RatingService define las operaciones de calificaciones
type RatingService interface {
	CreateRating(req domain.CreateRatingRequest) error
	GetUserRatings(userID int64, page, limit int) ([]domain.RatingDTO, int, error)
}

type ratingService struct {
	ratingRepo repository.RatingRepository
	userRepo   repository.UserRepository
}

// NewRatingService crea una nueva instancia del servicio de calificaciones
func NewRatingService(ratingRepo repository.RatingRepository, userRepo repository.UserRepository) RatingService {
	return &ratingService{
		ratingRepo: ratingRepo,
		userRepo:   userRepo,
	}
}

// CreateRating crea una nueva calificación y actualiza los promedios del usuario
func (s *ratingService) CreateRating(req domain.CreateRatingRequest) error {
	// 1. Validar que no existe un rating duplicado
	exists, err := s.ratingRepo.ExistsRating(req.RaterID, req.TripID, req.RatedUserID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("ya existe una calificación para este viaje")
	}

	// 2. Crear el rating en la base de datos
	ratingDAO := &dao.RatingDAO{
		RaterID:     req.RaterID,
		RatedUserID: req.RatedUserID,
		TripID:      req.TripID,
		RoleRated:   req.RoleRated,
		Score:       req.Score,
		Comment:     req.Comment,
	}

	if err := s.ratingRepo.Create(ratingDAO); err != nil {
		return err
	}

	// 3. Calcular nuevos promedios del usuario calificado
	avgDriver, avgPassenger, totalDriver, totalPassenger, err := s.ratingRepo.CalculateAverages(req.RatedUserID)
	if err != nil {
		return err
	}

	// 4. Actualizar los promedios del usuario
	if err := s.userRepo.UpdateRatings(req.RatedUserID, avgDriver, avgPassenger, totalDriver, totalPassenger); err != nil {
		return err
	}

	return nil
}

// GetUserRatings obtiene las calificaciones de un usuario con paginación
func (s *ratingService) GetUserRatings(userID int64, page, limit int) ([]domain.RatingDTO, int, error) {
	// Validar parámetros de paginación
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Calcular offset
	offset := (page - 1) * limit

	// Obtener ratings
	ratings, err := s.ratingRepo.FindByRatedUserID(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Convertir a DTOs
	ratingDTOs := make([]domain.RatingDTO, len(ratings))
	for i, rating := range ratings {
		ratingDTOs[i] = s.convertToDTO(&rating)
	}

	// Para el total, obtener count total de ratings del usuario
	// (simplificado: usar la longitud del resultado por ahora)
	total := len(ratingDTOs)

	return ratingDTOs, total, nil
}

// convertToDTO convierte un RatingDAO a RatingDTO
func (s *ratingService) convertToDTO(ratingDAO *dao.RatingDAO) domain.RatingDTO {
	return domain.RatingDTO{
		ID:          ratingDAO.ID,
		RaterID:     ratingDAO.RaterID,
		RatedUserID: ratingDAO.RatedUserID,
		TripID:      ratingDAO.TripID,
		RoleRated:   ratingDAO.RoleRated,
		Score:       ratingDAO.Score,
		Comment:     ratingDAO.Comment,
		CreatedAt:   ratingDAO.CreatedAt,
	}
}
