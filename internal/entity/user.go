package entity

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	SubjectID int       `json:"subject_id,omitempty"`
	Status    bool      `json:"status"`
	Day       int       `json:"day"`
	Month     int       `json:"month"`
	Year      int       `json:"year"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	SubjectID int    `json:"subject_id,omitempty"`
}

type UpdateUserRequest struct {
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	SubjectID int    `json:"subject_id,omitempty"`
	Status    bool   `json:"status,omitempty"`
}

type UserResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	SubjectID int       `json:"subject_id,omitempty"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserListResponse struct {
	Data []UserResponse `json:"data"`
	Meta struct {
		Limit      int `json:"limit"`
		Page       int `json:"page"`
		Total      int `json:"total"`
		TotalPages int `json:"totalPages"`
	} `json:"meta"`
}

type ErrorResponse struct {
	// StatusCode int    `json:"status_code"`
	Message string `json:"message"`
}
