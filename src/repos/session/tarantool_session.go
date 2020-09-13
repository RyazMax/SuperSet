package session

import (
	"fmt"
	"log"

	"../../models"
	"github.com/tarantool/go-tarantool"
)

// TarantoolRepo tarantool implementation of sessions.Repo
type TarantoolRepo struct {
	conn      *tarantool.Connection
	spaceName string
}

// Init inits tarantool repo connection
func (tr *TarantoolRepo) Init(host string, port int) (err error) {
	tr.conn, err = tarantool.Connect(fmt.Sprintf(host+":%d", port), tarantool.Opts{User: "go", Pass: "go"})
	tr.spaceName = "sessions"
	return err
}

// Drop closes connection
func (tr *TarantoolRepo) Drop() {
	tr.conn.Close()
}

// GetByID returns Session by it's ID
func (tr *TarantoolRepo) GetByID(id string) (*models.Session, error) {
	var s []models.Session
	err := tr.conn.SelectTyped(tr.spaceName, 0, 0, 1, tarantool.IterEq, []interface{}{id}, &s)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(s) == 0 {
		return nil, nil
	}

	return &s[0], nil
}

// DeleteByID deletes Session by it's ID
func (tr *TarantoolRepo) DeleteByID(id string) error {
	resp, err := tr.conn.Delete(tr.spaceName, 0, []interface{}{id})
	if err != nil {
		log.Println(err, resp)
	}
	return err
}

// Insert inserts Session into tarantool
func (tr *TarantoolRepo) Insert(sess *models.Session) error {
	var resp interface{}
	err := tr.conn.InsertTyped(tr.spaceName, sess, &resp)
	if err != nil {
		log.Println(err, resp)
		return err
	}
	return nil
}
