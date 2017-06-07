package main

import "dht-crawler/dht"

func main() {
	dht := dht.NewDHT("127.0.0.1", 6881)
	dht.Run()
}
