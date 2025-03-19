package dto

import "time"

type UserResponseDto struct {
	ID        int       `json:"id" example:"1"`
	Name      string    `json:"name" example:"ton"`
	Email     string    `json:"email" example:"email"`
	SubjectID int       `json:"subject_id" example:"1"`
	Status    bool      `json:"status" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2025-03-10T11:40:50.207+07:00"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-03-10T11:40:50.207+07:00"`
}

// MetaData represents pagination metadata
type MetaData struct {
	Page       int `json:"page" example:"1"`
	Limit      int `json:"limit" example:"10"`
	Total      int `json:"total" example:"3"`
	TotalPages int `json:"totalPages" example:"1"`
}

// UserListResponseDto represents the paginated response for users
type UserListResponseDto struct {
	Data []UserResponseDto `json:"data"`
	Meta MetaData          `json:"meta"`
}
