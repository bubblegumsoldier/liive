package types

type LoginRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}

type LoginResponse struct {
    Token string `json:"token"`
} 