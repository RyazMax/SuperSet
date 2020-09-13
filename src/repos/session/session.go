package session

import (
	"../../models"
)

// Repo represents repo interface for sessions
type Repo interface {
	Init(string, int) error
	Drop()

	GetByID(string) (*models.Session, error)
	DeleteByID(string) error
	Insert(*models.Session) error
}
