package dto

type CreateUserDto struct {
	Name      string `json:"name"`
	Email     string `json:"email,omitempty"`
	SubjectID int    `json:"subject_id,omitempty"`
}

type UpdateUserDto struct {
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	SubjectID int    `json:"subject_id,omitempty"`
	Status    bool   `json:"status,omitempty"`
}
