package dht

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"
)

// TIDLength 交易号长
var TIDLength = 2

// TokenLength Token 长度
var TokenLength = 2

// entropy 随机生成 len 长度的字符串
func entropy(len int) string {
	rand.Seed(time.Now().UnixNano())
	var buff bytes.Buffer

	for i := 0; i < len; i++ {
		rnd := byte(rand.Intn(255))
		buff.WriteByte(rnd)
	}
	return buff.String()
}

// decodeCompactNodeInfo Compact Node Info 解码
func decodeCompactNodeInfo(nodes []interface{}) []*kNode {
	return nil
}

// KRPC krpc 协议
type kRPC struct {
	nid     string                   // 本节点ID
	udpConn *net.UDPConn             // UDP 描述符
	ktable  *kTable                  // DHT 路由表
	logger  chan<- map[string]string // info_hash 传输
}

// NewKRPC 新建 krpc 协议
func newKRPC(dht *DHT) *kRPC {
	// 生成 nid
	rnd := entropy(20)
	t := sha1.New()
	io.WriteString(t, rnd)
	nid := fmt.Sprintf("%x", t.Sum(nil))
	fmt.Println("[NID] " + nid)

	return &kRPC{nid: nid, udpConn: dht.udpConn, ktable: dht.ktable, logger: dht.logger}
}

// sendKRPC 发送 KRPC 请求
func (krpc *kRPC) sendKRPC(msg map[string]interface{}, address *net.UDPAddr) {
	message, err := encodeBencode(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = krpc.udpConn.WriteToUDP(message, address)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// sendFindNode 发送 find_node 请求
func (krpc *kRPC) sendFindNode(target string, address *net.UDPAddr) {
	tid := entropy(TIDLength)

	msg := map[string]interface{}{
		"t": tid,
		"y": "q",
		"q": "find_node",
		"a": map[string]interface{}{
			"id":     krpc.nid,
			"target": target,
		},
	}

	krpc.sendKRPC(msg, address)
}

// responseFindNode 处理 find_node 响应
func (krpc *kRPC) responseFindNode(msg map[string]interface{}, address *net.UDPAddr) {
	// 处理消息
	tid, ok := msg["t"].(string)
	if !ok {
		fmt.Println("Message 'find_node' missing 'tid'")
		return
	}

	r, ok := msg["r"].(map[string]interface{})
	if !ok {
		fmt.Println("Message 'find_node' missing 'r'")
		return
	}

	nodes, ok := r["nodes"].([]interface{})
	if !ok {
		fmt.Println("Message 'find_node' missing 'nodes'")
		return
	}

	compactNodeInfos := decodeCompactNodeInfo(nodes)

	// 路由表更新
	for _, node := range compactNodeInfos {
		krpc.ktable.push(node)
	}
}

// requestPing 处理 ping 请求
func (krpc *kRPC) requestPing(msg map[string]interface{}, address *net.UDPAddr) {
	krpc.sendKRPC(msg, address)
}

// requestFindNode 处理 find_node 请求
func (krpc *kRPC) requestFindNode(msg map[string]interface{}, address *net.UDPAddr) {

}

// responseAnnouncePeer 处理 announce_peer 请求
func (krpc *kRPC) requestAnnouncePeer(msg map[string]interface{}, address *net.UDPAddr) {

}

// responseGetPeers 处理 get_peers 请求
func (krpc *kRPC) requestGetPeers(msg map[string]interface{}, address *net.UDPAddr) {

}

// responseError 发送 error 信息
func (krpc *kRPC) sendError(msg map[string]interface{}, address *net.UDPAddr) {

}
