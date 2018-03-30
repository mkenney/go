package main

import (
	"fmt"
	"time"
)

func testModMap() {
	testMap := modmap{}
	var start time.Time
	var end time.Time

	start = time.Now()
	for a := 0; a < length; a++ {
		testMap.set(a, a*100)
	}
	end = time.Now()
	fmt.Println("modmap set time: ", end.Sub(start))

	start = time.Now()
	for a := 0; a < length; a++ {
		val, found := testMap.get(a)
		if a*100 != val {
			panic(fmt.Sprintf("val %d = %v (%v)", a, val, found))
		}
	}
	end = time.Now()
	fmt.Println("modmap get time: ", end.Sub(start))

	start = time.Now()
	for a := 0; a < length; a++ {
		testMap.del(a)
	}
	end = time.Now()
	fmt.Println("modmap del time: ", end.Sub(start))
}

type modmap struct {
	keys   [10][]interface{}
	values []interface{}
}

func (m *modmap) del(k interface{}) bool {
	bucket := k.(int) % 10
	for mk := range m.keys[bucket] {
		if m.keys[bucket][mk] == k {
			mapMux.Lock()
			m.keys[bucket] = append(m.keys[bucket][:mk], m.keys[bucket][mk+1:]...)
			m.values = append(m.values[:mk], m.values[mk+1:]...)
			mapMux.Unlock()
			return true
		}
	}
	return false
}

func (m *modmap) get(k interface{}) (interface{}, bool) {
	bucket := k.(int) % 10
	for _, mk := range m.keys[bucket] {
		if mk == k {
			mapMux.Lock()
			defer mapMux.Unlock()
			switch mk.(type) {
			case int:
				return m.values[mk.(int)].(int), true
			case string:
				return m.values[mk.(int)].(string), true
			}
		}
	}
	return nil, false
}

func (m *modmap) set(k interface{}, v interface{}) bool {
	bucket := k.(int) % 10
	for mk := range m.keys[bucket] {
		if mk == k {
			mapMux.Lock()
			m.values[mk] = v
			mapMux.Unlock()
			return true
		}
	}
	mapMux.Lock()
	m.keys[bucket] = append(m.keys[bucket], k)
	m.values = append(m.values, v)
	mapMux.Unlock()
	return false
}
