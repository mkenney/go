/*
Package model is a data-driven modeling abstraction
*/
package model

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"testing"
)

type testInterface interface {
	add(a, b int) int
}

type testStruct struct {
	Name    string
	Options []string
}

func (ts *testStruct) add(a, b int) int {
	return a + b
}

func assert(t *testing.T, expect, actual interface{}, err error) {
	var eSlice bool
	rt := reflect.TypeOf(expect)
	switch rt.Kind() {
	case reflect.Slice:
		eSlice = true
	default:
		eSlice = false
	}

	if nil != err {
		t.Errorf("%s", err)
	} else if (eSlice && !reflect.DeepEqual(expect, actual)) && expect != actual {
		t.Fatalf("expected %v, %v found", expect, actual)
	}
}

func assertError(t *testing.T, expect, actual interface{}, err error) {
	if nil == err {
		t.Fatalf("expected an error, got nil")
	}
}

func triggerError(t *testing.T, err error) {
	t.Errorf("%s", err)
}

func makeFilteringModel() *Model {
	var dataMap map[string]string
	model := NewModel()

	dataMap = make(map[string]string)
	dataMap["a"] = "value1"
	dataMap["b"] = "value2"
	dataMap["c"] = "value3"
	dataMap["d"] = "value4"
	model.Set("data-1", dataMap)

	dataMap = make(map[string]string)
	dataMap["a"] = "value1"
	dataMap["b"] = "value1"
	dataMap["c"] = "value1"
	dataMap["d"] = "value1"
	model.Set("data-2", dataMap)

	dataMap = make(map[string]string)
	dataMap["a"] = "value2"
	dataMap["b"] = "value2"
	dataMap["c"] = "value2"
	dataMap["d"] = "value2"
	model.Set("data-2", dataMap)

	model.Set("data-3", "fooBarBaz00ntz")

	model.Set("data-4", 3.14159265)

	return model
}

func TestNew(t *testing.T) {
	var expect interface{}
	var actual interface{}

	model := NewModel()

	expect = 0
	actual = len(model.Data())
	assert(t, expect, actual, nil)

	expect = ""
	actual = model.UID()
	assert(t, expect, actual, nil)
}

func TestSetGet(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error

	model := NewModel()
	model.Set("var", "val")
	model.Set("int", 3)
	model.Set("float", 3.14)

	strct := new(testStruct)
	strct.Name = "Joe"
	strct.Options = make([]string, 10)
	model.Set("struct", strct)

	err = nil
	expect = "val"
	actual, err = model.Get("var")
	assert(t, expect, actual, err)

	err = nil
	expect = 3
	actual, err = model.Get("int")
	assert(t, expect, actual, err)

	err = nil
	expect = 3.14
	actual, err = model.Get("float")
	assert(t, expect, actual, err)

	err = nil
	expect = strct
	actual, err = model.Get("struct")
	assert(t, expect, actual, err)
}

func TestData(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error

	model := NewModel()
	model.Set("var", "val")
	model.Set("int", 3)
	model.Set("float", 3.14)
	strct := new(testStruct)
	strct.Name = "Joe"
	strct.Options = make([]string, 10)
	model.Set("struct", strct)

	data := model.Data()

	err = nil
	expect = "val"
	actual = data["var"]
	assert(t, expect, actual, err)

	err = nil
	expect = 3
	actual = data["int"]
	assert(t, expect, actual, err)

	err = nil
	expect = 3.14
	actual = data["float"]
	assert(t, expect, actual, err)

	err = nil
	expect = strct
	actual = data["struct"]
	assert(t, expect, actual, err)
}

func TestDeleteHas(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error

	model := NewModel()
	model.Set("var", "val")
	model.Set("int", 3)
	model.Set("float", 3.14)
	strct := new(testStruct)
	strct.Name = "Joe"
	strct.Options = make([]string, 10)
	model.Set("struct", strct)

	err = nil
	expect = true
	actual = model.Has("var")
	assert(t, expect, actual, err)

	err = nil
	expect = true
	actual = model.Has("int")
	assert(t, expect, actual, err)

	err = nil
	expect = true
	actual = model.Has("float")
	assert(t, expect, actual, err)

	err = nil
	expect = true
	actual = model.Has("struct")
	assert(t, expect, actual, err)

	err = nil
	expect = true
	actual, err = model.Delete("var")
	assert(t, expect, actual, err)
	err = nil
	expect = false
	actual = model.Has("var")
	assert(t, expect, actual, err)
	err = nil
	expect = false
	actual, err = model.Delete("var")
	assertError(t, expect, actual, err)

	err = nil
	expect = true
	actual, err = model.Delete("int")
	assert(t, expect, actual, err)
	err = nil
	expect = false
	actual = model.Has("int")
	assert(t, expect, actual, err)
	err = nil
	expect = false
	actual, err = model.Delete("int")
	assertError(t, expect, actual, err)

	err = nil
	expect = true
	actual, err = model.Delete("float")
	assert(t, expect, actual, err)
	err = nil
	expect = false
	actual = model.Has("float")
	assert(t, expect, actual, err)
	err = nil
	expect = false
	actual, err = model.Delete("float")
	assertError(t, expect, actual, err)

	err = nil
	expect = true
	actual, err = model.Delete("struct")
	assert(t, expect, actual, err)
	err = nil
	expect = false
	actual = model.Has("struct")
	assert(t, expect, actual, err)
	err = nil
	expect = false
	actual, err = model.Delete("struct")
	assertError(t, expect, actual, err)
}

