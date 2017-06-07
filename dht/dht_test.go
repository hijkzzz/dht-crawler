package dht

import (
	"testing"
)

func Test_run(t *testing.T) {
	dht := NewDHT("127.0.0.1", 34567)
	dht.Run()
}
