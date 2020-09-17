package task_manager

import (
	"errors"
	"log"

	"../models"
	"../repos/labeledtask"
	"../repos/schema"
	"../repos/task"
	"github.com/google/uuid"
)

// SimpleManager is simple implementation of task manager
type SimpleManager struct {
	TaskRepo    task.Repo
	LabeledRepo labeledtask.Repo
	SchemaRepo  schema.Repo
}

func (sm *SimpleManager) Init(tr task.Repo, ltr labeledtask.Repo, sr schema.Repo) error {
	sm.TaskRepo = tr
	sm.LabeledRepo = ltr
	sm.SchemaRepo = sr

	return nil
}

// PutTask put Task in Project proj tube and return ID of task
func (sm *SimpleManager) PutTask(proj string, t *models.Task) (id int, err error) {
	t.ID = int(uuid.New().ID())
	if t.ProjectID == 0 {
		panic(errors.New("Task should have project id"))
	}
	err = sm.TaskRepo.Insert(proj, t)
	if err != nil {
		log.Println("Can't insert task", err)
		return 0, err
	}

	return t.ID, nil
}

// TakeTask return task from tube with schema
func (sm *SimpleManager) TakeTask(projs []string) (*models.TaskWithSchema, error) {
	task, err := sm.TaskRepo.TakeTask(projs)
	if err != nil {
		log.Println("Can't take task", err)
		return nil, err
	}
	if task == nil {
		return nil, nil
	}

	log.Println(task)
	schema, err := sm.SchemaRepo.GetByProjectID(task.Tsk.ProjectID)
	if err != nil {
		log.Println("Can't take schema of task", err)
		return nil, err
	}
	if schema == nil {
		log.Panicf("Can't find schema[%d]", task.Tsk.ProjectID)
	}
	return &models.TaskWithSchema{Tsk: task, Schema: schema}, nil
}

// LabelTask ackquire task and put its labeled data into special space
func (sm *SimpleManager) LabelTask(tws *models.TaskAggr, lt *models.LabeledTask) error {
	err := sm.TaskRepo.AckTask(tws.Tsk.ProjectID, tws.ID)
	if err != nil {
		log.Println("Can't ack task", err)
		return err
	}

	lt.ProjectID = tws.Tsk.ProjectID
	lt.OriginID = tws.Tsk.ID
	err = sm.LabeledRepo.Insert(lt)
	if err != nil {
		log.Println("Can't insert labeled task", err)
		return err
	}

	return nil
}