func TestEmpty(t *testing.T) {
	var expect interface{}
	var actual interface{}

	model := NewModel()
	model.Set("emptyString", "")
	model.Set("emptyNil", nil)

	expect = true
	actual = model.Empty("emptyString")
	assert(t, expect, actual, nil)

	expect = true
	actual = model.Empty("emptyNil")
	assert(t, expect, actual, nil)

	model.Set("fullInt", 0)
	model.Set("fullString1", "0")
	model.Set("fullString2", " ")
	model.Set("fullFloat", 0.0)

	expect = false
	actual = model.Empty("fullInt")
	assert(t, expect, actual, nil)

	expect = false
	actual = model.Empty("fullString1")
	assert(t, expect, actual, nil)

	expect = false
	actual = model.Empty("fullString2")
	assert(t, expect, actual, nil)

	expect = false
	actual = model.Empty("fullFloat")
	assert(t, expect, actual, nil)

}

func TestFilter(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error
	var model *Model
	var data map[string]interface{}

	model = makeFilteringModel()
	data = model.Filter(func(key, val interface{}) bool {
		valMap, ok := val.(map[string]string)
		if ok {
			if "value3" == valMap["c"] {
				return true
			}
		}
		return false
	})
	//
	expect = 1
	actual = len(data)
	assert(t, expect, actual, nil)
	//
	err = nil
	expect = "value3"
	valMap, _ := data["data-1"].(map[string]string)
	actual = valMap["c"]
	assert(t, expect, actual, err)

	model = makeFilteringModel()
	data = model.Filter(func(key, val interface{}) bool {
		valMap, ok := val.(map[string]string)
		if ok {
			if "value2" == valMap["b"] {
				return true
			}
		}
		return false
	})
	//
	err = nil
	expect = 2
	actual = len(data)
	assert(t, expect, actual, err)
	//
	err = nil
	expect = "value2"
	valMap, _ = data["data-1"].(map[string]string)
	actual = valMap["b"]
	assert(t, expect, actual, err)
	//
	err = nil
	expect = "value2"
	valMap, _ = data["data-2"].(map[string]string)
	actual = valMap["b"]
	assert(t, expect, actual, err)

}

func TestStatic(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error
	var model *Model
	//var data map[string]interface{}

	model = makeFilteringModel()
	model.Set("Test", "Test")
	err = nil
	expect = false
	actual = model.IsStatic()
	assert(t, expect, actual, err)

	model.MakeStatic()
	err = nil
	expect = true
	actual = model.IsStatic()
	assert(t, expect, actual, err)

	err = nil
	expect = true
	actual = model.Set("newVal", "woot")
	if nil != actual {
		actual = true
	}
	assert(t, expect, actual, err)

	model2 := makeFilteringModel()
	err = nil
	expect = true
	actual = model.Merge(model2)
	if nil != actual {
		actual = true
	}
	assert(t, expect, actual, err)

	err = nil
	expect = true
	actual = model.Reset()
	if nil != actual {
		actual = true
	}
	assert(t, expect, actual, err)

	tmp := make(map[string]interface{})
	tmp["woot"] = 123
	err = nil
	expect = true
	actual = model.SetData(tmp)
	if nil != actual {
		actual = true
	}
	assert(t, expect, actual, err)

	err = nil
	expect = true
	actual = model.SetUID("new-uid")
	if nil != actual {
		actual = true
	}
	assert(t, expect, actual, err)

	err = nil
	expect = true
	actual = model.SetUIDType(ConstTypeInt)
	if nil != actual {
		actual = true
	}
	assert(t, expect, actual, err)

	err = nil
	expect = true
	actual = model.SetUIDType(ConstTypeString)
	if nil != actual {
		actual = true
	}
	assert(t, expect, actual, err)
}

func TestJSON(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error
	var model *Model

	model = makeFilteringModel()
	data := model.Data()
	bytes, _ := json.Marshal(data)
	expect = string(bytes)
	actual, err = model.JSON()
	assert(t, expect, actual, err)
}

