package models

// ProjectSchema is schema of project DataSet
type ProjectSchema struct {
	_msgpack struct{} `msgpack:",asArray"`

	ProjectID    int
	InputSchema  IInputSchema
	OutputSchema IOutputSchema
}
