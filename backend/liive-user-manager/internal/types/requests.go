package types

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=30"`
	Password  string `json:"password" validate:"required,min=8"`
}

type UpdateProfileRequest struct {
	Username  string `json:"username" validate:"omitempty,min=3,max=30"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=8"`
} 