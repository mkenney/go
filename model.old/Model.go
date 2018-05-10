/*
Package model is a data-driven modeling abstraction
*/
package model

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

/*
ConstIdentifierTypeInt defines the int type for the identifier
*/
const ConstIdentifierTypeInt = 0

/*
ConstIdentifierTypeString defines the string type for the identifier
*/
const ConstIdentifierTypeString = 1

/*
Model is a struct that models arbitrary data
*/
type Model struct {

	/*
		Local data storage.
	*/
	data map[interface{}]interface{}

	/*
		Optional, identifier of this data object.
		This is used as the index value in model collections
	*/
	identifier string

	/*
		The data type of the identifier value.
		Either ConstIdentifierTypeInt or ConstIdentifierTypeString
	*/
	identifierType int

	/*
	   The read vs read/write mode.
	   If true, methods that modify data should fail
	*/
	isStatic bool

	/*
	   The view pointer to use for filtering data output
	*/
	view *View
}

/*
NewModel initializes and returns a pointer to a Model
*/
func NewModel() (model *Model) {
	model = new(Model)
	model.data = make(map[interface{}]interface{})
	model.view = new(View)
	return model
}

/*
Data returns the internal data map
*/
func (ma *Model) Data() map[interface{}]interface{} {
	return ma.data
}

/*
Delete removes a locally stored value by index
*/
func (ma *Model) Delete(idx string) (bool, error) {
	if ma.IsStatic() {
		return false, fmt.Errorf("static models cannot be modified")
	}
	if !ma.Has(idx) {
		return false, fmt.Errorf("the specified index '%v' does not exist", idx)
	}

	delete(ma.data, idx)

	return true, nil
}

/*
IsEmpty checks to see if a value should be considered "empty"
True if the value does not exist, is nil or is an empty string
*/
func (ma *Model) IsEmpty(idx string) bool {
	_, ok := ma.data[idx]
	if !ok || "" == ma.data[idx] || nil == ma.data[idx] {
		return true
	}
	return false
}

/*
Filter filters elements of the data using a callback function and returns the
result
*/
func (ma *Model) Filter(callback func(key, val interface{}) bool) map[interface{}]interface{} {
	retVal := make(map[interface{}]interface{})
	for k, v := range ma.data {
		if callback(k, v) {
			retVal[k] = v
		}
	}
	return retVal
}

/*
Get returns a locally stored value by index
*/
func (ma *Model) Get(idx string) (interface{}, error) {
	if !ma.Has(idx) {
		return nil, fmt.Errorf("the specified index '%s' does not exist", idx)
	}
	return ma.data[idx], nil
}

/*
GetUIDType returns a locally stored identifiers type
*/
func (ma *Model) GetUIDType() int {
	return ma.identifierType
}

/*
Has checks to see if a value has been set
*/
func (ma *Model) Has(idx string) bool {
	_, ok := ma.data[idx]
	if ok {
		return true
	}
	return false
}

/*
Import imports a JSON document as model data
*/
func (ma *Model) Import(jsonString string) (bool, error) {
	var unmarshalledData interface{}

	// Accept "structured" JSON strings
	jsonString = strings.Replace(jsonString, "\n", "", -1)
	jsonString = strings.Replace(jsonString, "\t", "", -1)

	bytes := []byte(jsonString)
	err := json.Unmarshal(bytes, &unmarshalledData)
	if nil != err {
		return false, fmt.Errorf(ErrorCannotDecodeJsonString)
	}

	result, err := importHelper(unmarshalledData)
	fmt.Printf("result: %v\n", result)
	ma.SetData(result)

	return true, err
}
func importHelper(data interface{}) (map[interface{}]interface{}, error) {
	var err error
	var modelData map[interface{}]interface{}
	retVal := make(map[interface{}]interface{})

	switch typedData := data.(type) {
	case map[string]interface{}:
		for k, v := range typedData {
			switch v.(type) {
			case map[string]interface{}:
				modelData, err = importHelper(v)
				model := NewModel()
				model.SetData(modelData)
				retVal[k] = model
			case []interface{}:
				modelData, err = importHelper(v)
				model := NewModel()
				model.SetData(modelData)
				retVal[k] = model
			case string:
				retVal[k] = v
			case float64:
				retVal[k] = v
			case int:
				retVal[k] = v
			case bool:
				retVal[k] = v
			case nil:
				retVal[k] = v
			}
		}
	case []interface{}:
		for k, v := range typedData {
			switch v.(type) {
			case map[string]interface{}:
				modelData, err = importHelper(v)
				model := NewModel()
				model.SetData(modelData)
				retVal[k] = model
			case []interface{}:
				modelData, err = importHelper(v)
				model := NewModel()
				model.SetData(modelData)
				retVal[k] = model
			case string:
				retVal[k] = v
			case float64:
				retVal[k] = v
			case int:
				retVal[k] = v
			case bool:
				retVal[k] = v
			case nil:
				retVal[k] = v
			}
		}
	case string:
		retVal[0] = typedData
	case float64:
		retVal[0] = typedData
	case int:
		retVal[0] = typedData
	case bool:
		retVal[0] = typedData
	case nil:
		retVal[0] = typedData
	default:
		err = fmt.Errorf("unknown data structure found, some data was not imported")
	}

	return retVal, err
}

