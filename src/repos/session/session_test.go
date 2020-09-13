package session

import (
	"testing"

	"../../models"
)

var (
	repo        Repo
	testSession = models.Session{
		ID:        "somehash",
		UserLogin: "UserLogin",
	}
)

func TestSessionRepos(t *testing.T) {
	repos := []Repo{
		&TarantoolRepo{},
	}

	for _, repo := range repos {
		err := repo.Init("127.0.0.1", 6666)
		if err != nil {
			t.Errorf("Init on %T failed with %v", repo, err)
		}

		err = repo.Insert(&testSession)
		if err != nil {
			t.Errorf("Insert on %T failed with %v", repo, err)
		}

		err = repo.Insert(&testSession)
		if err == nil {
			t.Errorf("Insert after insert %T failed on nil error", repo)
		}

		session, err := repo.GetByID(testSession.ID)
		if err != nil {
			t.Errorf("GetByID on %T failed with %v", repo, err)
		}
		if !testSession.IsEqual(session) {
			t.Errorf("GetByID on %T expected %v, got %v", repo, testSession, session)
		}

		err = repo.DeleteByID(testSession.ID)
		if err != nil {
			t.Errorf("DeleteByID on %T failed with %v", repo, err)
		}

		session, err = repo.GetByID(testSession.ID)
		if err != nil {
			t.Errorf("GetByID after delete on %T failed with %v", repo, err)
		}
		if session != nil {
			t.Errorf("GetByID after delete on %T expected %v, got %v", repo, nil, session)
		}
	}
}
