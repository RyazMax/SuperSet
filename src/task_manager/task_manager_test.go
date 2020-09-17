package task_manager

import (
	"testing"

	"../models"
	"../repos/labeledtask"
	"../repos/project"
	"../repos/schema"
	"../repos/task"
)

var (
	host     = "127.0.0.1"
	port     = 6666
	testTask = models.Task{
		ProjectID: 1,
		DataJSON:  "{\"text\": \"smth\"}",
	}
	testLabel   = models.LabeledTask{}
	testProject = models.Project{
		ID:   1,
		Name: "Tester",
	}
	testSchema = models.ProjectSchema{
		ProjectID:    1,
		InputSchema:  models.TextInputSchema{},
		OutputSchema: models.ClassOutputSchema{},
	}
)

func TestManager(t *testing.T) {
	managers := []TaskManager{
		&SimpleManager{},
	}

	tr := task.TarantoolRepo{}
	ltr := labeledtask.TarantoolRepo{}
	sr := schema.TarantoolRepo{}
	pr := project.TarantoolRepo{}

	tr.Init(host, port)
	ltr.Init(host, port)
	sr.Init(host, port)
	pr.Init(host, port)

	for _, manager := range managers {
		err := manager.Init(&tr, &ltr, &sr)
		if err != nil {
			t.Errorf("Error on init %v", err)
		}

		pr.Insert(&testProject)
		sr.Insert(&testSchema)
		tr.CreateTube("Tester")

		_, err = manager.PutTask("Tester", &testTask)
		if err != nil {
			t.Errorf("PutTask failed with %v", err)
		}

		tws, err := manager.TakeTask([]string{"Tester"})
		if err != nil {
			t.Errorf("TakeTask failed with %v", err)
		}
		if tws == nil {
			t.Errorf("TakeTask failed, expected not nil TaskWithSchema")
		}
		if tws.Schema == nil {
			t.Errorf("TakeTask failed, expected not nil Schema")
		}

		err = manager.LabelTask(tws, &testLabel)
		if err != nil {
			t.Errorf("LabelTask failed with %v", err)
		}

		tr.DropTube("Tester")
		pr.DeleteByID(testProject.ID)
		sr.DeleteByProjectID(testSchema.ProjectID)
	}
}
