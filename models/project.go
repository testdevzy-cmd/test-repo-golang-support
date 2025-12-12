package models

import (
	"fmt"
	"time"
)

// Type alias for Project
type ProjectID = string
type ProjectList = []Project

// ProjectStatus represents the status of a project
type ProjectStatus string

// Project status constants
const (
	ProjectStatusActive   ProjectStatus = "active"
	ProjectStatusArchived ProjectStatus = "archived"
	ProjectStatusDraft    ProjectStatus = "draft"
)

// Project represents a project in the system
// Demonstrates struct composition through embedding
type Project struct {
	BaseEntity            // Embedded struct (composition)
	Timestamps            // Embedded struct for soft delete
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      ProjectStatus `json:"status"`
	OwnerID     UserID    `json:"owner_id"`
	OrgID       OrgID     `json:"org_id"`
}

// =====================================
// Value Receiver Methods on Project
// =====================================

// DisplayName returns the project display name (value receiver)
func (p Project) DisplayName() string {
	return p.Name
}

// IsActive checks if project is active (value receiver)
func (p Project) IsActive() bool {
	return p.Status == ProjectStatusActive
}

// IsArchived checks if project is archived (value receiver)
func (p Project) IsArchived() bool {
	return p.Status == ProjectStatusArchived
}

// String implements Stringer interface (value receiver)
func (p Project) String() string {
	return fmt.Sprintf("Project{ID: %s, Name: %s, Status: %s}", p.ID, p.Name, p.Status)
}

// =====================================
// Pointer Receiver Methods on Project
// =====================================

// UpdateName updates the project name (pointer receiver)
func (p *Project) UpdateName(name string) {
	p.Name = name
	p.UpdatedAt = time.Now()
}

// UpdateDescription updates the description (pointer receiver)
func (p *Project) UpdateDescription(desc string) {
	p.Description = desc
	p.UpdatedAt = time.Now()
}

// SetStatus sets the project status (pointer receiver)
func (p *Project) SetStatus(status ProjectStatus) {
	p.Status = status
	p.UpdatedAt = time.Now()
}

// Archive archives the project (pointer receiver)
func (p *Project) Archive() {
	p.Status = ProjectStatusArchived
	now := time.Now()
	p.DeletedAt = &now
	p.UpdatedAt = now
}

// Activate activates the project (pointer receiver)
func (p *Project) Activate() {
	p.Status = ProjectStatusActive
	p.DeletedAt = nil
	p.UpdatedAt = time.Now()
}

// =====================================
// Constructor Function
// =====================================

// NewProject creates a new Project with initialized fields
func NewProject(id, name, ownerID, orgID string) *Project {
	now := time.Now()
	return &Project{
		BaseEntity: BaseEntity{
			ID:        id,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:    name,
		Status:  ProjectStatusDraft,
		OwnerID: ownerID,
		OrgID:   orgID,
	}
}


