package models

const (
	// AdminUser - user that has ability to create projects
	AdminUser = iota
	// RegularUser - just regular user
	RegularUser
)

// User is service user profile
type User struct {
	_msgpack struct{} `msgpack:",asArray"`

	ID             uint
	Login          string
	PasswordShadow string
	Type           int
}

// IsEqual checks weather users fields are equal
func (u *User) IsEqual(o *User) bool {
	if o == nil {
		return false
	}
	return u.ID == o.ID && u.Login == o.Login && u.PasswordShadow == o.PasswordShadow && u.Type == o.Type
}