/*
IsStatic returns whether the static flag has been set or not
*/
func (ma *Model) IsStatic() bool {
	return ma.isStatic
}

/*
JSON recursively converts the internal data storage array to a JSON string
Exposed as a public method to give access to the json_encode() parameters
*/
func (ma *Model) JSON() (string, error) {
	strings, slices := jsonHelper(ma)
	var bytes []byte
	var err error

	if 0 != len(strings) {
		bytes, err = json.Marshal(strings)
	} else {
		bytes, err = json.Marshal(slices)
	}
	return string(bytes), err
}
func jsonHelper(ma *Model) (map[string]interface{}, []interface{}) {
	jsonMap := make(map[string]interface{})
	jsonSlice := make([]interface{}, 0)

	for k, v := range ma.Data() {
		switch v.(type) {
		case *Model:
			model, _ := v.(*Model)
			switch keyType := k.(type) {
			case string:
				strings, slices := jsonHelper(model)
				if 0 != len(strings) {
					jsonMap[keyType] = strings
				}
				if 0 != len(slices) {
					jsonMap[keyType] = slices
				}
			case int:
				strings, slices := jsonHelper(model)
				if 0 != len(strings) {
					jsonSlice = append(jsonSlice, strings)
				}
				if 0 != len(slices) {
					jsonSlice = append(jsonSlice, slices)
				}
			case nil:
				strings, slices := jsonHelper(model)
				if 0 != len(strings) {
					jsonMap[""] = strings
				}
				if 0 != len(slices) {
					jsonSlice = append(jsonSlice, slices)
				}
			}
		default:
			switch k.(type) {
			case string:
				jsonMap[k.(string)] = v
			case int:
				var tmpSlice []interface{}
				// This is to preserve array value order
				if len(jsonSlice) < k.(int)+1 {
					tmpSlice = make([]interface{}, k.(int)+1)
					for tmpk, tmpv := range jsonSlice {
						tmpSlice[tmpk] = tmpv
					}
					jsonSlice = tmpSlice
				}
				jsonSlice[k.(int)] = v
			}
		}
	}

	return jsonMap, jsonSlice
}

/*
Keys returns a sorted slice of model keys
*/
func (ma *Model) Keys() []string {
	keys := make([]string, 0)
	for k := range ma.data {
		switch k.(type) {
		case float64:
			keys = append(keys, strconv.Itoa(k.(int)))
		case int:
			keys = append(keys, strconv.Itoa(k.(int)))
		default:
			keys = append(keys, k.(string))
		}
	}
	sort.Strings(keys)
	return keys
}

/*
Len returns the number of items in this model
*/
func (ma *Model) Len() int {
	return len(ma.data)
}

/*
MakeStatic will set the static flag, making this model read-only
*/
func (ma *Model) MakeStatic() {
	ma.isStatic = true
}

var modelsDepth int

/*
Models recursively converts any Model instances in the internal
data storage map to nested map structures values and returns the result
*/
func (ma *Model) Models() *Model {

	if modelsDepth > 50 {
		panic("recursion too deep")
	}
	modelsDepth++

	retVal := NewModel()
	for k, v := range ma.Data() {
		// instanceof Model
		if model, ok := v.(Model); ok {
			retVal.Set(k, model.Models())

		} else {
			// Get data type
			rt := reflect.TypeOf(v)
			switch rt.Kind() {
			case reflect.Slice:
				fallthrough
			case reflect.Array:
				fallthrough
			case reflect.Map:
				// Rangeable data types
				// Converts all keys to strings... watch for collisions
				newV := NewModel()
				rangeV := reflect.ValueOf(v)
				for _, key := range rangeV.MapKeys() {
					newV.Set(key.String(), rangeV.MapIndex(key))
				}
				retVal.Set(k, newV.Models())

			// Store as-is
			default:
				retVal.Set(k, v)
			}
		}
	}

	return retVal
}

