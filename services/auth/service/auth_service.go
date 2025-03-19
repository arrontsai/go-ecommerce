package service

import (
	"context"
	"errors"

	"github.com/arrontsai/ecommerce/pkg/middleware"
	"github.com/arrontsai/ecommerce/pkg/models"
	"github.com/arrontsai/ecommerce/services/auth/repository"
)

// AuthService defines the interface for authentication service operations
type AuthService interface {
	Register(ctx context.Context, req models.UserRegistration) (*models.UserResponse, error)
	Login(ctx context.Context, req models.UserLogin) (string, *models.UserResponse, error)
	GetUserByID(ctx context.Context, id string) (*models.UserResponse, error)
}

// DefaultAuthService implements AuthService
type DefaultAuthService struct {
	userRepo  repository.UserRepository
	jwtSecret string
	jwtExpiry int
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo repository.UserRepository, jwtSecret string, jwtExpiry int) AuthService {
	return &DefaultAuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

// Register registers a new user
func (s *DefaultAuthService) Register(ctx context.Context, req models.UserRegistration) (*models.UserResponse, error) {
	// Check if the user already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("用戶電子郵件已存在")
	}

	// Create the user
	user, err := models.NewUser(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		return nil, err
	}

	// Save the user
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Return the user response
	response := user.ToResponse()
	return &response, nil
}

// Login authenticates a user and returns a JWT token
func (s *DefaultAuthService) Login(ctx context.Context, req models.UserLogin) (string, *models.UserResponse, error) {
	// Find the user
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", nil, err
	}
	if user == nil {
		return "", nil, errors.New("用戶不存在")
	}

	// Check the password
	if !user.CheckPassword(req.Password) {
		return "", nil, errors.New("密碼不正確")
	}

	// Generate a JWT token
	token, err := middleware.GenerateJWT(user.ID, user.Role, s.jwtSecret, s.jwtExpiry/3600) // Convert seconds to hours
	if err != nil {
		return "", nil, err
	}

	// Return the token and user response
	response := user.ToResponse()
	return token, &response, nil
}

// GetUserByID gets a user by ID
func (s *DefaultAuthService) GetUserByID(ctx context.Context, id string) (*models.UserResponse, error) {
	// Find the user
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用戶不存在")
	}

	// Return the user response
	response := user.ToResponse()
	return &response, nil
}
