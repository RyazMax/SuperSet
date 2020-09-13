package user

import "../../models"

// Repo interface to interruct with data
type Repo interface {
	Init(string, int) error
	Drop()

	GetByID(int) (*models.User, error)
	GetByLogin(string) (*models.User, error)

	DeleteByID(int) error
	DeleteByLogin(string) error

	Insert(*models.User) error
}
