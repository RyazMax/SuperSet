package user

import (
	"fmt"

	"github.com/tarantool/go-tarantool"
)

type TarantoolUserRepo struct {
	conn      *tarantool.Connection
	spaceName string
}

func (tp *TarantoolUserRepo) Init(host string, port int) (err error) {
	tp.conn, err = tarantool.Connect(fmt.Sprintf(host+":%d", port), tarantool.Opts{})
	tp.spaceName = "users"
	return err
}

func (tp *TarantoolUserRepo) GetById(id int) (User, error) {
	var u User
	err := tp.conn.SelectTyped(tp.spaceName, 0, 0, 0, 0, id, u)
	if err != nil {
		return u, err
	}

	return u, nil
}

func (tp *TarantoolUserRepo) GetByLogin(login string) (User, error) {
	var u User
	err := tp.conn.SelectTyped(tp.spaceName, "login", 0, 0, 0, login, u)
	if err != nil {
		return u, err
	}

	return u, nil
}
