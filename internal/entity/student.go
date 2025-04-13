// Package entity defines main entities for business logic.
package entity

// Student represents a student in the educational institution.
type Student struct {
	ID      int    `json:"id"`
	GroupID int    `json:"group_id"`
	Name    string `json:"name"`
	Email   string `json:"email,omitempty"` // Email is omitted in responses as per requirements
}

// Group represents an academic group.
type Group struct {
	ID        int     `json:"id"`
	ParentID  *int    `json:"parent_id,omitempty"`
	Name      string  `json:"name"`
	SubGroups []Group `json:"subGroups,omitempty"`
}

// StudentCreateRequest represents request body for creating a student.
type StudentCreateRequest struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	GroupID int    `json:"group_id" validate:"required"`
}

// GroupCreateRequest represents request body for creating a group.
type GroupCreateRequest struct {
	Name     string `json:"name" validate:"required"`
	ParentID *int   `json:"parent_id"`
}
