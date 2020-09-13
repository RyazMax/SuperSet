package models

// ProjectGrant means that user has ability to label project
type ProjectGrant struct {
	_msgpack  struct{} `msgpack:",asArray"`
	ProjectID int
	UserID    int
}

// IsEqual checks equality
func (pg *ProjectGrant) IsEqual(other *ProjectGrant) bool {
	if other == nil {
		return false
	}
	return pg.ProjectID == other.ProjectID && pg.UserID == other.UserID
}
