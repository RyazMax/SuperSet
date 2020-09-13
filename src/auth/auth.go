package auth

import (
	"../models"
	"../repos/session"
	"../repos/user"
)

// Auth interface for something that make basic auth functions
type Auth interface {
	Init(ur user.Repo, sr session.Repo) error

	Login(login string, pass string) (*models.Session, error)
	Logout(id string) error

	CheckSession(id string) (*models.Session, error)

	NewUser(*models.User) (*models.Session, error)
	DeleteUser(login string) error
}
