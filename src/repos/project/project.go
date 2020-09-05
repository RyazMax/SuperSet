package project

import "../../models"

// Repo interface to interruct with data
type Repo interface {
	Init(string, int) error

	GetByID(int) (*models.Project, error)
	GetByName(string) (*models.Project, error)
	GetByOwnerID(int) ([]models.Project, error)

	Insert(*models.Project) error

	DeleteByID(int) error
	DeleteByName(string) error
}
