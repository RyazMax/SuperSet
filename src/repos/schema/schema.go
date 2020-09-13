package schema

import (
	"../../models"
)

// Repo is repo for project schemas
type Repo interface {
	Init(string, int) error
	Drop()

	GetByProjectID(int) (*models.ProjectSchema, error)
	DeleteByProjectID(int) error
	Insert(*models.ProjectSchema) error
}
