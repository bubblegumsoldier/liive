package types

type UserResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	IsActive  bool   `json:"isActive"`
}

type ErrorResponse struct {
	Error string `json:"error"`
} 