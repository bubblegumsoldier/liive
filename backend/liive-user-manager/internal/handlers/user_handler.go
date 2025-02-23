package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/liive/backend/liive-user-manager/internal/types"
	"github.com/liive/backend/liive-user-manager/internal/service"
	"github.com/liive/backend/shared/pkg/models"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, username and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body types.RegisterRequest true "User registration details"
// @Success 201 {object} types.UserResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 409 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /register [post]
func (h *UserHandler) Register(c echo.Context) error {
	var req types.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request format"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
	}

	user, err := h.userService.Register(c.Request().Context(), req)
	if err != nil {
		switch err {
		case service.ErrEmailExists:
			return c.JSON(http.StatusConflict, types.ErrorResponse{Error: "Email already exists"})
		case service.ErrUsernameExists:
			return c.JSON(http.StatusConflict, types.ErrorResponse{Error: "Username already exists"})
		default:
			return c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to register user"})
		}
	}

	return c.JSON(http.StatusCreated, toUserResponse(user))
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update user profile information
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body types.UpdateProfileRequest true "User profile update details"
// @Success 200 {object} types.UserResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Failure 409 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /api/profile [put]
func (h *UserHandler) UpdateProfile(c echo.Context) error {
	userID := c.Get("user_id").(uint) // Set by auth middleware
	var req types.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request format"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
	}

	user, err := h.userService.UpdateProfile(c.Request().Context(), userID, req)
	if err != nil {
		switch err {
		case service.ErrUsernameExists:
			return c.JSON(http.StatusConflict, types.ErrorResponse{Error: "Username already exists"})
		case service.ErrUserNotFound:
			return c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "User not found"})
		default:
			return c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to update profile"})
		}
	}

	return c.JSON(http.StatusOK, toUserResponse(user))
}

// UpdatePassword godoc
// @Summary Update user password
// @Description Update user password with current and new password
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param user body types.UpdatePasswordRequest true "Password update details"
// @Success 200
// @Failure 400 {object} types.ErrorResponse
// @Failure 401 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /api/password [put]
func (h *UserHandler) UpdatePassword(c echo.Context) error {
	userID := c.Get("user_id").(uint) // Set by auth middleware
	var req types.UpdatePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid request format"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
	}

	if err := h.userService.UpdatePassword(c.Request().Context(), userID, req); err != nil {
		switch err {
		case service.ErrInvalidPassword:
			return c.JSON(http.StatusUnauthorized, types.ErrorResponse{Error: "Current password is incorrect"})
		case service.ErrUserNotFound:
			return c.JSON(http.StatusNotFound, types.ErrorResponse{Error: "User not found"})
		default:
			return c.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to update password"})
		}
	}

	return c.NoContent(http.StatusOK)
}

func toUserResponse(user *models.User) types.UserResponse {
	return types.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
	}
} 