package dht

// kTable 路由表
type kTable struct {
	nodes []*kNode
}

// newKTable 新建路由表
func newKTable() *kTable {
	return &kTable{make([]*kNode, 0, 8192)}
}

// size KTable 的大小
func (ktable *kTable) size() int {
	return len(ktable.nodes)
}
