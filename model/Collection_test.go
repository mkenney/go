/*
Package model is a data-driven modeling abstraction
*/
package model

import (
	"fmt"
	"testing"
)

func makeFilteringCollection() *Collection {
	var dataMap map[interface{}]interface{}
	var model *Model
	collection := NewCollection()

	dataMap = make(map[interface{}]interface{})
	dataMap["a"] = "value1.1"
	dataMap["b"] = "value2.1"
	dataMap["c"] = "value3.1"
	dataMap["d"] = "value4.1"
	model = NewModel()
	model.SetData(dataMap)
	collection.Push(model)

	dataMap = make(map[interface{}]interface{})
	dataMap["a"] = "value1.2"
	dataMap["b"] = "value2.2"
	dataMap["c"] = "value3.2"
	dataMap["d"] = "value4.2"
	model = NewModel()
	model.SetData(dataMap)
	collection.Push(model)

	dataMap = make(map[interface{}]interface{})
	dataMap["a"] = "value1.3"
	dataMap["b"] = "value2.3"
	dataMap["c"] = "value3.3"
	dataMap["d"] = "value4.3"
	model = NewModel()
	model.SetData(dataMap)
	collection.Push(model)

	return collection
}

func TestNewCollection(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error

	collection := NewCollection()
	expect = 0
	actual = len(collection.Data())
	assert(t, expect, actual, err)
}

func TestCollectionPush(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error

	collection := NewCollection()

	model := NewModel()
	model.Set("foo", "bar")
	collection.Push(model)

	expect = 1
	actual = len(collection.Data())
	assert(t, expect, actual, err)

	model = NewModel()
	model.Set("foo", "bar")
	collection.Push(model)

	expect = 2
	actual = len(collection.Data())
	assert(t, expect, actual, err)
}

func TestCollectionPop(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error

	collection := NewCollection()

	model := NewModel()
	model.Set("foo", "bar")
	collection.Push(model)

	model = NewModel()
	model.Set("baz", 1)
	collection.Push(model)

	model, err = collection.Pop()
	expect = 1
	actual = len(collection.Data())
	assert(t, expect, actual, err)

	model, err = collection.Pop()
	expect = 0
	actual = len(collection.Data())
	assert(t, expect, actual, err)

	expect = fmt.Errorf(ErrorCollectionIsEmpty)
	model, actual = collection.Pop()
	assert(t, expect, actual, err)
}

func TestCollectionDelete(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error

	collection := NewCollection()

	model := NewModel()
	model.Set("foo", "bar")
	collection.Push(model)

	model = NewModel()
	model.Set("baz", 1)
	collection.Push(model)

	expect = 2
	actual = len(collection.Data())
	assert(t, expect, actual, err)

	expect = fmt.Errorf(ErrorCollectionIndexDoesNotExist, 2)
	_, actual = collection.Delete(2)
	assert(t, expect, actual, err)

	expect = true
	actual, err = collection.Delete(1)
	assert(t, expect, actual, err)

	expect = 1
	actual = len(collection.Data())
	assert(t, expect, actual, err)

	expect = true
	actual, err = collection.Delete(0)
	assert(t, expect, actual, err)

	expect = 0
	actual = len(collection.Data())
	assert(t, expect, actual, err)
}

func TestCollectionIsEmpty(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error

	collection := NewCollection()

	expect = true
	actual = collection.IsEmpty()
	assert(t, expect, actual, err)

	model := NewModel()
	model.Set("foo", "bar")
	collection.Push(model)

	expect = false
	actual = collection.IsEmpty()
	assert(t, expect, actual, err)

}

func TestCollectionFilter(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error
	var data []*Model

	collection := makeFilteringCollection()

	data = collection.Filter(func(key int, model *Model) bool {
		val, _ := model.Get("c")
		if "value3.2" == val {
			return true
		}
		return false
	})

	//
	expect = 1
	actual = len(data)
	assert(t, expect, actual, err)

	//
	err = nil
	expect = "value3.2"
	actual, _ = data[0].Get("c")
	assert(t, expect, actual, err)

	data = collection.Filter(func(key int, model *Model) bool {
		val, _ := model.Get("a")
		if "value1.1" == val {
			return true
		}
		return false
	})

	//
	expect = 1
	actual = len(data)
	assert(t, expect, actual, err)

	//
	err = nil
	expect = "value1.1"
	actual, _ = data[0].Get("a")
	assert(t, expect, actual, err)
}
