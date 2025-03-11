package dto

type CreateSubjectDto struct {
	Name string `json:"name"`
}

type UpdateSubjectDto struct {
	Name   string `json:"name,omitempty"`
	Status bool   `json:"status,omitempty"`
}
