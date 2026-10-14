package dto

import "time"

type CreateUserRequest struct {
	Name        string  `json:"name"`
	Email       *string `json:"email"`
	PhoneNumber string  `json:"phone_number"`
	Password    string  `json:"password"`
	RoleID      string  `json:"role_id"`
	GroupID     string  `json:"group_id"`
}

type UpdateUserRequest struct {
	Name        *string `json:"name"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	Password    *string `json:"password"`
	RoleID      *string `json:"role_id"`
	GroupID     *string `json:"group_id"`
}

type UserData struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       *string   `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	RoleID      string    `json:"role_id"`
	GroupID     string    `json:"group_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserListResponse struct {
	Status  int        `json:"status"`
	Message string     `json:"message"`
	Data    []UserData `json:"data"`
}

type UserDetailResponse struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Data    UserData `json:"data"`
}

type MessageResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}