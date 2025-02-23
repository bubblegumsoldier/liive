package service

import (
	"context"
	"errors"

	"github.com/liive/backend/liive-user-manager/internal/types"
	"github.com/liive/backend/shared/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrEmailExists     = errors.New("email already exists")
	ErrUsernameExists  = errors.New("username already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) Register(ctx context.Context, req types.RegisterRequest) (*models.User, error) {
	// Check if email exists
	var count int64
	if err := s.db.Model(&models.User{}).Where("email = ?", req.Email).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrEmailExists
	}

	// Check if username exists
	if err := s.db.Model(&models.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, ErrUsernameExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    req.Email,
		Username: req.Username,
		Password: string(hashedPassword),
		IsActive: true,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID uint, req types.UpdateProfileRequest) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	if req.Username != "" && req.Username != user.Username {
		var count int64
		if err := s.db.Model(&models.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, ErrUsernameExists
		}
		user.Username = req.Username
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) UpdatePassword(ctx context.Context, userID uint, req types.UpdatePasswordRequest) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return ErrInvalidPassword
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
} 