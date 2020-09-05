package user

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
	tp.spaceName = "users"
	return err
}

// GetByID gets User by id
func (tp *TarantoolRepo) GetByID(id int) (*models.User, error) {
	var u []models.User
	err := tp.conn.SelectTyped(tp.spaceName, 0, 0, 1, tarantool.IterEq, []interface{}{id}, &u)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(u) == 0 {
		return nil, err
	}

	return &u[0], err
}

// GetByLogin gets User by login
func (tp *TarantoolRepo) GetByLogin(login string) (*models.User, error) {
	var u []models.User
	err := tp.conn.SelectTyped(tp.spaceName, "login", 0, 1, tarantool.IterEq, []interface{}{login}, &u)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(u) == 0 {
		return nil, nil
	}

	return &u[0], nil
}

// Insert inserts User to tarantool
func (tp *TarantoolRepo) Insert(u *models.User) error {
	var resp interface{}
	err := tp.conn.InsertTyped(tp.spaceName, u, &resp)
	if err != nil {
		log.Println(err, resp)
	}
	return err
}

// DeleteByID deletes user from space
func (tp *TarantoolRepo) DeleteByID(id int) error {
	_, err := tp.conn.Delete(tp.spaceName, 0, []interface{}{id})
	if err != nil {
		log.Println(err)
	}
	return err
}

// DeleteByLogin deletes user from space using its login
func (tp *TarantoolRepo) DeleteByLogin(login string) error {
	_, err := tp.conn.Delete(tp.spaceName, "login", []interface{}{login})
	if err != nil {
		log.Println(err)
	}
	return err
}
