package models

// IOutputSchema describes output schema of dataset
type IOutputSchema interface {
	OutputType() string
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

// ClassOutputSchema for classification task
type ClassOutputSchema struct {
	ClassNames []string
}

// OutputType implementation of IOutputSchema
func (cs ClassOutputSchema) OutputType() string {
	return ClassOutputType
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
