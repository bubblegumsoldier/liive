package handlers

import (
    "net/http"

    "github.com/labstack/echo/v4"
    "github.com/liive/backend/liive-auth/internal/service"
    "github.com/liive/backend/liive-auth/internal/types"
)

type AuthHandler struct {
    authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
    return &AuthHandler{
        authService: authService,
    }
}

// Login godoc
// @Summary Login user
// @Description Login with email and password to get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body types.LoginRequest true "Login credentials"
// @Success 200 {object} types.LoginResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /login [post]
func (h *AuthHandler) Login(c echo.Context) error {
    var req types.LoginRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request format"})
    }

    if err := c.Validate(&req); err != nil {
        return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
    }

    token, err := h.authService.Login(c.Request().Context(), req.Email, req.Password)
    if err != nil {
        switch err {
        case service.ErrInvalidCredentials:
            return c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "Invalid email or password"})
        default:
            return c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to login"})
        }
    }

    return c.JSON(http.StatusOK, types.LoginResponse{Token: token})
} 