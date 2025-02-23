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

type Claims struct {
    UserID    uint   `json:"user_id"`
    Email     string `json:"email"`
    Roles     []string `json:"roles"`
    jwt.RegisteredClaims
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
