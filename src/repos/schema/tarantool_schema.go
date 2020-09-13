package schema

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"../../models"
	"github.com/tarantool/go-tarantool"
)

// TarantoolRepo tarantool implementatino of Repo
type TarantoolRepo struct {
	conn      *tarantool.Connection
	spaceName string
}

type tarantoolSchema struct {
	_msgpack   struct{} `msgpack:",asArray"`
	ProjectID  int
	InputType  string
	InputJSON  string
	OutputType string
	OutputJSON string
}

func tarantoolSchemaFromSchema(ps *models.ProjectSchema) (*tarantoolSchema, error) {
	ts := tarantoolSchema{
		ProjectID:  ps.ProjectID,
		InputType:  ps.InputSchema.InputType(),
		OutputType: ps.OutputSchema.OutputType(),
	}

	data, err := json.Marshal(ps.InputSchema)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	ts.InputJSON = string(data)

	data, err = json.Marshal(ps.OutputSchema)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	ts.OutputJSON = string(data)

	return &ts, nil
}

func schemaFromTarantoolSchema(ts *tarantoolSchema) (*models.ProjectSchema, error) {
	// Magic with reflection whooo
	ps := models.ProjectSchema{
		ProjectID: ts.ProjectID,
	}

	iSchema := models.InputTypeToStructMap(ts.InputType)
	oSchema := models.OutputTypeToStructMap(ts.OutputType)
	smth := reflect.New(reflect.TypeOf(iSchema)).Interface().(models.IInputSchema)
	err := json.Unmarshal([]byte(ts.InputJSON), smth)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	ps.InputSchema = smth
	smth2 := reflect.New(reflect.TypeOf(oSchema)).Interface().(models.IOutputSchema)
	err = json.Unmarshal([]byte(ts.OutputJSON), &smth2)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	ps.OutputSchema = smth2
	return &ps, nil
}

// Init inits
func (tr *TarantoolRepo) Init(host string, port int) (err error) {
	tr.conn, err = tarantool.Connect(fmt.Sprintf(host+":%d", port), tarantool.Opts{User: "go", Pass: "go"})
	tr.spaceName = "project_schemas"
	return err
}

// GetByProjectID gets
func (tr *TarantoolRepo) GetByProjectID(id int) (*models.ProjectSchema, error) {
	var ts []tarantoolSchema
	err := tr.conn.SelectTyped(tr.spaceName, 0, 0, 1, tarantool.IterEq, []interface{}{id}, &ts)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(ts) == 0 {
		return nil, nil
	}

	return schemaFromTarantoolSchema(&ts[0])
}

// Drop drops
func (tr *TarantoolRepo) Drop() {
	tr.conn.Close()
}

// Insert inserts
func (tr *TarantoolRepo) Insert(s *models.ProjectSchema) error {
	ts, err := tarantoolSchemaFromSchema(s)
	if err != nil {
		return err
	}
	var resp interface{}
	err = tr.conn.InsertTyped(tr.spaceName, ts, &resp)
	if err != nil {
		log.Println(err, resp)
	}
	return err
}

// DeleteByProjectID deletes
func (tr *TarantoolRepo) DeleteByProjectID(id int) error {
	resp, err := tr.conn.Delete(tr.spaceName, 0, []interface{}{id})
	if err != nil {
		log.Println(err, resp)
	}
	return err
}
