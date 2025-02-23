package service

import (
    "context"
    "errors"

    "github.com/liive/backend/shared/pkg/auth"
    "github.com/liive/backend/shared/pkg/models"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

var (
    ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
    db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
    return &AuthService{
        db: db,
    }
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
    var user models.User
    if err := s.db.Preload("Roles").Where("email = ?", email).First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return "", ErrInvalidCredentials
        }
        return "", err
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return "", ErrInvalidCredentials
    }

    token, err := auth.GenerateToken(&user)
    if err != nil {
        return "", err
    }

    // Update last login time
    user.LastLogin = s.db.NowFunc()
    if err := s.db.Save(&user).Error; err != nil {
        // Log error but don't fail the login
        // TODO: Add proper logging
    }

    return token, nil
} 