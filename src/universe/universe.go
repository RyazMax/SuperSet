package universe

import (
	"log"

	"../auth"
	"../repos/project"
	"../repos/session"
	"../repos/user"
)

// Universe is singleton object for this app
type Universe struct {
	UserRepo    user.Repo
	ProjectRepo project.Repo
	SessionRepo session.Repo
	Auth        auth.Auth
}

var single *Universe
var inited bool

func deferNotInited(f func()) {
	if !inited {
		f()
	}
}

// Init inits
func Init(host string, port int) error {
	userRepo := user.TarantoolRepo{}
	err := userRepo.Init(host, port)
	if err != nil {
		return err
	}
	defer deferNotInited(userRepo.Drop)

	projectRepo := project.TarantoolRepo{}
	err = projectRepo.Init(host, port)
	if err != nil {
		return err
	}
	defer deferNotInited(projectRepo.Drop)

	sessionRepo := session.TarantoolRepo{}
	err = sessionRepo.Init(host, port)
	if err != nil {
		return err
	}
	defer deferNotInited(sessionRepo.Drop)

	authInstance := auth.SimpleAuth{}
	err = authInstance.Init(&userRepo, &sessionRepo)
	if err != nil {
		return err
	}

	single = &Universe{
		UserRepo:    &userRepo,
		ProjectRepo: &projectRepo,
		SessionRepo: &sessionRepo,
		Auth:        &authInstance,
	}
	inited = true
	return nil
}

// Get return single
func Get() *Universe {
	if single == nil {
		log.Fatal("Universe is Not inited, you probably forgot to init it")
		return nil
	}

	return single
}
