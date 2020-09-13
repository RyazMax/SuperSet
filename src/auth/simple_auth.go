package auth

import (
	"log"

	"../models"
	"../repos/session"
	"../repos/user"
	"github.com/google/uuid"
)

// SimpleAuth simple implementation of auth
type SimpleAuth struct {
	userRepo    user.Repo
	sessionRepo session.Repo
}

// Init inits
func (auth *SimpleAuth) Init(ur user.Repo, sr session.Repo) error {
	auth.userRepo = ur
	auth.sessionRepo = sr
	return nil
}

// Login logins
func (auth *SimpleAuth) Login(login string, pass string) (*models.Session, error) {
	user, err := auth.userRepo.GetByLogin(login)
	if err != nil {
		log.Println("Can't login", err)
		return nil, err
	}

	if user == nil || user.PasswordShadow != pass {
		return nil, nil
	}

	session := &models.Session{ID: uuid.New().String(), UserLogin: login}
	err = auth.sessionRepo.Insert(session)
	if err != nil {
		log.Println("Can't login", err)
		return nil, err
	}

	return session, nil
}

// Logout logouts
func (auth *SimpleAuth) Logout(id string) error {
	err := auth.sessionRepo.DeleteByID(id)
	if err != nil {
		log.Println("Can't logout", err)
	}
	return err
}

// CheckSession checks
func (auth *SimpleAuth) CheckSession(id string) (*models.Session, error) {
	sess, err := auth.sessionRepo.GetByID(id)
	if err != nil {
		log.Println("Can't check session", err)
		return nil, err
	}
	return sess, nil
}

// NewUser new user
func (auth *SimpleAuth) NewUser(u *models.User) (*models.Session, error) {
	err := auth.userRepo.Insert(u)
	if err != nil {
		log.Println("Can't create user", err)
		return nil, err
	}

	session := &models.Session{ID: uuid.New().String(), UserLogin: u.Login}
	err = auth.sessionRepo.Insert(session)
	if err != nil {
		log.Println("Can't login", err)
		return nil, err
	}

	return session, nil
}

// DeleteUser deletes
func (auth *SimpleAuth) DeleteUser(login string) error {
	err := auth.userRepo.DeleteByLogin(login)
	if err != nil {
		log.Println("Can't delete", err)
	}
	return err
}
