package dht

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// BootstrapNodes 初始节点
var BootstrapNodes = []*kNode{
	newKNode("", "router.bittorrent.com", 6881),
	newKNode("", "dht.transmissionbt.com", 6881),
	newKNode("", "router.utorrent.com", 6881)}

// DHT BEP005 服务实现
type DHT struct {
	bindHost  string                 // 监听地址
	bindPort  int                    // 监听端口
	logger    chan map[string]string // 传输 info_hash
	ktable    *kTable                // 路由表
	krpc      *kRPC                  // KRPC 协议
	udpConn   *net.UDPConn           // UDP 连接
	waitGroup *sync.WaitGroup        // 等待子线程
}

// NewDHT 新建 DHT 服务器, seed 作为种子生成 ID
func NewDHT(host string, port int) *DHT {
	// 监听 UDP 端口
	udpAddress, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}

	udpConn, err := net.ListenUDP("udp", udpAddress)
	if err != nil {
		panic(err)
	}

	// 新建 DHT 服务器
	dht := &DHT{bindHost: host,
		bindPort:  port,
		logger:    make(chan map[string]string, 8192),
		ktable:    newKTable(),
		udpConn:   udpConn,
		waitGroup: new(sync.WaitGroup)}

	// krpc 协议初始化
	dht.krpc = newKRPC(dht)

	return dht
}

// Run 运行 DHT 服务器
func (dht *DHT) Run() {
	dht.waitGroup.Add(3)
	defer dht.udpConn.Close()

	// 线程1, 更新路由表
	go dht.updateKtable()

	// 线程2, 处理 UDP 报文
	go dht.receiveMessages()

	// 线程3，处理 info_hash
	go dht.processInfoHash()

	dht.waitGroup.Wait()
}

// receiveMessages 处理 UDP 报文
func (dht *DHT) receiveMessages() {
	defer dht.waitGroup.Done()

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

		// 报文类型判断
		y, ok := message["y"].(string)
		if !ok {
			fmt.Println("KRPC request missing y field")
			dht.krpc.sendError(message, 203, raddr)
			return
		}

		if y == "q" { //处理请求报文
			q, ok := message["q"].(string)
			if !ok {
				fmt.Println("KRPC request missing q field")
				dht.krpc.sendError(message, 203, raddr)
				return
			}

			switch q {
			case "ping":
				dht.krpc.requestPing(message, raddr)
			case "find_node":
				dht.krpc.requestFindNode(message, raddr)
			case "get_peers":
				dht.krpc.requestGetPeers(message, raddr)
			case "announce_peer":
				// 收集 announce_peer 的 info_hash
				dht.krpc.requestAnnouncePeer(message, raddr)
			default:
				dht.krpc.sendError(message, 203, raddr)
				fmt.Println("KRPC not support 'q' " + q)
			}
		} else if y == "r" { //处理响应报文
			dht.krpc.responseFindNode(message, raddr)
		} else if y == "e" { //处理错误报文
			fmt.Println("KRPC value of 'y' is 'e' ")
			fmt.Println(message)
		} else {
			dht.krpc.sendError(message, 204, raddr)
			fmt.Println("KRPC value of 'y' error " + y)
		}
	}
}

// findNewNodes 更新路由表
func (dht *DHT) updateKtable() {
	defer dht.waitGroup.Done()

	for true {
		len := dht.ktable.size()
		if len == 0 {
			for _, node := range BootstrapNodes {
				dht.krpc.sendFindNode(getNeigborID(node.nid, dht.krpc.nid, 0), node.getUDPAddr())
			}

		} else {
			for len > 0 {
				len--
				node := dht.ktable.pop()
				dht.krpc.sendFindNode(getNeigborID(node.nid, dht.krpc.nid, 10), node.getUDPAddr())
			}
		}
		time.Sleep(1 * time.Second)
	}
}

// processInfoHash 处理 info_hash
func (dht *DHT) processInfoHash() {
	defer dht.waitGroup.Done()

	for true {
		message := <-dht.logger
		fmt.Println(message["info_hash"])
	}
}
