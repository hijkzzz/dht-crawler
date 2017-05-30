package dht

import "errors"

// encodeBencode bencode 编码
func encodeBencode(msg map[string]interface{}) ([]byte, error) {
	return nil, errors.New("this is a new error")
}

// decodeBencode bencode 解码
func decodeBencode(msg []byte) (map[string]interface{}, error) {
	return make(map[string]interface{}), errors.New("this is a new error")
}
