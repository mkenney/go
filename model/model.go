package model

import (
	"encoding/json"
	"fmt"
)

/*
New returns a referece to a Model.
*/
func New() (*Model, error) {
	return &Model{
		make(map[string]interface{}),
		false,
	}, nil
}

/*
Model represents a data model.
*/
type Model struct {
	data   map[string]interface{}
	static bool
}

/*
Del removes a locally stored value by index. If the model is static or
the specified value does not exist return an error, otherwise return the
deleted value.
*/
func (model *Model) Del(idx string) (interface{}, error) {
	if model.IsStatic() {
		return nil, fmt.Errorf(ErrCannotModifyStaticModel)
	}
	if !model.Has(idx) {
		return nil, fmt.Errorf(ErrInvalidIndex, idx)
	}
	data := model.data[idx]
	delete(model.data, idx)
	return data, nil
}

/*
Get returns the value stored at the specified index, or an error if the
index doesn't exist.
*/
func (model *Model) Get(idx string) (interface{}, error) {
	if !model.Has(idx) {
		return nil, fmt.Errorf(ErrInvalidIndex, idx)
	}
	return model.data[idx], nil
}

/*
Has checks to see if if the specified index exists.
*/
func (model *Model) Has(idx string) bool {
	_, ok := model.data[idx]
	if ok {
		return true
	}
	return false
}

/*
IsStatic returns the static flag value.
*/
func (model *Model) IsStatic() bool {
	return model.static
}

/*
Lock marks a model as read-only and no futher modifications can be made
to the data.
*/
func (model *Model) Lock() {
	model.static = true
}

/*
MarshalJSON implements json.Marshaler.
*/
func (model *Model) MarshalJSON() ([]byte, error) {
	return json.Marshal(model.data)
}

/*
Set stores a value at the specified index.
*/
func (model *Model) Set(idx string, val interface{}) error {
	if model.IsStatic() {
		return fmt.Errorf(ErrCannotModifyStaticModel)
	}
	model.data[idx] = val
	return nil
}

/*
UnmarshalJSON implements json.Unmarshaler.
*/
func (model *Model) UnmarshalJSON(byts []byte) error {
	if model.IsStatic() {
		return fmt.Errorf(ErrCannotModifyStaticModel)
	}

	json.Unmarshal(byts, &model.data)
	for _, v := range model.data {
		getType(v)
	}
	return nil

	//return json.Unmarshal(byts, &model.data)
}

func getType(val interface{}) dataType {
	switch typ := val.(type) {
	case *Model:
		fmt.Printf("\n\nModel: TYP: %+v\n\n", typ)
		return typeModel
	case bool:
		fmt.Printf("\n\nbool: TYP: %+v\n\n", typ)
		return typeBool
	case complex128:
		fmt.Printf("\n\ncomplex128: TYP: %+v\n\n", typ)
		return typeComplex128
	case complex64:
		fmt.Printf("\n\ncomplex64: TYP: %+v\n\n", typ)
		return typeComplex64
	case float32:
		fmt.Printf("\n\nfloat32: TYP: %+v\n\n", typ)
		return typeFloat32
	case float64:
		fmt.Printf("\n\nfloat64: TYP: %+v\n\n", typ+1)
		return typeFloat64
	case int:
		fmt.Printf("\n\nint: TYP: %+v\n\n", typ)
		return typeInt
	case int16:
		fmt.Printf("\n\nint16: TYP: %+v\n\n", typ)
		return typeInt16
	case int32: // also matches `rune` data
		fmt.Printf("\n\nru: TYP: %+v\n\n", typ)
		return typeInt32
	case int64:
		fmt.Printf("\n\nint64: TYP: %+v\n\n", typ)
		return typeInt64
	case int8:
		fmt.Printf("\n\nint8: TYP: %+v\n\n", typ)
		return typeInt8
	case string:
		fmt.Printf("\n\nstring: TYP: %+v\n\n", typ)
		return typeString
	case uint:
		fmt.Printf("\n\nuint: TYP: %+v\n\n", typ)
		return typeUint
	case uint16:
		fmt.Printf("\n\nuint16: TYP: %+v\n\n", typ)
		return typeUint16
	case uint32:
		fmt.Printf("\n\nuint32: TYP: %+v\n\n", typ)
		return typeUint32
	case uint64:
		fmt.Printf("\n\nuint64: TYP: %+v\n\n", typ)
		return typeUint64
	case uint8: // also matches `byte` data
		fmt.Printf("\n\nby: TYP: %+v\n\n", typ)
		return typeUint8
	case uintptr:
		fmt.Printf("\n\nuintptr: TYP: %+v\n\n", typ)
		return typeUintptr
	default:
		fmt.Printf("\n\ndefault: TYP: %+v\n\n", typ)
		return typeNA
	}
}

type dataType int

const (
	typeNA dataType = iota
	typeModel
	typeBool
	typeComplex128
	typeComplex64
	typeFloat32
	typeFloat64
	typeInt
	typeInt16
	typeInt32 // also matches `rune` data
	typeInt64
	typeInt8
	typeString
	typeUint
	typeUint16
	typeUint32
	typeUint64
	typeUint8 // also matches `byte` data
	typeUintptr
)
