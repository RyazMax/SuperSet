package project_manager

import (
	"testing"

	"../models"
	"../repos/grant"
	"../repos/project"
	"../repos/schema"
	"../repos/user"
)

var (
	host        = "127.0.0.1"
	port        = 6666
	testProject = ProjectAggr{
		Project: models.Project{
			Name:    "Test",
			OwnerID: 2,
		},
		Schema: models.ProjectSchema{
			InputSchema:  models.TextInputSchema{},
			OutputSchema: models.TextOutputSchema{},
		},
	}
	testOwnerUser = models.User{
		ID:    2,
		Login: "Jack",
	}
	testUser = models.User{
		ID:    1,
		Login: "Max",
	}
)

func TestProjectManagers(t *testing.T) {
	managers := []ProjectManager{
		&SimpleManager{},
	}

	ur := &user.TarantoolRepo{}
	sr := &schema.TarantoolRepo{}
	gr := &grant.TarantoolRepo{}
	pr := &project.TarantoolRepo{}
	ur.Init(host, port)
	sr.Init(host, port)
	gr.Init(host, port)
	pr.Init(host, port)
	for _, manager := range managers {
		err := manager.Init(ur, pr, gr, sr)
		if err != nil {
			t.Errorf("Init failed on %T, with error %v", manager, err)
		}

		ur.Insert(&testOwnerUser)
		ok, err := manager.Create(&testProject)
		if err != nil {
			t.Errorf("Create on %T failed with error %v", manager, err)
		}
		if !ok {
			t.Errorf("Create on %T failed with not ok", manager)
		}

		ok, err = manager.Create(&testProject)
		if !ok {
			t.Errorf("Create after create on %T failed, with not ok", manager)
		}
		if err == nil {
			t.Errorf("Create after create on %T failed, with nil error", manager)
		}

		// AddUser to table
		ur.Insert(&testUser)
		ok, err = manager.AddGrant(testOwnerUser.Login, testProject.Project.Name, testUser.Login)
		if err != nil {
			t.Errorf("AddGrant on %T failed with error %v", manager, err)
		}
		if !ok {
			t.Errorf("AddGrant on %T failed with not ok result", manager)
		}

		ok, err = manager.DeleteGrant(testOwnerUser.Login, testProject.Project.Name, testUser.Login)
		if err != nil {
			t.Errorf("DeleteGrant on %T failed with error %v", manager, err)
		}
		if !ok {
			t.Errorf("DeleteGrant on %T failed with not ok result", manager)
		}

		ok, err = manager.AddGrant(testUser.Login, testProject.Project.Name, testOwnerUser.Login)
		if !ok {
			t.Errorf("AddGrant with wrong owner failed on %T, with no ok resp", manager)
		}
		if err == nil {
			t.Errorf("AddGrant with wrong owner failed on %T, expected not nil error", manager)
		}

		ok, err = manager.DeleteByName(testOwnerUser.Login, testProject.Project.Name)
		if err != nil {
			t.Errorf("DeleteByName failed on %T, with error %v", manager, err)
		}
		if !ok {
			t.Errorf("DeleteByName failed on %T, with not ok resp", manager)
		}

		err = ur.DeleteByLogin(testOwnerUser.Login)
		if err != nil {
			t.Errorf("Teardown failed on %T, error %v", manager, err)
		}
		err = ur.DeleteByLogin(testUser.Login)
		if err != nil {
			t.Errorf("Teardown failed on %T, error %v", manager, err)
		}
	}
}
