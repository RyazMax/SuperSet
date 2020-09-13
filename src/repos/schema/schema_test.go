package schema

import (
	"testing"

	"../../models"
)

var testSchema = models.ProjectSchema{
	ProjectID:    1,
	InputSchema:  models.TextInputSchema{},
	OutputSchema: models.ClassOutputSchema{ClassNames: []string{"Cat", "Dog"}},
}

func TestSchemasRepos(t *testing.T) {
	repos := []Repo{
		&TarantoolRepo{},
	}

	for _, repo := range repos {
		err := repo.Init("127.0.0.1", 6666)
		if err != nil {
			t.Errorf("Init failed on %T, with error %v", repo, err)
		}

		err = repo.Insert(&testSchema)
		if err != nil {
			t.Errorf("Insert failed on %T, with error %v", repo, err)
		}

		schema, err := repo.GetByProjectID(testSchema.ProjectID)
		if err != nil {
			t.Errorf("GetByProjectID failed on %T, with error %v", repo, err)
		}

		if schema.ProjectID != testSchema.ProjectID {
			t.Errorf("GetByProjectID failed on %T, expected %v, got %v", repo, testSchema, schema)
		}

		err = repo.DeleteByProjectID(testSchema.ProjectID)
		if err != nil {
			t.Errorf("DeleteByProjectID failed on %T, with error %v", repo, err)
		}

		schema, err = repo.GetByProjectID(testSchema.ProjectID)
		if err != nil {
			t.Errorf("GetByProjectID after Delete failed on %T, with error %v", repo, err)
		}
		if schema != nil {
			t.Errorf("GetByProjectID after Delete failed on %T, expected nil, got %v", repo, schema)
		}
	}
}
