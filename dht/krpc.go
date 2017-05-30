package dht

import (
	"fmt"
	"net"
)

// KRPC krpc 协议
type KRPC struct {
	udpConn *net.UDPConn
}

// NewKRPC 新建 krpc 协议
func NewKRPC(conn *net.UDPConn) *KRPC {
	return &KRPC{udpConn: conn}
}

// sendKRPC 发送 KRPC 请求
func (krpc *KRPC) sendKRPC(msg map[string]interface{}, host string, port int) {
	udpAddress, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}

	encodeMsg, err := EncodeBencode(msg)
	if err != nil {
		panic(err)
	}

	_, err = krpc.udpConn.WriteToUDP(encodeMsg, udpAddress)
	if err != nil {
		panic(err)
	}
}

// SendFindNode find_node 请求
func (krpc *KRPC) SendFindNode() {

}

// ResponsePing ping 响应
func (krpc *KRPC) ResponsePing(msg map[string]interface{}, address *net.UDPAddr) {

}

// ResponseFindNode find_node 响应
func (krpc *KRPC) ResponseFindNode(msg map[string]interface{}, address *net.UDPAddr) {

}

// ResponseAnnouncePeer announce_peer 响应
func (krpc *KRPC) ResponseAnnouncePeer(msg map[string]interface{}, address *net.UDPAddr) {

}

// ResponseGetPeers get_peers 响应
func (krpc *KRPC) ResponseGetPeers(msg map[string]interface{}, address *net.UDPAddr) {

}
