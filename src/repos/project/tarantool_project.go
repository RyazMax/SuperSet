package project

import (
	"fmt"
	"log"

	"../../models"
	"github.com/tarantool/go-tarantool"
)

// TarantoolRepo struct to connect to Tarantool implements Repo
type TarantoolRepo struct {
	conn      *tarantool.Connection
	spaceName string
}

// Init connects to Tarantool and initialies Repo
func (tp *TarantoolRepo) Init(host string, port int) (err error) {
	tp.conn, err = tarantool.Connect(fmt.Sprintf(host+":%d", port), tarantool.Opts{User: "go", Pass: "go"})
	tp.spaceName = "projects"
	return err
}

// Drop closes connection
func (tp *TarantoolRepo) Drop() {
	tp.conn.Close()
}

// GetByID gets project by id
func (tp *TarantoolRepo) GetByID(id int) (*models.Project, error) {
	var p []models.Project
	err := tp.conn.SelectTyped(tp.spaceName, 0, 0, 1, tarantool.IterEq, []interface{}{id}, &p)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(p) == 0 {
		return nil, err
	}

	return &p[0], err
}

// GetByName gets project by name
func (tp *TarantoolRepo) GetByName(name string) (*models.Project, error) {
	var p []models.Project
	err := tp.conn.SelectTyped(tp.spaceName, "name", 0, 1, tarantool.IterEq, []interface{}{name}, &p)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(p) == 0 {
		return nil, err
	}

	return &p[0], err
}

// GetByOwnerID selects projects with same OwnerID
func (tp *TarantoolRepo) GetByOwnerID(id int) ([]models.Project, error) {
	var p []models.Project
	err := tp.conn.SelectTyped(tp.spaceName, "ownerid", 0, tarantool.KeyLimit, tarantool.IterEq, []interface{}{id}, &p)
	if err != nil {
		log.Println(err)
	}
	return p, err
}

// GetAllowed gets projects which are allowed to user with ID = id
func (tp *TarantoolRepo) GetAllowed(id int) ([]models.Project, error) {
	var p []models.Project
	err := tp.conn.CallTyped("get_allowed_projects", []interface{}{id}, &p)
	if err != nil {
		log.Println(err)
	}
	return p, err
}

// Insert inserts project
func (tp *TarantoolRepo) Insert(p *models.Project) error {
	var resp interface{}
	err := tp.conn.InsertTyped(tp.spaceName, p, &resp)
	if err != nil {
		log.Println(err, resp)
	}
	return err
}

// DeleteByID deletes project by ID
func (tp *TarantoolRepo) DeleteByID(id int) error {
	_, err := tp.conn.Delete(tp.spaceName, 0, []interface{}{id})
	if err != nil {
		log.Println(err)
	}
	return err
}

// DeleteByName deletes project by its name
func (tp *TarantoolRepo) DeleteByName(name string) error {
	_, err := tp.conn.Delete(tp.spaceName, "name", []interface{}{name})
	if err != nil {
		log.Println(err)
	}
	return err
}
