package dto

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int          `json:"expires_in"`
	User        UserResponse `json:"user"`
}

type UserResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Email       *string `json:"email"`
	PhoneNumber string  `json:"phone_number"`
	RoleID      string  `json:"role_id"`
	GroupID     string  `json:"group_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}