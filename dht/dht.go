package dht

import (
	"fmt"
	"net"
	"time"
)

// BootstrapNodes 初始节点
var BootstrapNodes = []*kNode{
	newKNode(-1, "router.bittorrent.com", 6881),
	newKNode(-1, "dht.transmissionbt.com", 6881),
	newKNode(-1, "router.utorrent.com", 6881)}

// TIDLength 交易号长度
var TIDLength = 2

// ReJoinDHTInterval 重加入间隔(秒)
var ReJoinDHTInterval = 3

// TokenLength Token 长度
var TokenLength = 2

// DHT BEP005 服务实现
type DHT struct {
	bindHost string                 // 监听地址
	bindPort int                    // 监听端口
	logger   chan map[string]string // info_hash 导出
	ktable   *kTable                // 路由表
	krpc     *kRPC                  // KRPC 协议
	udpConn  *net.UDPConn           // UDP 连接
}

// NewDHT 新建 DHT 服务器
func NewDHT(host string, port int, logger chan map[string]string) *DHT {
	// 监听 UDP 端口
	udpAddress, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}

	udpConn, err := net.ListenUDP("udp", udpAddress)
	if err != nil {
		panic(err)
	}

	return &DHT{bindHost: host, bindPort: port, logger: logger,
		ktable: newKTable(), krpc: newKRPC(udpConn), udpConn: udpConn}
}

// Run 运行 DHT 服务器
func (dht *DHT) Run() {
	// 线程1, 更新路由表
	go dht.findNewNodes()

	// 线程2, 处理 UDP 报文
	go dht.receiveMessages()
}

// receiveMessages 处理 UDP 报文
func (dht *DHT) receiveMessages() {
	defer dht.udpConn.Close()

	buff := make([]byte, 65536)
	for true {
		// 读取 UDP 数据
		n, raddr, err := dht.udpConn.ReadFromUDP(buff)
		if err != nil {
			fmt.Println(err)
			return
		}

		// UDP 数据解码
		message, err := decodeBencode(buff[:n])
		if err != nil {
			fmt.Println(err)
			return
		}

		// 请求格式判断
		method, ok := message["q"].(string)
		if !ok {
			fmt.Println("KRPC request missing q field")
			return
		}

		switch method {
		case "ping":
			dht.krpc.responsePing(message, raddr)
		case "find_node":
			dht.krpc.responseFindNode(message, raddr)
		case "get_peers":
			dht.krpc.responseGetPeers(message, raddr)
		case "announce_peer":
			dht.krpc.responseAnnouncePeer(message, raddr)
		default:
			dht.krpc.responseError(message, raddr)
			fmt.Println("KRPC not support q " + method)
		}
	}
}

// findNewNodes 更新路由表
func (dht *DHT) findNewNodes() {
	for true {
		if dht.ktable.size() == 0 {
			for _, node := range BootstrapNodes {
				dht.krpc.sendFindNode(node)
			}
		} else {
			for _, node := range dht.ktable.nodes {
				dht.krpc.sendFindNode(node)
			}

			time.Sleep(1 * time.Second)
		}
	}
}