/*
Merge a model into this one
*/
func (ma *Model) Merge(model *Model) error {
	if ma.IsStatic() {
		return fmt.Errorf(ErrorCannotModifyStaticModel)
	}
	for k, v := range model.Data() {
		ma.data[k] = v
	}
	return nil
}

/*
Reduce iteratively reduces the data to a single value using a callback function
and returns that value
*/
func (ma *Model) Reduce(
	callback func(interface{}, interface{}) interface{},
	initial interface{}) interface{} {

	for _, v := range ma.data {
		initial = callback(initial, v)
	}
	return initial
}

/*
Reset deletes all values from the internal data store
*/
func (ma *Model) Reset() error {
	if ma.IsStatic() {
		return fmt.Errorf(ErrorCannotModifyStaticModel)
	}
	ma.data = make(map[interface{}]interface{})
	return nil
}

/*
Search the elements of this model for a given value and return the first
corresponding index if successful. If needle is a callback, each element
is passed in. If the element is not found, return false
*/
func (ma *Model) Search(needle interface{}) (interface{}, error) {
	for k, v := range ma.data {
		if v == needle {
			return k, nil
		}
	}
	return "", fmt.Errorf("No matching value found")
}

/*
Set a named value locally
*/
func (ma *Model) Set(idx interface{}, val interface{}) error {
	if ma.IsStatic() {
		return fmt.Errorf(ErrorCannotModifyStaticModel)
	}
	ma.data[idx] = val
	return nil
}

/*
SetData replaces the entire internal data storage array
*/
func (ma *Model) SetData(data map[interface{}]interface{}) error {
	if ma.IsStatic() {
		return fmt.Errorf(ErrorCannotModifyStaticModel)
	}
	ma.data = data
	return nil
}

/*
SetUID sets the model identifier property
*/
func (ma *Model) SetUID(identifier string) error {
	if ma.IsStatic() {
		return fmt.Errorf(ErrorCannotModifyStaticModel)
	}
	ma.identifier = identifier
	return nil
}

/*
SetUIDType sets the model identifierType property
*/
func (ma *Model) SetUIDType(identifierType int) error {
	if ma.IsStatic() {
		return fmt.Errorf(ErrorCannotModifyStaticModel)
	}
	switch identifierType {
	case ConstIdentifierTypeInt:
	case ConstIdentifierTypeString:
		ma.identifierType = identifierType
		return nil
	}
	return fmt.Errorf("Invalid identifier type '%d'", identifierType)
}

/*
SetView sets a view instance to be used with this model
*/
func (ma *Model) SetView(view *View) error {
	if ma.IsStatic() {
		return fmt.Errorf(ErrorCannotModifyStaticModel)
	}
	ma.view = view
	return nil
}

/*
String converts the data to a text representation, should generally be JSON by
default
*/
//func (ma *Model) String() string {
//	str, _ := ma.JSON()
//	return str
//}

/*
Transform applies a callback to all elements in this model and return the result
The current model is not modified
*/
func (ma *Model) Transform(callback func(interface{}) interface{}) map[interface{}]interface{} {
	retVal := make(map[interface{}]interface{})
	for k, v := range ma.data {
		retVal[k] = callback(v)
	}
	return retVal
}

/*
UID returns the model's unique identifier
Throw an exception if "name" has no meaning in your class.
*/
func (ma *Model) UID() interface{} {
	return ma.identifier
}

/*
Unique will remove any duplicates from this model
*/
func (ma *Model) Unique() error {
	if ma.IsStatic() {
		return fmt.Errorf(ErrorCannotModifyStaticModel)
	}
	keys := ma.Keys() // Iterate through the keys alpabetically
	uniques := make(map[interface{}]string)
	for _, k := range keys {
		uniques[ma.data[k]] = k
	}
	data := make(map[interface{}]interface{})
	for k, v := range uniques {
		data[v] = k
	}
	ma.data = data
	return nil
}

/*
View returns the View pointer for this model
*/
func (ma *Model) View() *View {
	return ma.view
}

/*
Walk will apply a user supplied function to every member of this model and STORE
the result
*/
func (ma *Model) Walk(callback func(interface{}) interface{}) error {
	if ma.IsStatic() {
		return fmt.Errorf(ErrorCannotModifyStaticModel)
	}
	ma.data = ma.Transform(callback)
	return nil
}
