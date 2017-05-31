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
