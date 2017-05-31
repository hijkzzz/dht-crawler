package dht

import (
	"fmt"
)

// getNeigborID 生成一个邻居节点
func getNeigborID(target string, nid string, end int) string {
	return target[:end] + nid[end:]
}

// KNode DHT 网络节点
type kNode struct {
	nid  string
	host string
	port int
}

// NewKNode 新建 DHT 网络节点
func newKNode(nid string, ip string, port int) *kNode {
	return &kNode{nid, ip, port}
}

func (knode *kNode) getHostPort() string {
	return fmt.Sprintf("%s:%d", knode.host, knode.port)
}
