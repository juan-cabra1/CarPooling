package service

import (
	"errors"
	"time"
	"users-api/internal/dao"
	"users-api/internal/domain"
	"users-api/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService define las operaciones de autenticación y gestión de usuarios
type AuthService interface {
	// Login y Registro
	Register(req domain.CreateUserRequest) (*domain.UserDTO, error)
	Login(req domain.LoginRequest) (*domain.LoginResponse, error)

	// JWT
	ValidateToken(tokenString string) (*jwt.Token, error)
	GenerateJWT(userID int64, email, role string) (string, error)

	// Verificación de email
	VerifyEmail(token string) error
	ResendVerificationEmail(email string) error

	// Gestión de contraseña
	RequestPasswordReset(email string) error
	ResetPassword(token, newPassword string) error
	ChangePassword(userID int64, currentPassword, newPassword string) error
}

type authService struct {
	userRepo     repository.UserRepository
	emailService EmailService
	jwtSecret    string
}

// NewAuthService crea una nueva instancia del servicio de autenticación
func NewAuthService(userRepo repository.UserRepository, emailService EmailService, jwtSecret string) AuthService {
	return &authService{
		userRepo:     userRepo,
		emailService: emailService,
		jwtSecret:    jwtSecret,
	}
}

// ==================== LOGIN Y REGISTRO ====================

// Register crea un nuevo usuario y envía email de verificación
func (s *authService) Register(req domain.CreateUserRequest) (*domain.UserDTO, error) {
	// Verificar si el email ya existe
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("el email ya está registrado")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Hashear la contraseña con bcrypt cost 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, err
	}

	// Parsear la fecha de nacimiento
	birthdate, err := time.Parse("2006-01-02", req.Birthdate)
	if err != nil {
		return nil, errors.New("formato de fecha inválido, usar YYYY-MM-DD")
	}

	// Generar token de verificación
	verificationToken, err := s.emailService.GenerateToken()
	if err != nil {
		return nil, err
	}

	// Crear el usuario
	userDAO := &dao.UserDAO{
		Email:                  req.Email,
		EmailVerified:          false,
		EmailVerificationToken: &verificationToken,
		Name:                   req.Name,
		Lastname:               req.Lastname,
		PasswordHash:           string(hashedPassword),
		Role:                   "user",
		Phone:                  req.Phone,
		Street:                 req.Street,
		Number:                 req.Number,
		PhotoURL:               req.PhotoURL,
		Sex:                    req.Sex,
		Birthdate:              birthdate,
	}

	if err := s.userRepo.Create(userDAO); err != nil {
		return nil, err
	}

	// Enviar email de verificación de forma asíncrona
	go s.emailService.SendVerificationEmail(req.Email, verificationToken)

	// Convertir a DTO y retornar
	return s.convertToDTO(userDAO), nil
}

// Login autentica a un usuario y retorna un JWT
func (s *authService) Login(req domain.LoginRequest) (*domain.LoginResponse, error) {
	// Buscar usuario por email
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("credenciales inválidas")
		}
		return nil, err
	}

	// Verificar la contraseña
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Generar JWT
	token, err := s.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		Token: token,
		User:  s.convertToDTO(user),
	}, nil
}

// ==================== JWT ====================

// GenerateJWT genera un token JWT con 24 horas de expiración
func (s *authService) GenerateJWT(userID int64, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // 24 horas
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ValidateToken valida un token JWT
func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenExpired
	}

	return token, nil
}

// ==================== VERIFICACIÓN DE EMAIL ====================

// VerifyEmail verifica el email de un usuario usando el token
func (s *authService) VerifyEmail(token string) error {
	// Buscar usuario por token de verificación
	user, err := s.userRepo.FindByEmailVerificationToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("token de verificación inválido")
		}
		return err
	}

	// Marcar email como verificado
	if err := s.userRepo.UpdateEmailVerified(user.ID, true); err != nil {
		return err
	}

	// Limpiar el token de verificación
	return s.userRepo.SaveEmailVerificationToken(user.ID, "")
}

// ResendVerificationEmail reenvía el email de verificación
func (s *authService) ResendVerificationEmail(email string) error {
	// Buscar usuario por email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("usuario no encontrado")
		}
		return err
	}

	// Verificar si ya está verificado
	if user.EmailVerified {
		return errors.New("el email ya está verificado")
	}

	// Generar nuevo token
	token, err := s.emailService.GenerateToken()
	if err != nil {
		return err
	}

	// Guardar nuevo token
	if err := s.userRepo.SaveEmailVerificationToken(user.ID, token); err != nil {
		return err
	}

	// Enviar email
	go s.emailService.SendVerificationEmail(user.Email, token)

	return nil
}

// ==================== GESTIÓN DE CONTRASEÑA ====================

// RequestPasswordReset envía un email con el token para restablecer contraseña
func (s *authService) RequestPasswordReset(email string) error {
	// Buscar usuario por email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// No revelar si el email existe o no por seguridad
		// Retornar nil para evitar enumeration attacks
		return nil
	}

	// Generar token de reset
	token, err := s.emailService.GenerateToken()
	if err != nil {
		return err
	}

	// Guardar token con expiración de 1 hora
	expiresAt := time.Now().Add(1 * time.Hour)
	if err := s.userRepo.SavePasswordResetToken(user.ID, token, expiresAt); err != nil {
		return err
	}

	// Enviar email de forma asíncrona
	go s.emailService.SendPasswordResetEmail(user.Email, token)

	return nil
}

// ResetPassword cambia la contraseña usando el token de reset
func (s *authService) ResetPassword(token, newPassword string) error {
	// Buscar usuario por token de reset (valida que no esté expirado)
	user, err := s.userRepo.FindByPasswordResetToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("token de reset inválido o expirado")
		}
		return err
	}

	// Validar longitud mínima de contraseña
	if len(newPassword) < 8 {
		return errors.New("la contraseña debe tener al menos 8 caracteres")
	}

	// Hashear nueva contraseña con bcrypt cost 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return err
	}

	// Actualizar contraseña
	if err := s.userRepo.UpdatePassword(user.ID, string(hashedPassword)); err != nil {
		return err
	}

	// Limpiar token de reset
	return s.userRepo.ClearPasswordResetToken(user.ID)
}

// ChangePassword permite al usuario cambiar su contraseña estando autenticado
func (s *authService) ChangePassword(userID int64, currentPassword, newPassword string) error {
	// Buscar usuario
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("usuario no encontrado")
	}

	// Verificar contraseña actual
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword))
	if err != nil {
		return errors.New("contraseña actual incorrecta")
	}

	// Validar longitud mínima de nueva contraseña
	if len(newPassword) < 8 {
		return errors.New("la nueva contraseña debe tener al menos 8 caracteres")
	}

	// Hashear nueva contraseña con bcrypt cost 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return err
	}

	// Actualizar contraseña
	return s.userRepo.UpdatePassword(user.ID, string(hashedPassword))
}

// ==================== HELPERS ====================

// convertToDTO convierte un UserDAO a UserDTO
func (s *authService) convertToDTO(userDAO *dao.UserDAO) *domain.UserDTO {
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
