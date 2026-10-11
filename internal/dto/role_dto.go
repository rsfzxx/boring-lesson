package dto

type RoleResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RoleListResponse struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Data    []RoleResponse `json:"data"`
}

type RoleDetailResponse struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    RoleResponse `json:"data"`
}