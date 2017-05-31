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

// sendFindNode find_node 请求
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

// responsePing ping 响应
func (krpc *kRPC) responsePing(msg map[string]interface{}, address *net.UDPAddr) {
	krpc.sendKRPC(msg, address)
}

// responseFindNode find_node 响应
func (krpc *kRPC) responseFindNode(msg map[string]interface{}, address *net.UDPAddr) {

}

// responseAnnouncePeer announce_peer 响应
func (krpc *kRPC) responseAnnouncePeer(msg map[string]interface{}, address *net.UDPAddr) {

}

// responseGetPeers get_peers 响应
func (krpc *kRPC) responseGetPeers(msg map[string]interface{}, address *net.UDPAddr) {

}

// responseError error 响应
func (krpc *kRPC) responseError(msg map[string]interface{}, address *net.UDPAddr) {

}
