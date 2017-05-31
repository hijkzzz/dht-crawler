package dht

import (
	"fmt"
	"net"
)

// KRPC krpc 协议
type kRPC struct {
	udpConn *net.UDPConn
	logger  chan<- map[string]string
}

// NewKRPC 新建 krpc 协议
func newKRPC(conn *net.UDPConn, logger chan<- map[string]string) *kRPC {
	return &kRPC{udpConn: conn, logger: logger}
}

// sendKRPC 发送 KRPC 请求
func (krpc *kRPC) sendKRPC(msg map[string]interface{}, host string, port int) {
	udpAddress, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Println(err)
		return
	}

	message, err := encodeBencode(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = krpc.udpConn.WriteToUDP(message, udpAddress)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// sendFindNode find_node 请求
func (krpc *kRPC) sendFindNode(target *kNode) {

}

// responsePing ping 响应
func (krpc *kRPC) responsePing(msg map[string]interface{}, address *net.UDPAddr) {

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
