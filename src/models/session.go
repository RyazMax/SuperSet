package models

// Session struct describes session of user
type Session struct {
	_msgpack struct{} `msgpack:",asArray"`

	ID        string
	UserLogin string
	// TODO Expires
}

// IsEqual checks if sessions are equal
func (s *Session) IsEqual(o *Session) bool {
	if o == nil {
		return false
	}
	return s.ID == o.ID && s.UserLogin == o.UserLogin
}
