package avro

import "github.com/stealthly/go-avro"

type Employee struct {
	Name string
	Boss *Employee
}

func NewEmployee() *Employee {
	return &Employee{}
}

func (this *Employee) Schema() avro.Schema {
	if _Employee_schema_err != nil {
		panic(_Employee_schema_err)
	}
	return _Employee_schema
}

// Generated by codegen. Please do not modify.
var _Employee_schema, _Employee_schema_err = avro.ParseSchema(`{
    "type": "record",
    "namespace": "avro",
    "name": "Employee",
    "fields": [
        {
            "name": "name",
            "type": "string"
        },
        {
            "name": "boss",
            "type": [
                "Employee",
                "null"
            ]
        }
    ]
}`)
