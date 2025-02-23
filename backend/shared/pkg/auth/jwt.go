package auth

import (
    "fmt"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/liive/backend/shared/pkg/models"
)

var jwtKey = []byte(getSecretKey())

func getSecretKey() string {
    key := os.Getenv("JWT_SECRET_KEY")
    if key == "" {
        key = "your-256-bit-secret" // Default key for development
    }
    return key
}

// JWTToken is a wrapper around jwt.Token to avoid type conflicts
type JWTToken = *jwt.Token

type Claims struct {
    UserID    uint   `json:"user_id"`
    Email     string `json:"email"`
    Roles     []string `json:"roles"`
    jwt.RegisteredClaims
}

// Valid implements the jwt.Claims interface
func (c *Claims) Valid() error {
    if c.UserID == 0 {
        return fmt.Errorf("missing user ID")
    }
    if c.Email == "" {
        return fmt.Errorf("missing email")
    }
    now := time.Now()
    if !c.ExpiresAt.IsZero() && now.After(c.ExpiresAt.Time) {
        return fmt.Errorf("token has expired")
    }
    if !c.NotBefore.IsZero() && now.Before(c.NotBefore.Time) {
        return fmt.Errorf("token is not valid yet")
    }
    return nil
}

func GenerateToken(user *models.User) (string, error) {
    // Create roles slice
    roles := make([]string, len(user.Roles))
    for i, role := range user.Roles {
        roles[i] = role.Name
    }

    // Create claims
    claims := &Claims{
        UserID: user.ID,
        Email:  user.Email,
        Roles:  roles,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    // Generate token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func ValidateToken(tokenStr string) (*Claims, error) {
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return jwtKey, nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    return claims, nil
}