func TestKeys(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error
	var model *Model

	model = makeFilteringModel()
	data := make([]string, len(model.Data()))
	for k := range model.Data() {
		data = append(data, k)
	}
	sort.Strings(data)
	expect = data
	actual = model.Keys()
	assert(t, expect, actual, err)
}

func TestLen(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error
	var model *Model

	model = makeFilteringModel()
	data := make([]string, len(model.Data()))
	expect = len(data)
	actual = model.Len()
	assert(t, expect, actual, err)
}

func TestModels(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error
	var model *Model

	// Not a real test...
	model = makeFilteringModel()
	data := make([]string, len(model.Data()))
	expect = len(data)
	actual = model.Len()
	assert(t, expect, actual, err)
}

func TestMerge(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error
	var model *Model

	// Not a real test...
	model = makeFilteringModel()
	moreModel := NewModel()
	moreModel.Set("MoreData", "somestring")
	model.Merge(moreModel)
	expect = "somestring"
	actual, err = model.Get("MoreData")
	assert(t, expect, actual, err)
}

func TestReduce(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error
	var model *Model

	// Not a real test...
	model = NewModel()
	model.Set("one", 1)
	model.Set("two", 2)
	model.Set("three", 3)
	model.Set("four", 4)
	model.Set("five", 5)

	expect = 15
	actual = model.Reduce(
		func(initial, value interface{}) interface{} {
			var realI, realV int
			switch i := initial.(type) {
			case int:
				realI = i
			}

			switch v := value.(type) {
			case int:
				realV = v
			}

			return realI + realV
		},
		nil)

	assert(t, expect, actual, err)
}

func TestReset(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error
	var model *Model

	// Not a real test...
	model = makeFilteringModel()

	model.Reset()
	expect = make(map[string]interface{})
	actual = model.Data()
	assert(t, expect, actual, err)
}

func TestSearch(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error
	var model *Model

	model = makeFilteringModel()

	expect = "data-3"
	actual, _ = model.Search("fooBarBaz00ntz")
	assert(t, expect, actual, err)
}

func TestSetData(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error
	var model *Model

	model = makeFilteringModel()

	data := make(map[string]interface{})
	data["w00t"] = 1
	model.SetData(data)
	expect = data
	actual = model.Data()
	assert(t, expect, actual, err)
}

func TestSetUID(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error
	var model *Model

	model = makeFilteringModel()

	model.SetUID("1")
	expect = "1"
	actual = model.UID()
	assert(t, expect, actual, err)
}

func TestSetUIDType(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error
	var model *Model

	model = makeFilteringModel()

	model.SetUIDType(ConstTypeInt)
	expect = ConstTypeInt
	actual = model.GetUIDType()
	assert(t, expect, actual, err)

	model.SetUIDType(ConstTypeString)
	expect = ConstTypeString
	actual = model.GetUIDType()
	assert(t, expect, actual, err)
}

func TestString(t *testing.T) {
	var expect interface{}
	var actual interface{}
	var err error
	var model *Model

	model = makeFilteringModel()

	expect, _ = model.JSON()
	actual = fmt.Sprintf("%s", model)
	assert(t, expect, actual, err)
}

func TestTransform(t *testing.T) {
	var expect map[string]interface{}
	var actual interface{}
	var err error
	var model *Model

	model = NewModel()
	model.Set("one", 1)
	model.Set("two", 2)
	model.Set("three", 3)
	model.Set("four", 4)
	model.Set("five", 5)

	expect = make(map[string]interface{})
	expect["one"] = 2
	expect["two"] = 3
	expect["three"] = 4
	expect["four"] = 5
	expect["five"] = 6
	actual = model.Transform(
		func(val interface{}) interface{} {
			typedVal := val.(int)
			return typedVal + 1
		})
	assert(t, expect, actual, err)
}

func TestUnique(t *testing.T) {
	var expect map[string]interface{}
	var actual interface{}
	var err error
	var model *Model

	model = NewModel()
	model.Set("one", 1)
	model.Set("two", 1)
	model.Set("three", 2)
	model.Set("four", 2)
	model.Set("five", 3)
	model.Unique()

	expect = make(map[string]interface{})
	expect["two"] = 1
	expect["three"] = 2
	expect["five"] = 3
	actual = model.Data()
	assert(t, expect, actual, err)
}

func TestWalk(t *testing.T) {
	var expect map[string]interface{}
	var actual interface{}
	var err error
	var model *Model

	model = NewModel()
	model.Set("one", 1)
	model.Set("two", 2)
	model.Set("three", 3)
	model.Set("four", 4)
	model.Set("five", 5)

	expect = make(map[string]interface{})
	expect["one"] = 2
	expect["two"] = 3
	expect["three"] = 4
	expect["four"] = 5
	expect["five"] = 6

	model.Walk(
		func(val interface{}) interface{} {
			typedVal := val.(int)
			return typedVal + 1
		})
	actual = model.Data()
	assert(t, expect, actual, err)
}
