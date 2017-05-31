package dht

import (
	"container/list"
	"sync"
)

// kTable 路由表, 线程安全
type kTable struct {
	nodes *list.List
	mutex *sync.Mutex
	cond  *sync.Cond
}

// newKTable 新建路由表
func newKTable() *kTable {
	ktable := &kTable{nodes: list.New(), mutex: new(sync.Mutex)}
	ktable.cond = sync.NewCond(ktable.mutex)
	return ktable
}

// size KTable 的大小
func (ktable *kTable) size() int {
	ktable.mutex.Lock()
	defer ktable.mutex.Unlock()

	return ktable.nodes.Len()
}

func (ktable *kTable) push(node *kNode) {
	ktable.mutex.Lock()
	defer ktable.mutex.Unlock()

	ktable.nodes.PushBack(node)
	ktable.cond.Broadcast()
}

func (ktable *kTable) pop() *kNode {
	ktable.mutex.Lock()
	defer ktable.mutex.Unlock()

	front := ktable.nodes.Front()
	for front == nil {
		ktable.cond.Wait()
		front = ktable.nodes.Front()
	}

	ktable.nodes.Remove(front)
	return front.Value.(*kNode)
}
