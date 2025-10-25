package service

import (
	"errors"
	"testing"
	"time"
	"users-api/internal/dao"
	"users-api/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockUserRepository es un mock del repositorio de usuarios
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *dao.UserDAO) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id int64) (*dao.UserDAO, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dao.UserDAO), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*dao.UserDAO, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dao.UserDAO), args.Error(1)
}

func (m *MockUserRepository) Update(user *dao.UserDAO) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateRatings(userID int64, avgDriverRating, avgPassengerRating float64, totalTripsDriver, totalTripsPassenger int) error {
	args := m.Called(userID, avgDriverRating, avgPassengerRating, totalTripsDriver, totalTripsPassenger)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) UpdatePassword(userID int64, newPasswordHash string) error {
	args := m.Called(userID, newPasswordHash)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateEmailVerified(userID int64, verified bool) error {
	args := m.Called(userID, verified)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmailVerificationToken(token string) (*dao.UserDAO, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dao.UserDAO), args.Error(1)
}

func (m *MockUserRepository) FindByPasswordResetToken(token string) (*dao.UserDAO, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dao.UserDAO), args.Error(1)
}

func (m *MockUserRepository) SaveEmailVerificationToken(userID int64, token string) error {
	args := m.Called(userID, token)
	return args.Error(0)
}

func (m *MockUserRepository) SavePasswordResetToken(userID int64, token string, expiresAt time.Time) error {
	args := m.Called(userID, token, expiresAt)
	return args.Error(0)
}

func (m *MockUserRepository) ClearPasswordResetToken(userID int64) error {
	args := m.Called(userID)
	return args.Error(0)
}

// Test 1: TestGetUserByID_Success
func TestGetUserByID_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	expectedUser := &dao.UserDAO{
		ID:                  1,
		Email:               "test@example.com",
		EmailVerified:       true,
		Name:                "Juan",
		Lastname:            "Pérez",
		PasswordHash:        "hashedpassword",
		Role:                "user",
		Phone:               "123456789",
		Street:              "Calle Falsa",
		Number:              123,
		PhotoURL:            "http://photo.url",
		Sex:                 "hombre",
		AvgDriverRating:     4.5,
		AvgPassengerRating:  4.8,
		TotalTripsPassenger: 10,
		TotalTripsDriver:    5,
		Birthdate:           time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	mockRepo.On("FindByID", int64(1)).Return(expectedUser, nil)

	// Execute
	result, err := service.GetUserByID(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "Juan", result.Name)
	assert.Equal(t, "Pérez", result.Lastname)
	assert.Equal(t, 4.5, result.AvgDriverRating)
	assert.Equal(t, 4.8, result.AvgPassengerRating)

	mockRepo.AssertExpectations(t)
}

// Test 2: TestGetUserByID_NotFound
func TestGetUserByID_NotFound(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	mockRepo.On("FindByID", int64(999)).Return(nil, gorm.ErrRecordNotFound)

	// Execute
	result, err := service.GetUserByID(999)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "usuario no encontrado", err.Error())

	mockRepo.AssertExpectations(t)
}

// Test 3: TestUpdateUser_Success
func TestUpdateUser_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	existingUser := &dao.UserDAO{
		ID:       1,
		Email:    "test@example.com",
		Name:     "Juan",
		Lastname: "Pérez",
		Phone:    "123456789",
		Street:   "Calle Falsa",
		Number:   123,
	}

	newName := "Carlos"
	newPhone := "987654321"
	updateReq := domain.UpdateUserRequest{
		Name:  &newName,
		Phone: &newPhone,
	}

	mockRepo.On("FindByID", int64(1)).Return(existingUser, nil)
	mockRepo.On("Update", mock.AnythingOfType("*dao.UserDAO")).Return(nil)

	// Execute
	result, err := service.UpdateUser(1, updateReq)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Carlos", result.Name)
	assert.Equal(t, "987654321", result.Phone)
	assert.Equal(t, "Pérez", result.Lastname) // No debe cambiar

	mockRepo.AssertExpectations(t)
}

// Test 4: TestUpdateUser_PartialUpdate
func TestUpdateUser_PartialUpdate(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	existingUser := &dao.UserDAO{
		ID:       1,
		Email:    "test@example.com",
		Name:     "Juan",
		Lastname: "Pérez",
		Phone:    "123456789",
		Street:   "Calle Falsa",
		Number:   123,
		PhotoURL: "old-photo.jpg",
	}

	// Solo actualizar PhotoURL, dejar todo lo demás igual
	newPhotoURL := "new-photo.jpg"
	updateReq := domain.UpdateUserRequest{
		PhotoURL: &newPhotoURL,
	}

	mockRepo.On("FindByID", int64(1)).Return(existingUser, nil)
	mockRepo.On("Update", mock.AnythingOfType("*dao.UserDAO")).Return(nil)

	// Execute
	result, err := service.UpdateUser(1, updateReq)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "new-photo.jpg", result.PhotoURL)
	// Verificar que los demás campos no cambiaron
	assert.Equal(t, "Juan", result.Name)
	assert.Equal(t, "Pérez", result.Lastname)
	assert.Equal(t, "123456789", result.Phone)
	assert.Equal(t, "Calle Falsa", result.Street)
	assert.Equal(t, 123, result.Number)

	mockRepo.AssertExpectations(t)
}

// Test adicional: TestUpdateUser_UserNotFound
func TestUpdateUser_UserNotFound(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	newName := "Carlos"
	updateReq := domain.UpdateUserRequest{
		Name: &newName,
	}

	mockRepo.On("FindByID", int64(999)).Return(nil, gorm.ErrRecordNotFound)

	// Execute
	result, err := service.UpdateUser(999, updateReq)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "usuario no encontrado", err.Error())

	mockRepo.AssertExpectations(t)
}

// Test adicional: TestUpdateUser_RepositoryError
func TestUpdateUser_RepositoryError(t *testing.T) {
	// Setup
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	existingUser := &dao.UserDAO{
		ID:   1,
		Name: "Juan",
	}

	newName := "Carlos"
	updateReq := domain.UpdateUserRequest{
		Name: &newName,
	}

	mockRepo.On("FindByID", int64(1)).Return(existingUser, nil)
	mockRepo.On("Update", mock.AnythingOfType("*dao.UserDAO")).Return(errors.New("database error"))

	// Execute
	result, err := service.UpdateUser(1, updateReq)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "database error", err.Error())

	mockRepo.AssertExpectations(t)
}
