package labeledtask

import (
	"fmt"
	"log"

	"../../models"
	"github.com/tarantool/go-tarantool"
)

// TarantoolRepo tarantool implementation of Repo
type TarantoolRepo struct {
	spaceName string
	conn      *tarantool.Connection
}

func (tr *TarantoolRepo) Init(host string, port int) (err error) {
	tr.conn, err = tarantool.Connect(fmt.Sprintf(host+":%d", port), tarantool.Opts{User: "go", Pass: "go"})
	tr.spaceName = "labeledtasks"
	return err
}

func (tr *TarantoolRepo) Insert(t *models.LabeledTask) error {
	var resp interface{}
	err := tr.conn.InsertTyped(tr.spaceName, t, &resp)
	if err != nil {
		log.Println(err, resp)
	}
	return err
}

func (tr *TarantoolRepo) GetByProjectID(pid int) ([]models.LabeledTask, error) {
	var t []models.LabeledTask
	err := tr.conn.SelectTyped(tr.spaceName, "ProjectID", 0, tarantool.KeyLimit, tarantool.IterEq, []interface{}{pid}, &t)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return t, nil
}

func (tr *TarantoolRepo) GetByOriginID(pid, oid int) (*models.LabeledTask, error) {
	var t []models.LabeledTask
	err := tr.conn.SelectTyped(tr.spaceName, "OriginID", 0, tarantool.KeyLimit, tarantool.IterEq, []interface{}{pid, oid}, &t)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(t) == 0 {
		return nil, nil
	}
	return &t[0], nil
}

func (tr *TarantoolRepo) GetGreaterTime(pid int, ts uint64) ([]models.LabeledTask, error) {
	var t []models.LabeledTask
	err := tr.conn.CallTyped("task_greater_ts", []interface{}{pid, ts}, &t)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return t, nil
}

func (tr *TarantoolRepo) DeleteByID(id int) error {
	resp, err := tr.conn.Delete(tr.spaceName, 0, []interface{}{id})
	if err != nil {
		log.Println(err, resp)
	}
	return err
}

func (tr *TarantoolRepo) Drop() {
	tr.conn.Close()
}
