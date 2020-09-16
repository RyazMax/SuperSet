package project_manager

import (
	"../models"
	"../repos/grant"
	"../repos/project"
	"../repos/schema"
	"../repos/task"
	"../repos/user"
)

// ProjectAggr aggregate of required entities
type ProjectAggr struct {
	Project models.Project
	Schema  models.ProjectSchema
}

// ProjectManager manages all things about creating projects
type ProjectManager interface {
	Init(ur user.Repo, pr project.Repo, gr grant.Repo, sr schema.Repo, tr task.Repo) error
	Create(pa *ProjectAggr) (bool, error)
	DeleteByName(string, string) (bool, error)

	AddGrant(string, string, string) (bool, error)
	DeleteGrant(string, string, string) (bool, error)
}
