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
	GetAllUsers(page, limit int, roleFilter, search string) ([]*domain.UserDTO, int64, error)
	GetUserByID(id int64) (*domain.UserDTO, error)
	GetUserProfile(id int64) (*domain.UserDTO, error)
	UpdateUser(id int64, req domain.UpdateUserRequest) (*domain.UserDTO, error)
	DeleteUser(id int64) error
	ForceReauthentication(id int64) error
}

type userService struct {
	userRepo     repository.UserRepository
	emailService EmailService
}

// NewUserService crea una nueva instancia del servicio de usuarios
func NewUserService(userRepo repository.UserRepository, emailService EmailService) UserService {
	return &userService{
		userRepo:     userRepo,
		emailService: emailService,
	}
}

// GetAllUsers obtiene todos los usuarios con paginación y filtros (solo admin)
func (s *userService) GetAllUsers(page, limit int, roleFilter, search string) ([]*domain.UserDTO, int64, error) {
	users, total, err := s.userRepo.FindAllWithPagination(page, limit, roleFilter, search)
	if err != nil {
		return nil, 0, err
	}

	// Convertir a DTOs
	userDTOs := make([]*domain.UserDTO, len(users))
	for i, user := range users {
		userDTOs[i] = s.convertToDTO(user)
	}

	return userDTOs, total, nil
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

// ForceReauthentication desverifica el email y reenvía el email de verificación
func (s *userService) ForceReauthentication(id int64) error {
	// Verificar que el usuario existe
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("usuario no encontrado")
		}
		return err
	}

	// Generar nuevo token de verificación
	token, err := s.emailService.GenerateToken()
	if err != nil {
		return err
	}

	// Guardar nuevo token y desverificar email
	if err := s.userRepo.SaveEmailVerificationToken(id, token); err != nil {
		return err
	}

	if err := s.userRepo.UnverifyEmail(id, user.Email); err != nil {
		return err
	}

	// Enviar email de forma asíncrona
	go func() {
		if err := s.emailService.SendVerificationEmail(user.Email, token); err != nil {
			// El error ya está logueado en emailService
			return
		}
	}()

	return nil
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
