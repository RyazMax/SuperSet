package user

// This tests are for integration testing with tarantool
// That means it is required to start tarantool before via tarantool src/tarantool/init.lua
// Also that's why their not atomic

import (
	"testing"

	"../../models"
)

var repo Repo
var testUser = models.User{
	ID:             1,
	Login:          "test",
	PasswordShadow: "qwerty",
	Type:           models.AdminUser,
}

func TestUserRepos(t *testing.T) {
	repos := []Repo{
		&TarantoolRepo{},
	}

	for _, repo := range repos {
		err := repo.Init("127.0.0.1", 6666)
		if err != nil {
			t.Errorf("Init on %T failed with %v", repo, err)
		}

		err = repo.Insert(&testUser)
		if err != nil {
			t.Errorf("Insert on %T failed with %v", repo, err)
		}

		user, err := repo.GetByID(int(testUser.ID))
		if err != nil {
			t.Errorf("GetByID on %T failed with %v", repo, err)
		}
		if !testUser.IsEqual(user) {
			t.Errorf("GetByID on %T expected %v, got %v", repo, testUser, user)
		}

		user, err = repo.GetByLogin(testUser.Login)
		if err != nil {
			t.Errorf("GetByLogin on %T failed with %v", repo, err)
		}
		if !testUser.IsEqual(user) {
			t.Errorf("GetByLogin on %T, expected %v, got %v", repo, testUser, user)
		}

		err = repo.DeleteByID(int(testUser.ID))
		if err != nil {
			t.Errorf("DeleteByID on %T failed with %v", repo, err)
		}

		err = repo.DeleteByLogin(testUser.Login)
		if err != nil {
			t.Errorf("DeleteByLogin on %T failed with %v", repo, err)
		}
	}
}
