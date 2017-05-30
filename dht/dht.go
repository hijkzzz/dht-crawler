package dht

import (
	"fmt"
	"net"
)

// BootstrapNodes 初始节点
var BootstrapNodes = []*KNode{
	NewKNode(-1, "router.bittorrent.com", 6881),
	NewKNode(-1, "dht.transmissionbt.com", 6881),
	NewKNode(-1, "router.utorrent.com", 6881)}

// TIDLength 交易号长度
var TIDLength = 2

// ReJoinDHTInterval 重加入间隔(秒)
var ReJoinDHTInterval = 3

// TokenLength Token 长度
var TokenLength = 2

// DHT BEP005 服务实现
type DHT struct {
	bindHost string
	bindPort int
	logger   chan map[string]string
	ktable   *KTable
	krpc     *KRPC
	udpConn  *net.UDPConn
}

// NewDHT 新建 DHT 服务器
func NewDHT(host string, port int, logger chan map[string]string) *DHT {
	return &DHT{bindHost: host, bindPort: port, logger: logger, ktable: NewKTable()}
}

// Run 运行 DHT 服务器
func (dht *DHT) Run() {
	// 监听 UDP 端口
	udpAddress, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", dht.bindHost, dht.bindPort))
	if err != nil {
		panic(err)
	}

	udpConn, err := net.ListenUDP("udp", udpAddress)
	if err != nil {
		panic(err)
	}

	defer udpConn.Close()

	dht.udpConn = udpConn
	dht.krpc = NewKRPC(udpConn)

	// 从 Bootstrap 交朋友
	go dht.joinDHT()

	// 接受 UDP 数据
	buff := make([]byte, 8192)
	for true {
		n, raddr, err := udpConn.ReadFromUDP(buff)
		if err != nil {
			panic(err)
		}

		// UDP 数据解码
		decode_msg, err := DecodeBencode(buff[:n])
		if err != nil {
			panic(err)
		}

		// 请求格式判断
	}
}

// joinDHT 从 Bootstrap 交朋友
func (dht *DHT) joinDHT() {
}

// updateKTable 自动刷新路由表
func (dht *DHT) updateKTable() {

}
