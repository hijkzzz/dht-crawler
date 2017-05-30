package dht

// KNode DHT 网络节点
type KNode struct {
	nid  int
	host string
	port int
}

// NewKNode 新建 DHT 网络节点
func NewKNode(nid int, ip string, port int) *KNode {
	return &KNode{nid, ip, port}
}
