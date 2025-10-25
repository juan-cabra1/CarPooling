package service

import (
	"errors"
	"users-api/internal/dao"
	"users-api/internal/domain"
	"users-api/internal/repository"

	"gorm.io/gorm"
)

// UserService define las operaciones de gestión de usuarios
type UserService interface {
	GetUserByID(id int64) (*domain.UserDTO, error)
	GetUserProfile(id int64) (*domain.UserDTO, error)
	UpdateUser(id int64, req domain.UpdateUserRequest) (*domain.UserDTO, error)
	DeleteUser(id int64) error
}

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService crea una nueva instancia del servicio de usuarios
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// GetUserByID obtiene un usuario por su ID
func (s *userService) GetUserByID(id int64) (*domain.UserDTO, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	return s.convertToDTO(user), nil
}

// GetUserProfile es un alias de GetUserByID usado para /users/me
func (s *userService) GetUserProfile(id int64) (*domain.UserDTO, error) {
	return s.GetUserByID(id)
}

// UpdateUser actualiza los datos de un usuario
func (s *userService) UpdateUser(id int64, req domain.UpdateUserRequest) (*domain.UserDTO, error) {
	// Buscar usuario existente
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	// Actualizar solo los campos que no son nulos
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Lastname != nil {
		user.Lastname = *req.Lastname
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.Street != nil {
		user.Street = *req.Street
	}
	if req.Number != nil {
		user.Number = *req.Number
	}
	if req.PhotoURL != nil {
		user.PhotoURL = *req.PhotoURL
	}

	// Guardar cambios
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return s.convertToDTO(user), nil
}

// DeleteUser elimina un usuario
func (s *userService) DeleteUser(id int64) error {
	// Verificar que el usuario existe
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("usuario no encontrado")
		}
		return err
	}

	return s.userRepo.Delete(id)
}

// convertToDTO convierte un UserDAO a UserDTO
func (s *userService) convertToDTO(userDAO *dao.UserDAO) *domain.UserDTO {
	return &domain.UserDTO{
		ID:                  userDAO.ID,
		Email:               userDAO.Email,
		EmailVerified:       userDAO.EmailVerified,
		Name:                userDAO.Name,
		Lastname:            userDAO.Lastname,
		Role:                userDAO.Role,
		Phone:               userDAO.Phone,
		Street:              userDAO.Street,
		Number:              userDAO.Number,
		PhotoURL:            userDAO.PhotoURL,
		Sex:                 userDAO.Sex,
		AvgDriverRating:     userDAO.AvgDriverRating,
		AvgPassengerRating:  userDAO.AvgPassengerRating,
		TotalTripsPassenger: userDAO.TotalTripsPassenger,
		TotalTripsDriver:    userDAO.TotalTripsDriver,
		Birthdate:           userDAO.Birthdate,
		CreatedAt:           userDAO.CreatedAt,
		UpdatedAt:           userDAO.UpdatedAt,
	}
}
