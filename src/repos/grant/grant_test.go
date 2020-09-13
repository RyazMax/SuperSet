package grant

import (
	"testing"

	"../../models"
)

var testGrant = models.ProjectGrant{
	ProjectID: 1,
	UserID:    3,
}

func TestGrantRepos(t *testing.T) {
	repos := []Repo{
		&TarantoolRepo{},
	}

	for _, repo := range repos {
		err := repo.Init("127.0.0.1", 6666)
		if err != nil {
			t.Errorf("Init failed on %T, with error %v", repo, err)
		}

		err = repo.Insert(&testGrant)
		if err != nil {
			t.Errorf("Insert failed on %T, with error %v", repo, err)
		}

		grant, err := repo.GetByPairID(testGrant.ProjectID, testGrant.UserID)
		if err != nil {
			t.Errorf("GetByPairID failed on %T, with error %v", repo, err)
		}
		if !testGrant.IsEqual(grant) {
			t.Errorf("GetByPairID failed on %T, expected %v, got %v", repo, testGrant, grant)
		}

		err = repo.DeleteByPairID(testGrant.ProjectID, testGrant.UserID)
		if err != nil {
			t.Errorf("DeleteByPairID failed on %T, with error %v", repo, err)
		}

		grant, err = repo.GetByPairID(testGrant.ProjectID, testGrant.UserID)
		if err != nil {
			t.Errorf("GetByPairID after Delete failed on %T, with error %v", repo, err)
		}
		if grant != nil {
			t.Errorf("GetByPairID after Delete failed on %T, expected nil grant, got %v", repo, grant)
		}
	}
}
