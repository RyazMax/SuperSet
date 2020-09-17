package task_manager

import (
	"../models"
	"../repos/labeledtask"
	"../repos/schema"
	"../repos/task"
)

type TaskManager interface {
	Init(tr task.Repo, ltr labeledtask.Repo, sr schema.Repo) error

	PutTask(string, *models.Task) (int, error)
	TakeTask([]string) (*models.TaskWithSchema, error)
	LabelTask(*models.TaskWithSchema, *models.LabeledTask) error
}
