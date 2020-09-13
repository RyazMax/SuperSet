package models

// IInputSchema is interface for input schemas
type IInputSchema interface {
	InputType() string
}

const (
	TextInputType  = "PlainText"
	TableInputType = "TableData"
	ImageInputType = "Image"
)

// InputTypeToStructMap is mapping of string type to implementation of IInputSchema
func InputTypeToStructMap(t string) IInputSchema {
	switch t {
	case TextInputType:
		return TextInputSchema{}
	case TableInputType:
		return TableInputSchema{}
	case ImageInputType:
		return ImageInputSchema{}
	}
	return nil
}

// TextInputSchema plain text input schema
type TextInputSchema struct {
}

// InputType implements IInputSchema
func (ts TextInputSchema) InputType() string {
	return TextInputType
}

// TableInputSchema table values input implementation
type TableInputSchema struct {
	ColumnsNumber int
	ColumnsNames  []string
}

// InputType implements IInputSchema
func (ts TableInputSchema) InputType() string {
	return TableInputType
}

// ImageInputSchema is for images inputs
type ImageInputSchema struct {
}

// InputType implements IInputSchema
func (is ImageInputSchema) InputType() string {
	return ImageInputType
}
