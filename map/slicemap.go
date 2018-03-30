package main

import (
	"fmt"
	"time"
)

func testSliceMap() {
	testMap := slicemap{}
	var start time.Time
	var end time.Time

	start = time.Now()
	for a := 0; a < length; a++ {
		testMap.set(a, a*100)
	}
	end = time.Now()
	fmt.Println("slicemap set time: ", end.Sub(start))

	start = time.Now()
	for a := 0; a < length; a++ {
		val, found := testMap.get(a)
		if a*100 != val {
			panic(fmt.Sprintf("val %d = %v (%v)", a, val, found))
		}
	}
	end = time.Now()
	fmt.Println("slicemap get time: ", end.Sub(start))

	start = time.Now()
	for a := 0; a < length; a++ {
		testMap.del(a)
	}
	end = time.Now()
	fmt.Println("slicemap del time: ", end.Sub(start))
}

type slicemap struct {
	keys   []interface{}
	values []interface{}
}

func (m *slicemap) del(k interface{}) bool {
	for mk := range m.keys {
		if m.keys[mk] == k {
			mapMux.Lock()
			m.keys = append(m.keys[:mk], m.keys[mk+1:]...)
			m.values = append(m.values[:mk], m.values[mk+1:]...)
			mapMux.Unlock()
			return true
		}
	}
	return false
}

func (m *slicemap) get(k interface{}) (interface{}, bool) {

	for mk := range m.keys {
		if mk == k {
			mapMux.Lock()
			defer mapMux.Unlock()
			switch m.values[mk].(type) {
			case int:
				return m.values[mk].(int), true
			case string:
				return m.values[mk].(string), true
			}
		}
	}
	return nil, false
}

func (m *slicemap) set(k interface{}, v interface{}) bool {
	for mk := range m.keys {
		if mk == k {
			mapMux.Lock()
			m.values[mk] = v
			mapMux.Unlock()
			return true
		}
	}
	mapMux.Lock()
	m.keys = append(m.keys, k)
	m.values = append(m.values, v)
	mapMux.Unlock()
	return false
}
