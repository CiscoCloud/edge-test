package main

import "github.com/stealthly/go-avro"

var _FullName_schema, _FullName_schema_err = avro.ParseSchema(`{
    "type": "record",
    "namespace": "com.example",
    "name": "FullName",
    "doc": "A persons first and last names",
    "fields": [
        {
            "name": "first",
            "doc": "first name",
            "type": "string"
        },
        {
            "name": "last",
            "default": "bar",
            "type": "string"
        }
    ]
}`)

/* A persons first and last names */
type FullName struct {
	/* first name */
	First string

	Last string
}

func NewFullName() *FullName {
	return &FullName{
		Last: "bar",
	}
}

func (this *FullName) Schema() avro.Schema {
	if _FullName_schema_err != nil {
		panic(_FullName_schema_err)
	}
	return _FullName_schema
}
