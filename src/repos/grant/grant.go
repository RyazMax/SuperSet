package grant

import "../../models"

// Repo smth
type Repo interface {
	Init(string, int) error
	Drop()

	Insert(*models.ProjectGrant) error
	GetByPairID(int, int) (*models.ProjectGrant, error)
	DeleteByPairID(int, int) error
	DeleteByProjectID(int) error
}
