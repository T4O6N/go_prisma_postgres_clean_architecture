package dto

// NOTE - create user dto
type CreateUserDto struct {
	Name      string `json:"name"`
	Email     string `json:"email,omitempty"`
	SubjectID int    `json:"subject_id,omitempty"`
}

// NOTE - update user dto
type UpdateUserDto struct {
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	SubjectID int    `json:"subject_id,omitempty"`
	Status    bool   `json:"status,omitempty"`
}

// NOTE - user query dto
type UserQueryDto struct {
	Page  int    `form:"page" binding:"omitempty,min=1"`          // Default 1, must be >= 1
	Limit int    `form:"limit" binding:"omitempty,min=1,max=100"` // Default 10, must be between 1-100
	Name  string `form:"name" binding:"omitempty"`                // Optional name filter
}
