package models

// Task task for labeling
type Task struct {
	_msgpack struct{} `msgpack:",asArray"`

	ID        int
	ProjectID int // Used to find schema
	DataJSON  string
}

// TaskAggr have info for schema
type TaskAggr struct {
	_msgpack struct{} `msgpack:",asArray"`

	ID  int
	Tsk Task
}

// LabeledTask stored after labeling
type LabeledTask struct {
	_msgpack struct{} `msgpack:",asArray"`

	ID         int
	ProjectID  int
	OriginID   int
	AnswerJSON string
	Timestamp  uint64
}
