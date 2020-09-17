package universe

import (
	"log"

	"../auth"
	"../project_manager"
	"../repos/grant"
	"../repos/labeledtask"
	"../repos/project"
	"../repos/schema"
	"../repos/session"
	"../repos/task"
	"../repos/user"
	"../task_manager"
)

// Universe is singleton object for this app
type Universe struct {
	UserRepo       user.Repo
	ProjectRepo    project.Repo
	SessionRepo    session.Repo
	GrantRepo      grant.Repo
	SchemaRepo     schema.Repo
	TaskRepo       task.Repo
	LabeledRepo    labeledtask.Repo
	Auth           auth.Auth
	ProjectManager project_manager.ProjectManager
	TaskManager    task_manager.TaskManager
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
	defer deferNotInited(grantRepo.Drop)

	schemaRepo := schema.TarantoolRepo{}
	err = schemaRepo.Init(host, port)
	if err != nil {
		return err
	}
	defer deferNotInited(schemaRepo.Drop)

	taskRepo := task.TarantoolRepo{}
	err = taskRepo.Init(host, port)
	if err != nil {
		return err
	}
	defer deferNotInited(taskRepo.Drop)

	labeledRepo := labeledtask.TarantoolRepo{}
	err = labeledRepo.Init(host, port)
	if err != nil {
		return err
	}
	defer deferNotInited(labeledRepo.Drop)

	authInstance := auth.SimpleAuth{}
	err = authInstance.Init(&userRepo, &sessionRepo)
	if err != nil {
		return err
	}

	projectManager := project_manager.SimpleManager{}
	err = projectManager.Init(&userRepo, &projectRepo, &grantRepo, &schemaRepo, &taskRepo)
	if err != nil {
		return err
	}

	taskManager := task_manager.SimpleManager{}
	err = taskManager.Init(&taskRepo, &labeledRepo, &schemaRepo)
	if err != nil {
		return err
	}

	single = &Universe{
		UserRepo:       &userRepo,
		ProjectRepo:    &projectRepo,
		SessionRepo:    &sessionRepo,
		GrantRepo:      &grantRepo,
		SchemaRepo:     &schemaRepo,
		TaskRepo:       &taskRepo,
		LabeledRepo:    &labeledRepo,
		Auth:           &authInstance,
		ProjectManager: &projectManager,
		TaskManager:    &taskManager,
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
