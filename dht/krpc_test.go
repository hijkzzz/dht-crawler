package dht

import (
	"crypto/sha1"
	"fmt"
	"io"
	"net"
	"testing"
)

func TestEntropy(t *testing.T) {

	for i := 0; i < 100; i++ {
		rnd := entropy(20)
		t := sha1.New()
		io.WriteString(t, rnd)
		fmt.Println(fmt.Sprintf("%x", t.Sum(nil)))
	}
}

func TestDecodeCompactNodeInfo(t *testing.T) {
	// 生成 ID
	rnd := entropy(20)
	sha := sha1.New()
	io.WriteString(sha, rnd)
	id := fmt.Sprintf("%x", sha.Sum(nil))[:20]
	fmt.Println(id)

	// 主机地址 192.168.1.1 低位 -> 高位
	host := []byte{0x01, 0x01, 0xa8, 0xc0}

	// 端口 255
	port := []byte{0x00, 0xFF}

	compactInfo := id + string(host) + string(port)

	kNodes := decodeCompactNodeInfo(compactInfo + compactInfo)
	fmt.Println(len(kNodes))

	for _, kNode := range kNodes {
		fmt.Println(kNode.getHostPort())
	}
}

func Test_RequestGetPeers(t *testing.T) {
	//返回net.UDPConn 用于发送测试数据
	udpConn := udpConnOpen("127.0.0.1", 45678)

	msg := map[string]interface{}{
		"t": "aa",
		"y": "q",
		"q": "get_peers",
		"a": map[string]interface{}{
			"id":        "mnopqrstuvwxyz12345",
			"info_hash": "mnopqrstuvwxyz12345",
		},
	}
	//发送测试数据，并打印接受数据
	sendUDPMessage(udpConn, msg, "127.0.0.1", 34567)
}

func Test_ResponseFindNode(t *testing.T) {

}

func Test_RequestPing(t *testing.T) {

}

func Test_RequestFindNode(t *testing.T) {

}

func Test_RequestAnnouncePeers(t *testing.T) {

}

func Test_SendFindNode(t *testing.T) {

}

func Test_SendError(t *testing.T) {

}

func udpConnOpen(ip string, port int) *net.UDPConn {
	udpAddress, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		panic(err)
	}

	udpConn, err := net.ListenUDP("udp", udpAddress)
	if err != nil {
		panic(err)
	}

	return udpConn
}

func sendUDPMessage(udpConn *net.UDPConn, msg map[string]interface{}, ip string, port int) {
	defer udpConn.Close()

	message, err := encodeBencode(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		panic(err)
	}
	_, err = udpConn.WriteToUDP(message, addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	buff := make([]byte, 65536)
	n, _, err := udpConn.ReadFromUDP(buff)
	if err != nil {
		fmt.Println(err)
		return
	}

	reMsg, err := decodeBencode(buff[:n])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(reMsg)
}
