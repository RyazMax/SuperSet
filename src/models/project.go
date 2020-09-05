package models

// Project type describes project that includes dataset to be labeled
type Project struct {
	_msgpack struct{} `msgpack:",asArray"`

	ID      int
	Name    string
	OwnerID int
}

// IsEqual checks if projects are equal
func (p *Project) IsEqual(o *Project) bool {
	if o == nil {
		return false
	}
	return p.ID == o.ID && p.Name == o.Name && p.OwnerID == o.OwnerID
}
