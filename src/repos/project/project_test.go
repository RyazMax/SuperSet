package project

// This tests are for integration testing with tarantool
// That means it is required to start tarantool before via tarantool src/tarantool/init.lua
// Also that's why their not atomic

import (
	"testing"

	"../../models"
)

var repo Repo
var testProject = models.Project{
	ID:      1,
	Name:    "Mail.ru",
	OwnerID: 1,
}

var testProject2 = models.Project{
	ID:      2,
	Name:    "Yandex",
	OwnerID: 1,
}

func TestProjectRepos(t *testing.T) {
	repos := []Repo{
		&TarantoolRepo{},
	}

	for _, repo := range repos {
		err := repo.Init("127.0.0.1", 6666)
		if err != nil {
			t.Errorf("Init on %T failed with %v", repo, err)
		}

		err = repo.Insert(&testProject)
		if err != nil {
			t.Errorf("Insert on %T failed with %v", repo, err)
		}

		err = repo.Insert(&testProject2)
		if err != nil {
			t.Errorf("Insert on %T failed with %v", repo, err)
		}

		proj, err := repo.GetByID(testProject.ID)
		if err != nil {
			t.Errorf("GetByID on %T failed with %v", repo, err)
		}
		if !testProject.IsEqual(proj) {
			t.Errorf("GetByID on %T expected %v, got %v", repo, testProject, proj)
		}

		proj, err = repo.GetByName(testProject.Name)
		if err != nil {
			t.Errorf("GetByName on %T failed with %v", repo, err)
		}
		if !testProject.IsEqual(proj) {
			t.Errorf("GetByName on %T, expected %v, got %v", repo, testProject, proj)
		}

		projs, err := repo.GetByOwnerID(testProject.OwnerID)
		if err != nil {
			t.Errorf("GetByOwnerID on %T failed with %v", repo, err)
		}
		if !testProject.IsEqual(&projs[0]) || !testProject2.IsEqual(&projs[1]) {
			t.Errorf("GetByOwnerID on %T, expected %v, got %v", repo, []models.Project{testProject, testProject2}, projs)
		}

		err = repo.DeleteByID(testProject.ID)
		if err != nil {
			t.Errorf("DeleteByID on %T failed with %v", repo, err)
		}

		err = repo.DeleteByName(testProject2.Name)
		if err != nil {
			t.Errorf("DeleteByName on %T failed with %v", repo, err)
		}
	}
}
