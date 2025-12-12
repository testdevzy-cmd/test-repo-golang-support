package interfaces

import (
	"context"

	"github.com/test-repo-golang-support/models"
)

// ProjectReader interface for project read operations
type ProjectReader interface {
	ReadProject(ctx context.Context, id string) (*models.Project, error)
	ReadAllProjects(ctx context.Context) (models.ProjectList, error)
}

// ProjectWriter interface for project write operations
type ProjectWriter interface {
	WriteProject(ctx context.Context, project *models.Project) error
	DeleteProject(ctx context.Context, id string) error
}

// ProjectReadWriter combines ProjectReader and ProjectWriter
// Demonstrates interface embedding
type ProjectReadWriter interface {
	ProjectReader // Embedded interface
	ProjectWriter // Embedded interface
}

// ProjectRepository is a comprehensive interface for project data access
type ProjectRepository interface {
	ProjectReadWriter // Embedded composite interface
	CountProjects(ctx context.Context) (int, error)
	ProjectExists(ctx context.Context, id string) (bool, error)
}


