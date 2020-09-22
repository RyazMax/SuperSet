package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
)

// IInputSchema is interface for input schemas
type IInputSchema interface {
	InputType() string

	Validate(string, multipart.File) (*Task, error)
	SaveName(string, int) string

	FormatInputData(*Task, map[string]interface{})
	Init(r *http.Request) IInputSchema
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

func (ts TextInputSchema) Validate(name string, file multipart.File) (*Task, error) {
	ext := filepath.Ext(name)
	if ext != ".txt" {
		return nil, errors.New("Текстовая схема требует .txt файлы")
	}
	text, _ := ioutil.ReadAll(file)

	task := Task{
		DataJSON: string(text), // Not JSON but ok:)
	}
	return &task, nil
}

func (ts TextInputSchema) SaveName(name string, id int) string {
	return ""
}

func (ts TextInputSchema) FormatInputData(t *Task, data map[string]interface{}) {
	data["Text"] = t.DataJSON
}

func (ts TextInputSchema) Init(r *http.Request) IInputSchema {
	return ts
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

func (ts TableInputSchema) Validate(name string, file multipart.File) (*Task, error) {
	ext := filepath.Ext(name)
	if ext != ".csv" {
		return nil, errors.New("Текстовая схема требует .csv файлы")
	}
	text, _ := ioutil.ReadAll(file)

	// Do something with csv?

	task := Task{
		DataJSON: string(text), // Not JSON but ok:)
	}
	return &task, nil
}

func (ts TableInputSchema) SaveName(name string, id int) string {
	return ""
}

func (ts TableInputSchema) FormatInputData(t *Task, data map[string]interface{}) {
	data["ColNames"] = ts.ColumnsNames
	var s []string
	json.Unmarshal([]byte(t.DataJSON), &s)
	data["ColVals"] = s
}

func (ts TableInputSchema) Init(r *http.Request) IInputSchema {
	ts.ColumnsNumber, _ = strconv.Atoi(r.FormValue("col_count"))
	for i := 0; i < ts.ColumnsNumber; i++ {
		ts.ColumnsNames = append(ts.ColumnsNames, r.FormValue("column_"+strconv.Itoa(i+1)))
	}
	return ts
}

// ImageInputSchema is for images inputs
type ImageInputSchema struct {
}

// InputType implements IInputSchema
func (is ImageInputSchema) InputType() string {
	return ImageInputType
}

// Validate create new task on schema if possible
func (is ImageInputSchema) Validate(name string, file multipart.File) (*Task, error) {
	ext := filepath.Ext(name)
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return nil, errors.New("Текстовая схема требует .jpg файлы")
	}
	// Do something with csv?

	task := Task{DataJSON: name}
	return &task, nil
}

func (is ImageInputSchema) SaveName(name string, id int) string {
	return fmt.Sprintf("%s-%d", name, id)
}

func (is ImageInputSchema) FormatInputData(t *Task, data map[string]interface{}) {
	data["ImgSrc"] = fmt.Sprintf("%s-%d", t.DataJSON, t.ID)
}

func (is ImageInputSchema) Init(r *http.Request) IInputSchema {
	return is
}
