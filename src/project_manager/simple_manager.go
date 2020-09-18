package project_manager

import (
	"errors"
	"log"

	"github.com/google/uuid"

	"../models"
	"../repos/grant"
	"../repos/project"
	"../repos/schema"
	"../repos/task"
	"../repos/user"
)

// SimpleManager is simple implementation of ProjectManager
type SimpleManager struct {
	userRepo    user.Repo
	projectRepo project.Repo
	schemaRepo  schema.Repo
	grantRepo   grant.Repo
	taskRepo    task.Repo
}

func (sm *SimpleManager) isOwner(uname string, oid int) *models.User {
	user, err := sm.userRepo.GetByID(oid)
	if err != nil {
		log.Println("Can't delete", err)
		return nil
	}
	if user == nil {
		log.Println("Owner of ", oid, " do not exists")
		return nil
	}
	if user.Login != uname {
		return nil
	}
	return user
}

// Init init
func (sm *SimpleManager) Init(ur user.Repo, pr project.Repo, gr grant.Repo, sr schema.Repo, tr task.Repo) error {
	sm.userRepo = ur
	sm.projectRepo = pr
	sm.schemaRepo = sr
	sm.grantRepo = gr
	sm.taskRepo = tr

	user, _ := sm.userRepo.GetByLogin("admin")
	_, err := sm.Create(&ProjectAggr{
		Project: models.Project{OwnerID: int(user.ID), Name: "Animals"},
		Schema: models.ProjectSchema{
			InputSchema:  models.ImageInputSchema{},
			OutputSchema: models.ClassOutputSchema{ClassNames: []string{"Кошка", "Cобака", "Корова"}},
		}})
	if err != nil {
		log.Println(err)
	}
	return nil
}

// Create creates
func (sm *SimpleManager) Create(pa *ProjectAggr) (bool, error) {
	existing, err := sm.projectRepo.GetByName(pa.Project.Name)
	if err != nil {
		return false, err
	}
	if existing != nil {
		return true, errors.New("Проект с таким именем уже существует")
	}

	owner, err := sm.userRepo.GetByID(pa.Project.OwnerID)
	if err != nil {
		return false, err
	}
	if owner == nil {
		return true, errors.New("Создатель проекта не существует")
	}

	pa.Project.ID = int(uuid.New().ID())
	pa.Schema.ProjectID = pa.Project.ID

	err = sm.projectRepo.Insert(&pa.Project)
	if err != nil {
		return false, err
	}
	err = sm.schemaRepo.Insert(&pa.Schema)
	if err != nil {
		log.Println("Can't insert schema", err)
		return false, err
	}
	err = sm.grantRepo.Insert(&models.ProjectGrant{ProjectID: pa.Project.ID, UserID: int(owner.ID)})
	if err != nil {
		log.Println("Can't insert grant", err)
		return false, err
	}
	err = sm.taskRepo.CreateTube(pa.Project.Name)
	if err != nil {
		log.Println("Can't insert tube", err)
		return false, err
	}

	return true, nil
}

// DeleteByName deletes project with pname and checks if user with uname can do this
func (sm *SimpleManager) DeleteByName(uname, pname string) (bool, error) {
	project, err := sm.projectRepo.GetByName(pname)
	if err != nil {
		log.Println("Can't delete", err)
		return false, err
	}
	if project == nil {
		return true, errors.New("Проект с таким именем не существует")
	}

	if sm.isOwner(uname, project.OwnerID) == nil {
		return true, errors.New("Нет прав для удаления проекта")
	}

	err = sm.projectRepo.DeleteByName(pname)
	if err != nil {
		return false, err
	}
	err = sm.schemaRepo.DeleteByProjectID(project.ID)
	if err != nil {
		return false, err
	}
	err = sm.grantRepo.DeleteByProjectID(project.ID)
	if err != nil {
		return false, err
	}
	err = sm.taskRepo.DropTube(pname)
	if err != nil {
		return false, err
	}

	return true, nil
}

// AddGrant tries to add grant for User uname to project Pname and check owner oname
func (sm *SimpleManager) AddGrant(oname, pname, uname string) (bool, error) {
	project, err := sm.projectRepo.GetByName(pname)
	if err != nil {
		log.Println("Can't delete", err)
		return false, err
	}
	if project == nil {
		return true, errors.New("Проект с таким именем не существует")
	}
	if sm.isOwner(oname, project.OwnerID) == nil {
		return true, errors.New("Нет прав для удаления проекта")
	}

	user, err := sm.userRepo.GetByLogin(uname)
	if err != nil {
		log.Println("Can't add", err)
		return false, err
	}
	if user == nil {
		return true, errors.New("Пользователь не существует")
	}

	existing, err := sm.grantRepo.GetByPairID(project.ID, int(user.ID))
	if err != nil {
		log.Println("Can't add", err)
		return false, err
	}
	if existing != nil {
		return true, errors.New("У пользователя уже есть права")
	}

	err = sm.grantRepo.Insert(&models.ProjectGrant{ProjectID: project.ID, UserID: int(user.ID)})
	if err != nil {
		return false, err
	}

	return true, nil
}

// DeleteGrant tries to delete grant of User uname to project Pname and check owner oname
func (sm *SimpleManager) DeleteGrant(oname, pname, uname string) (bool, error) {
	project, err := sm.projectRepo.GetByName(pname)
	if err != nil {
		log.Println("Can't delete", err)
		return false, err
	}
	if project == nil {
		return true, errors.New("Проект с таким именем не существует")
	}
	if sm.isOwner(oname, project.OwnerID) == nil {
		return true, errors.New("Нет прав для удаления проекта")
	}

	user, err := sm.userRepo.GetByLogin(uname)
	if err != nil {
		log.Println("Can't add", err)
		return false, err
	}
	if user == nil {
		return true, errors.New("Пользователь не существует")
	}

	existing, err := sm.grantRepo.GetByPairID(project.ID, int(user.ID))
	if err != nil {
		log.Println("Can't add", err)
		return false, err
	}
	if existing == nil {
		return true, errors.New("У пользователя нет прав")
	}

	err = sm.grantRepo.DeleteByPairID(project.ID, int(user.ID))
	if err != nil {
		return false, err
	}

	return true, nil
}

// Instance return instance of project aggr
func Instance() *ProjectAggr {
	return &ProjectAggr{}
}
