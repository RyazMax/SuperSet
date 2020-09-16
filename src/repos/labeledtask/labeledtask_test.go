package labeledtask

import (
	"testing"

	"../../models"
)

var (
	host      = "127.0.0.1"
	port      = 6666
	testTasks = []models.LabeledTask{
		{ID: 1, ProjectID: 1, OriginID: 1, Timestamp: 5},
		{ID: 2, ProjectID: 1, OriginID: 2, Timestamp: 4},
		{ID: 3, ProjectID: 2, OriginID: 3, Timestamp: 6},
	}
)

func TestTaskRepo(t *testing.T) {
	repos := []Repo{
		&TarantoolRepo{},
	}

	for _, repo := range repos {
		err := repo.Init(host, port)
		if err != nil {
			t.Errorf("Init on %T failed, with error %v", repo, err)
		}

		for i, task := range testTasks {
			err = repo.Insert(&task)
			if err != nil {
				t.Errorf("Insert failed on %T[%d], with error %v", repo, i, err)
			}
		}

		task, err := repo.GetByOriginID(1, 1)
		if err != nil {
			t.Errorf("GetByOriginID failed with %v", err)
		}
		if task == nil {
			t.Errorf("GetByOriginID failed, expected %v, got nil", testTasks[0])
		}

		tasks, err := repo.GetByProjectID(1)
		if err != nil {
			t.Errorf("GetByProjectID failed with %v", err)
		}
		if len(tasks) != 2 {
			t.Errorf("Expected 2 tasks, got %d", len(tasks))
		}

		tasks, err = repo.GetGreaterTime(1, 4)
		if err != nil {
			t.Errorf("GetByProjectID failed with %v", err)
		}
		if len(tasks) != 1 {
			t.Errorf("Expected 1 tasks, got %d", len(tasks))
		}

		for i, task := range testTasks {
			err = repo.DeleteByID(task.ID)
			if err != nil {
				t.Errorf("Delete failed on %T[%d], with error %v", repo, i, err)
			}
		}

	}
}
