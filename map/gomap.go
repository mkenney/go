package main

import (
	"fmt"
	"time"
)

func testGoMap() {
	testMap := gomap{}
	var start time.Time
	var end time.Time

	start = time.Now()
	for a := 0; a < length; a++ {
		testMap.set(a, a*100)
	}
	end = time.Now()
	fmt.Println("gomap set time: ", end.Sub(start))

	start = time.Now()
	for a := 0; a < length; a++ {
		val, found := testMap.get(a)
		if a*100 != val {
			panic(fmt.Sprintf("val %d = %v (%v)", a, val, found))
		}
	}
	end = time.Now()
	fmt.Println("gomap get time: ", end.Sub(start))

	start = time.Now()
	for a := 0; a < length; a++ {
		testMap.del(a)
	}
	end = time.Now()
	fmt.Println("gomap del time: ", end.Sub(start))
}

var size = 1

type gomap struct {
	keys   [][]interface{}
	values []interface{}
}

func (m *gomap) del(k interface{}) bool {
	bucket, err := m.getBucket(k)
	if nil != err {
		return false
	}
	for bk, mk := range m.keys[bucket] {
		if mk == k {
			mapMux.Lock()
			m.keys[bucket] = append(m.keys[bucket][:bk], m.keys[bucket][bk+1:]...)
			m.values = append(m.values[:bk], m.values[bk+1:]...)
			mapMux.Unlock()
			return true
		}
	}
	if 0 == len(m.keys[bucket]) {
		m.keys = append(m.keys[:bucket], m.keys[bucket+1:]...)
	}
	return false
}

func (m *gomap) get(k interface{}) (interface{}, bool) {
	bucket, err := m.getBucket(k)
	if nil != err {
		return nil, false
	}
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

func (m *gomap) set(k interface{}, v interface{}) bool {
	var found bool
	fmt.Println("got here")
	return false
	bucket, err := m.getBucket(k)
	if nil != err {
		for bucket = range m.keys {
			if len(m.keys[bucket]) < size {
				found = true
				break
			}
		}
		if !found {
			bucket = m.alloc()
		}
	}

	for mk := range m.keys[bucket] {
		if m.keys[bucket][mk] == k {
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

func (m *gomap) alloc() int {
	m.keys = append(m.keys, make([]interface{}, 0))
	return len(m.keys) - 1
}

type circutBreaker struct {
	stop bool
}

func (m *gomap) getBucket(query interface{}) (int, error) {
	breaker := &circutBreaker{}
	done := make(chan int, len(m.keys))
	result := make(chan int)

	for bucketID, bucket := range m.keys {
		go func(
			bucketID int,
			bucket []interface{},
			query interface{},
			result chan int,
			done chan int,
			breaker *circutBreaker,
		) {
			for key := range bucket {
				if breaker.stop {
					fmt.Println("breaker.stop")
					done <- 1
					return
				}
				if query == bucket[key] {
					fmt.Println("result <- bucketID")
					result <- bucketID
					done <- 1
					return
				}
			}
		}(bucketID, bucket, query, result, done, breaker)
	}

	var bucket int
	var doneCount int
	numBuckets := len(m.keys)
	if numBuckets > 0 {
		for {
			select {
			case bucket = <-result:
				return bucket, nil
			case <-done:
				doneCount++
			}
			if doneCount == numBuckets {
				return 0, fmt.Errorf("bucket not found")
			}
		}
	}
	return 0, fmt.Errorf("bucket not found")
}
