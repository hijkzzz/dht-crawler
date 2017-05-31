package dht

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

func TestKTable(t *testing.T) {
	var w sync.WaitGroup
	w.Add(2)
	ktable := newKTable()
	// 线程1
	go func() {
		i := 0
		for true {
			fmt.Println("PUSH " + strconv.Itoa(i))
			ktable.push(&kNode{nid: strconv.Itoa(i)})
			i++
		}
		w.Done()
	}()

	// 线程 2
	go func() {
		for true {
			fmt.Println("POP " + ktable.pop().nid)
		}
		w.Done()
	}()

	w.Wait()
}
