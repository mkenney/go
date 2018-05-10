package model

import (
	"encoding/json"
	"testing"
)

func TestTest(t *testing.T) {
	model, _ := New()
	json.Unmarshal([]byte(`{"z":100}`), model)
	t.Errorf("\n\n%+v\n\n", model)
}

func TestModel(t *testing.T) {
	var err error

	model, err := New()
	if nil != err {
		t.Errorf("expected nil, got error: %+v", err)
	}

	err = model.Set("key", "value")
	if nil != err {
		t.Errorf("expected nil, got error: %+v", err)
	}

	if !model.Has("key") {
		t.Errorf("expected true, got false")
	} else if model.Has("doesn't exist") {
		t.Errorf("expected false, got true")
	}

	val, err := model.Get("key")
	if nil != err {
		t.Errorf("expected nil, got error: %+v", err)
	}
	valstr, ok := val.(string)
	if !ok {
		t.Errorf("expected string but could not cast value: %+v", val)
	} else if "value" != valstr {
		t.Errorf("expected 'value', got '%s'", valstr)
	}
	_, err = model.Get("doesn't exist")
	if nil == err {
		t.Errorf("expected error, got nil")
	}

	val, err = model.Del("key")
	if nil != err {
		t.Errorf("expected nil, got error: %+v", err)
	} else if valstr, ok := val.(string); !ok {
		t.Errorf("expected string but could not cast value: %+v", val)
	} else if "value" != valstr {
		t.Errorf("expected 'value', got '%s'", valstr)
	}
	_, err = model.Del("doesn't exist")
	if nil == err {
		t.Errorf("expected error, got nil")
	}

	bytesIn := []byte(`{"z":100}`)
	err = json.Unmarshal(bytesIn, model)
	if nil != err {
		t.Errorf("expected nil, got error: %+v", err)
	}
	bytesOut, err := json.Marshal(model)
	if nil != err {
		t.Errorf("expected nil, got error: %+v", err)
	}
	if string(bytesIn) != string(bytesOut) {
		t.Errorf("expected %s, got error: %s", string(bytesIn), string(bytesOut))
	}

	//t.Errorf(`
	//	model: %+v
	//	key: %+v
	//	keystr: %+v
	//`, model, val, valstr,
	//)

}

func TestStaticModel(t *testing.T) {
	var err error
	model, _ := New()
	model.Set("a", 1)
	model.Set("b", 2)
	model.Set("c", 3)
	model.Lock()

	_, err = model.Del("a")
	if nil == err {
		t.Errorf("expected error, got nil")
	}

	if !model.IsStatic() {
		t.Errorf("expected true, got false")
	}

	err = model.Set("a", 100)
	if nil == err {
		t.Errorf("expected error, got nil")
	}

	err = model.Set("z", 100)
	if nil == err {
		t.Errorf("expected error, got nil")
	}

	err = json.Unmarshal([]byte(`{"z":100}`), model)
	if nil == err {
		t.Errorf("expected error, got nil")
	}

	//t.Errorf(`
	//	model: %#v
	//`, model,
	//)
}
