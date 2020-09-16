package labeledtask

import (
	"../../models"
)

// Repo is interface
type Repo interface {
	Init(host string, port int) error

	Insert(*models.LabeledTask) error
	GetByOriginID(pid, oid int) (*models.LabeledTask, error)
	GetGreaterTime(pid int, ts uint64) ([]models.LabeledTask, error)
	GetByProjectID(pid int) ([]models.LabeledTask, error)

	DeleteByID(int) error

	Drop()
}
