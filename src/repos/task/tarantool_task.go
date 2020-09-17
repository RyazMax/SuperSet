package task

import (
	"fmt"
	"log"

	"../../models"
	"github.com/tarantool/go-tarantool"
)

type TarantoolRepo struct {
	conn *tarantool.Connection
}

func (tr *TarantoolRepo) Init(host string, port int) (err error) {
	tr.conn, err = tarantool.Connect(fmt.Sprintf(host+":%d", port), tarantool.Opts{User: "go", Pass: "go"})
	return err
}

func (tr *TarantoolRepo) Insert(proj string, t *models.Task) error {
	resp, err := tr.conn.Call("insert_task", []interface{}{proj, t})
	if err != nil {
		log.Println(err, resp)
	}
	return err
}

func (tr *TarantoolRepo) CreateTube(proj string) error {
	resp, err := tr.conn.Call("create_tube", []interface{}{proj})
	if err != nil {
		log.Println(err, resp)
	}
	return err
}

func (tr *TarantoolRepo) DropTube(proj string) error {
	resp, err := tr.conn.Call("drop_tube", []interface{}{proj})
	if err != nil {
		log.Println(err, resp)
	}
	return err
}

func (tr *TarantoolRepo) TakeTask(projs []string) (*models.TaskAggr, error) {
	var t []models.TaskAggr
	err := tr.conn.CallTyped("take_aggr_by_projects", []interface{}{projs}, &t)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(t) == 0 {
		return nil, nil
	}
	if t[0].ID == -1 {
		return nil, nil
	}
	return &t[0], nil
}

func (tr *TarantoolRepo) AckTask(pid, id int) error {
	resp, err := tr.conn.Call("ack_task", []interface{}{pid, id})
	if err != nil {
		log.Println(err, resp)
	}
	return err
}

func (tr *TarantoolRepo) Drop() {
	tr.conn.Close()
}
