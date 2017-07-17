/*
Package model is a data-driven modeling abstraction
*/
package model

import (
	"encoding/json"
	"fmt"
	"sort"
)

/*
Collection is a collection of Model instances
*/
type Collection struct {

	/*
		Local data storage.
	*/
	data []*Model

	/*
		Optional, identifier of this data object.
		This is used as the index value in model collections
	*/
	identifier string

	/*
		The data type of the identifier value.
		Either ConstIdentifierTypeInt or ConstTypeString
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
NewCollection initializes and returns a pointer to a Model
*/
func NewCollection() (collection *Collection) {
	collection = new(Collection)
	collection.data = make([]*Model, 0)
	collection.view = new(View)
	return collection
}

/*
Data returns the internal data map
*/
func (cn *Collection) Data() []*Model {
	return cn.data
}

/*
Delete removes a locally stored value by index
*/
func (cn *Collection) Delete(idx int) (bool, error) {
	if cn.IsStatic() {
		return false, fmt.Errorf("static models cannot be modified")
	}
	if len(cn.data) <= idx {
		return false, fmt.Errorf(ErrorCollectionIndexDoesNotExist, idx)
	}

	cn.data = append(cn.data[:idx], cn.data[idx+1:]...)

	return true, nil
}

/*
IsEmpty checks to see if a value should be considered "empty"
True if the value does not exist, is nil or is an empty string
*/
func (cn *Collection) IsEmpty() bool {
	return 0 == len(cn.data)
}

/*
Filter filters elements of the data using a callback function and returns the
result
*/
func (cn *Collection) Filter(callback func(key int, model *Model) bool) []*Model {
	retVal := make([]*Model, 0)
	for k, v := range cn.data {
		if callback(k, v) {
			retVal = append(retVal, v)
		}
	}
	return retVal
}

/*
Get returns a locally stored value by index
*/
func (cn *Collection) Get(idx int) (*Model, error) {
	if !cn.Has(idx) {
		return nil, fmt.Errorf("the specified index '%d' does not exist", idx)
	}
	return cn.data[idx], nil
}

/*
GetUIDType returns a locally stored identifiers type
*/
func (cn *Collection) GetUIDType() int {
	return cn.identifierType
}

/*
Has checks to see if a value has been set
*/
func (cn *Collection) Has(idx int) bool {
	if len(cn.data) > idx {
		return true
	}
	return false
}

/*
IsStatic returns whether the static flag has been set or not
*/
func (cn *Collection) IsStatic() bool {
	return cn.isStatic
}

/*
JSON recursively converts the internal data storage array to a JSON string
Exposed as a public method to give access to the json_encode() parameters
*/
func (cn *Collection) JSON() (string, error) {
	data := make([]interface{}, 0)
	for _, v := range cn.data {
		data = append(data, v.Data())
	}
	str, err := json.Marshal(data)
	return string(str), err
}

/*
Keys returns a sorted slice of model keys
*/
func (cn *Collection) Keys() []int {
	keys := make([]int, len(cn.data))
	for k := range cn.data {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	return keys
}

/*
Len returns the number of items in this model
*/
func (cn *Collection) Len() int {
	return len(cn.data)
}

/*
MakeStatic will set the static flag, making this model read-only
*/
func (cn *Collection) MakeStatic() {
	cn.isStatic = true
}

/*
Models recursively converts any Model instances in the internal
data storage map to nested map structures values and returns the result
*/
func (cn *Collection) Models() []*Model {
	return cn.data
}

/*
Merge a model into this one
*/
func (cn *Collection) Merge(model *Model) error {
	cn.data = append(cn.data, model)
	return nil
}

/*
Push adds a Model pointer to the stack
*/
func (cn *Collection) Push(model *Model) {
	cn.data = append(cn.data, model)
}

/*
Pop removes the last Model pointer from the stack and returns it
*/
func (cn *Collection) Pop() (*Model, error) {
	if 0 == len(cn.data) {
		return nil, fmt.Errorf(ErrorCollectionIsEmpty)
	}

	model := cn.data[len(cn.data)-1]
	cn.data = cn.data[:len(cn.data)-1]

	return model, nil
}

/*
Reduce iteratively reduces the data to a single value using a callback function
and returns that value
*/
func (cn *Collection) Reduce(
	callback func(interface{}, interface{}) interface{},
	initial interface{}) interface{} {

	for _, v := range cn.data {
		initial = callback(initial, v)
	}
	return initial
}

/*
Reset deletes all values from the internal data store
*/
func (cn *Collection) Reset() error {
	if cn.IsStatic() {
		return fmt.Errorf(ErrorCannotModifyStaticCollection)
	}
	cn.data = make([]*Model, 0)
	return nil
}

/*
Search the elements of this model for a given value and return the first
corresponding index if successful. If needle is a callback, each element
is passed in. If the element is not found, return false
*/
func (cn *Collection) Search(needle interface{}) (int, error) {
	for k, v := range cn.data {
		if v == needle {
			return k, nil
		}
	}
	return -1, fmt.Errorf("No matching value found")
}

/*
Set a named value locally
*/
func (cn *Collection) Set(idx int, val *Model) error {
	if cn.IsStatic() {
		return fmt.Errorf(ErrorCannotModifyStaticCollection)
	}
	cn.data[idx] = val
	return nil
}

/*
SetData replaces the entire internal data storage array
*/
func (cn *Collection) SetData(data []*Model) error {
	if cn.IsStatic() {
		return fmt.Errorf(ErrorCannotModifyStaticCollection)
	}
	cn.data = data
	return nil
}

/*
SetUID sets the model identifier property
*/
func (cn *Collection) SetUID(identifier string) error {
	if cn.IsStatic() {
		return fmt.Errorf(ErrorCannotModifyStaticCollection)
	}
	cn.identifier = identifier
	return nil
}

/*
SetUIDType sets the model identifierType property
*/
func (cn *Collection) SetUIDType(identifierType int) error {
	if cn.IsStatic() {
		return fmt.Errorf(ErrorCannotModifyStaticCollection)
	}
	switch identifierType {
	case ConstIdentifierTypeInt:
	case ConstIdentifierTypeString:
		cn.identifierType = identifierType
		return nil
	}
	return fmt.Errorf("Invalid identifier type '%d'", identifierType)
}

/*
SetView sets a view instance to be used with this model
*/
func (cn *Collection) SetView(view *View) error {
	if cn.IsStatic() {
		return fmt.Errorf(ErrorCannotModifyStaticCollection)
	}
	cn.view = view
	return nil
}

/*
String converts the data to a text representation, should generally be JSON by
default
*/
func (cn *Collection) String() string {
	str, _ := cn.JSON()
	return str
}

/*
Transform applies a callback to all elements in this model and return the result
The current model is not modified
*/
func (cn *Collection) Transform(callback func(*Model) *Model) []*Model {
	retVal := make([]*Model, 0)
	for k, v := range cn.data {
		retVal[k] = callback(v)
	}
	return retVal
}

/*
UID returns the model's unique identifier
Throw an exception if "name" has no meaning in your class.
*/
func (cn *Collection) UID() interface{} {
	return cn.identifier
}

/*
Unique will remove any duplicates from this model
*/
func (cn *Collection) Unique() error {
	if cn.IsStatic() {
		return fmt.Errorf(ErrorCannotModifyStaticCollection)
	}
	keys := cn.Keys() // Iterate through the keys alpabetically
	uniques := make(map[interface{}]int)
	for _, k := range keys {
		uniques[cn.data[k].identifier] = k
	}
	data := make([]*Model, 0)
	for _, k := range uniques {
		data = append(data, cn.data[k])
	}
	cn.data = data
	return nil
}

/*
View returns the View pointer for this model
*/
func (cn *Collection) View() *View {
	return cn.view
}

/*
Walk will apply a user supplied function to every member of this model and STORE
the result
*/
func (cn *Collection) Walk(callback func(*Model) *Model) error {
	if cn.IsStatic() {
		return fmt.Errorf(ErrorCannotModifyStaticCollection)
	}
	cn.data = cn.Transform(callback)
	return nil
}
