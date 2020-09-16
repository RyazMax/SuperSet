package task

import (
	"../../models"
)

type Repo interface {
	Init(host string, port int) error

	Insert(proj string, t *models.Task) error

	CreateTube(proj string) error
	DropTube(proj string) error

	TakeTask(projs []string) (*models.TaskAggr, error)
	AckTask(proj string, id int) error

	Drop()
}
