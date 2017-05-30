package dht

import (
	"fmt"
)

// KNode DHT 网络节点
type kNode struct {
	nid  int
	host string
	port int
}

// NewKNode 新建 DHT 网络节点
func newKNode(nid int, ip string, port int) *kNode {
	return &kNode{nid, ip, port}
}

func (knode *kNode) getHostPort() string {
	return fmt.Sprintf("%s:%d", knode.host, knode.port)
}
