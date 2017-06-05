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
var TokenLength = 8

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

// inet_ntoa 网络字节序 IP 转换为 ASCII
func inet_ntoa(ipnr []byte) net.IP {
	if len(ipnr) != 4 {
		panic("inet_ntoa lenghth of 'ipnr' should be 4")
	}

	// 192.168.1.1 低位 -> 高位
	return net.IPv4(ipnr[3], ipnr[2], ipnr[1], ipnr[0])
}

// decodeCompactNodeInfo Compact Node Info 解码
func decodeCompactNodeInfo(nodes interface{}) []*kNode {
	kNodes := make([]*kNode, 0, 8)

	temp, ok := nodes.(string)
	if !ok {
		fmt.Println("CompactNodeInfos is not string")
	}

	compactNodeInfos := []byte(temp)
	if len(compactNodeInfos)%26 != 0 {
		return kNodes
	}

	for i := 0; i < len(compactNodeInfos); i += 26 {
		// port 字节序转换
		var port int
		port += int(compactNodeInfos[i+25])
		port += int(compactNodeInfos[i+24]) << 8

		kNodes = append(kNodes, newKNode(
			string(compactNodeInfos[i:i+20]),
			inet_ntoa(compactNodeInfos[i+20:i+24]).String(),
			port,
		))
	}
	return kNodes
}

// KRPC krpc 协议
type kRPC struct {
	nid     string                   // sha1 生成的 node ID
	udpConn *net.UDPConn             // UDP 描述符
	ktable  *kTable                  // DHT 路由表
	logger  chan<- map[string]string // info_hash 传输
}

// NewKRPC 新建 krpc 协议, seed 作为种子生成 ID
func newKRPC(dht *DHT, seed string) *kRPC {
	// 生成 nid
	t := sha1.New()
	io.WriteString(t, seed)
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
		if len(node.nid) != 20 || node.port < 1 || node.port > 65535 || node.nid == krpc.nid {
			continue
		}
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
	a, ok := msg["a"].(map[string]interface{})
	if !ok {
		fmt.Println("Message 'announce_peers' missing 'a'")
		krpc.sendError(msg, 203, address)
		return
	}

	infoHash, ok := a["info_hash"].(string)
	if !ok {
		fmt.Println("Message 'get_peers' missing 'info_hash'")
		krpc.sendError(msg, 203, address)
		return
	}

	if len(infoHash) != 20 {
		fmt.Println("info_hash of message 'get_peers' is error")
		krpc.sendError(msg, 203, address)
		return
	}

	hashMsg := map[string]string{
		"info_hash": infoHash,
	}

	krpc.logger <- hashMsg

	fNid, ok := a["id"].(string)
	if !ok {
		fmt.Println("Message 'announce_peers' missing 'nid'")
		krpc.sendError(msg, 203, address)
		return
	}

	token, ok := a["token"].(string)
	if !ok {
		fmt.Println("Message 'announce_peers' missing 'token'")
		krpc.sendError(msg, 203, address)
		return
	}

	if fNid[0:8] != token {
		fmt.Println("'token' of message 'announce_peers' is error")
		krpc.sendError(msg, 203, address)
		return
	}

	//tid错误，只打印信息，不调用sendError
	tid, ok := msg["t"].(string)
	if !ok {
		fmt.Println("Message 'announce_peers' missing 'tid'")
		return
	}

	reMsg := map[string]interface{}{
		"t": tid,
		"y": "r",
		"r": map[string]interface{}{
			"id": krpc.nid,
		},
	}

	krpc.sendKRPC(reMsg, address)

}

// responseGetPeers 处理 get_peers 请求
func (krpc *kRPC) requestGetPeers(msg map[string]interface{}, address *net.UDPAddr) {
	a, ok := msg["a"].(map[string]interface{})
	if !ok {
		fmt.Println("Message 'get_peers' missing 'a'")
		krpc.sendError(msg, 203, address)
		return
	}

	//保存info_hash
	infoHash, ok := a["info_hash"].(string)
	if !ok {
		fmt.Println("Message 'get_peers' missing 'info_hash'")
		krpc.sendError(msg, 203, address)
		return
	}

	if len(infoHash) != 20 {
		fmt.Println("info_hash of message 'get_peers' is error")
		krpc.sendError(msg, 203, address)
		return
	}

	hashMsg := map[string]string{
		"info_hash": infoHash,
	}

	krpc.logger <- hashMsg

	//tid错误，只打印信息，不调用sendError
	tid, ok := msg["t"].(string)
	if !ok {
		fmt.Println("Message 'get_peers' missing 'tid'")
		return
	}

	fNid, ok := a["id"].(string)
	if !ok {
		fmt.Println("Message 'get_peers' missing 'nid'")
		krpc.sendError(msg, 203, address)
		return
	}

	token := fNid[0:8]

	nodes := ""

	reMsg := map[string]interface{}{
		"t": tid,
		"y": "r",
		"r": map[string]interface{}{
			"id":    krpc.nid,
			"token": token,
			"nodes": nodes,
		},
	}

	krpc.sendKRPC(reMsg, address)

}

// responseError 发送 error 信息, msg:krpc消息,errNum:错误码(203 协议错误)
func (krpc *kRPC) sendError(msg map[string]interface{}, errNum int, address *net.UDPAddr) {
	tid, ok := msg["t"].(string)
	if !ok {
		fmt.Println("Message 'announce_peers' missing 'tid'")
		return
	}

	errMsg := map[string]interface{}{
		"t": tid,
		"y": "e",
		"e": []interface{}{
			errNum, "Protocol errors, such as non-standard packages, invalid parameters, or wrong toke",
		},
	}

	krpc.sendKRPC(errMsg, address)

}
