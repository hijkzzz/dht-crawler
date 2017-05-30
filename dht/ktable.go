package dht

import (
	"container/list"
)

// KTable 路由表
type KTable struct {
	*list.List
}

// NewKTable 新建路由表
func NewKTable() *KTable {
	return &KTable{}
}

// Size KTable 的大小
func (ktable *KTable) Size() int {
	return ktable.List.Len()
}
