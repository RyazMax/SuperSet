package grant

import (
	"fmt"
	"log"

	"../../models"
	"github.com/tarantool/go-tarantool"
)

// TarantoolRepo tarantool implementation of repo
type TarantoolRepo struct {
	conn      *tarantool.Connection
	spaceName string
}

// Init connects to Tarantool and initialies Repo
func (tp *TarantoolRepo) Init(host string, port int) (err error) {
	tp.conn, err = tarantool.Connect(fmt.Sprintf(host+":%d", port), tarantool.Opts{User: "go", Pass: "go"})
	tp.spaceName = "project_grants"
	return err
}

// Drop closes connection
func (tp *TarantoolRepo) Drop() {
	tp.conn.Close()
}

// GetByPairID gets project by id
func (tp *TarantoolRepo) GetByPairID(pid, uid int) (*models.ProjectGrant, error) {
	var p []models.ProjectGrant
	err := tp.conn.SelectTyped(tp.spaceName, 0, 0, 1, tarantool.IterEq, []interface{}{pid, uid}, &p)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(p) == 0 {
		return nil, err
	}

	return &p[0], err
}

// Insert inserts project
func (tp *TarantoolRepo) Insert(p *models.ProjectGrant) error {
	var resp interface{}
	err := tp.conn.InsertTyped(tp.spaceName, p, &resp)
	if err != nil {
		log.Println(err, resp)
	}
	return err
}

// DeleteByPairID deletes project by ID
func (tp *TarantoolRepo) DeleteByPairID(pid, uid int) error {
	_, err := tp.conn.Delete(tp.spaceName, 0, []interface{}{pid, uid})
	if err != nil {
		log.Println(err)
	}
	return err
}

// DeleteByProjectID deletes all grants of choosen project
func (tp *TarantoolRepo) DeleteByProjectID(pid int) error {
	_, err := tp.conn.Call("delete_grants_by_project_id", []interface{}{pid})
	if err != nil {
		log.Println(err)
	}
	return err
}
