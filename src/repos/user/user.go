package user

const (
	ADMIN_USER = iota
	REGULAR_USER
)

type User struct {
	ID             int
	Login          string
	passwordShadow string
	Type           int
}

type UserRepo interface {
	Init(string, int) error

	GetById(int) (User, error)
	GetByLogin(string) (User, error)

	SetById(int) (User, error)
	SetByLogin(int) (User, error)
}
