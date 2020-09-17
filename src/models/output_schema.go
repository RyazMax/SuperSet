package models

import "net/http"

// IOutputSchema describes output schema of dataset
type IOutputSchema interface {
	OutputType() string

	FormatOutputData(*Task, map[string]interface{})
	FormatLabeledTask(r *http.Request) (*LabeledTask, error)
}

const (
	IntOutputType   = "Integer"
	FloatOutputType = "Float"
	ClassOutputType = "ClassData"
	TextOutputType  = "PlainText"
)

// OutputTypeToStructMap is mapping of string type to implementation of IOutputSchema
func OutputTypeToStructMap(t string) IOutputSchema {
	switch t {
	case IntOutputType:
		return IntOutputSchema{}
	case FloatOutputType:
		return FloatOutputSchema{}
	case ClassOutputType:
		return ClassOutputSchema{}
	case TextOutputType:
		return TextOutputSchema{}
	}
	return nil
}

// IntOutputSchema is struct for integer regression tasks
type IntOutputSchema struct {
	IsLimited   bool
	BottomLimit int
	TopLimit    int
}

// OutputType implementation of IOutputSchema
func (is IntOutputSchema) OutputType() string {
	return IntOutputType
}

func (is IntOutputSchema) FormatOutputData(t *Task, data map[string]interface{}) {
}

func (is IntOutputSchema) FormatLabeledTask(r *http.Request) (*LabeledTask, error) {
	t := LabeledTask{
		AnswerJSON: r.FormValue("IntegerData"),
	}
	return &t, nil
}

// FloatOutputSchema is struct for float regression tasks
type FloatOutputSchema struct {
	IsLimited   bool
	BottomLimit float64
	TopLimit    float64
}

// OutputType implementatino of IOutputSchema
func (fs FloatOutputSchema) OutputType() string {
	return FloatOutputType
}

func (fs FloatOutputSchema) FormatOutputData(t *Task, data map[string]interface{}) {
}

func (fs FloatOutputSchema) FormatLabeledTask(r *http.Request) (*LabeledTask, error) {
	t := LabeledTask{
		AnswerJSON: r.FormValue("FloatData"),
	}
	return &t, nil
}

// ClassOutputSchema for classification task
type ClassOutputSchema struct {
	ClassNames []string
}

// OutputType implementation of IOutputSchema
func (cs ClassOutputSchema) OutputType() string {
	return ClassOutputType
}

func (cs ClassOutputSchema) FormatOutputData(t *Task, data map[string]interface{}) {
	data["Classes"] = cs.ClassNames
}

func (cs ClassOutputSchema) FormatLabeledTask(r *http.Request) (*LabeledTask, error) {
	t := LabeledTask{
		AnswerJSON: r.FormValue("ClassLabel"),
	}
	return &t, nil
}

// TextOutputSchema for text answers
type TextOutputSchema struct {
	MinLength int
	MaxLength int
}

// OutputType implementation of IOutputSchema
func (ts TextOutputSchema) OutputType() string {
	return TextOutputType
}

func (ts TextOutputSchema) FormatOutputData(t *Task, data map[string]interface{}) {
}

func (ts TextOutputSchema) FormatLabeledTask(r *http.Request) (*LabeledTask, error) {
	t := LabeledTask{
		AnswerJSON: r.FormValue("TextData"),
	}
	return &t, nil
}
