package dto

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	Status      int    `json:"status"`
	Message     string `json:"message"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
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
	Status  int    `json:"status"`
	Message string `json:"message"`
}