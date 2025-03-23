package entity

import "time"

type Subject struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      []User    `json:"user"`
}

type CreateSubjectRequest struct {
	Name string `json:"name,omitempty"`
}

type UpdateSubjectRequest struct {
	Name   string `json:"name,omitempty"`
	Status bool   `json:"status,omitempty"`
}
