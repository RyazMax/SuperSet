package universe

import (
	"log"

	"../auth"
	"../project_manager"
	"../repos/grant"
	"../repos/project"
	"../repos/schema"
	"../repos/session"
	"../repos/user"
)

// Universe is singleton object for this app
type Universe struct {
	UserRepo       user.Repo
	ProjectRepo    project.Repo
	SessionRepo    session.Repo
	GrantRepo      grant.Repo
	SchemaRepo     schema.Repo
	Auth           auth.Auth
	ProjectManager project_manager.ProjectManager
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

	grantRepo := grant.TarantoolRepo{}
	err = grantRepo.Init(host, port)
	if err != nil {
		return err
	}
	deferNotInited(grantRepo.Drop)

	schemaRepo := schema.TarantoolRepo{}
	err = schemaRepo.Init(host, port)
	if err != nil {
		return err
	}
	deferNotInited(schemaRepo.Drop)

	authInstance := auth.SimpleAuth{}
	err = authInstance.Init(&userRepo, &sessionRepo)
	if err != nil {
		return err
	}

	projectManager := project_manager.SimpleManager{}
	err = projectManager.Init(&userRepo, &projectRepo, &grantRepo, &schemaRepo)
	if err != nil {
		return err
	}

	single = &Universe{
		UserRepo:       &userRepo,
		ProjectRepo:    &projectRepo,
		SessionRepo:    &sessionRepo,
		GrantRepo:      &grantRepo,
		SchemaRepo:     &schemaRepo,
		Auth:           &authInstance,
		ProjectManager: &projectManager,
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
