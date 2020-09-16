package task

import (
	"testing"

	"../../models"
	"../schema"
)

var (
	host     = "127.0.0.1"
	port     = 6666
	testTask = models.Task{
		ID:        1,
		ProjectID: 1,
		DataJSON:  "{text: 'Hello'}",
	}
	testSchema = models.ProjectSchema{
		ProjectID:    1,
		InputSchema:  models.TextInputSchema{},
		OutputSchema: models.TextOutputSchema{},
	}
	schemaRepo = schema.TarantoolRepo{}
)

func TestRepos(t *testing.T) {
	repos := []Repo{
		&TarantoolRepo{},
	}

	schemaRepo.Init(host, port)

	for _, repo := range repos {
		err := repo.Init(host, port)
		if err != nil {
			t.Errorf("Init failed on %T, with %v", repo, err)
		}

		err = repo.CreateTube("Tester")
		if err != nil {
			t.Errorf("CreateTube failed on %T, with %v", repo, err)
		}

		err = repo.Insert("Tester", &testTask)
		if err != nil {
			t.Errorf("Insert failed on %T, with error %v", repo, err)
		}

		schemaRepo.Insert(&testSchema)

		aggr, err := repo.TakeTask([]string{"Tester"})
		if err != nil {
			t.Errorf("TakeTask failed with %v", err)
		}
		if aggr == nil {
			t.Errorf("TakeTask failed expected not nil aggr")
		}

		schemaRepo.DeleteByProjectID(testSchema.ProjectID)
		err = repo.AckTask("Tester", aggr.ID)
		if err != nil {
			t.Errorf("AckTask failed with %v", err)
		}

		err = repo.DropTube("Tester")
		if err != nil {
			t.Errorf("DropTube failed with %v", err)
		}
	}
}
