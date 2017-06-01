package dht

import (
	"crypto/sha1"
	"fmt"
	"io"
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
