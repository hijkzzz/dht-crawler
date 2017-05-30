package dht

import "errors"

// EncodeBencode bencode 编码
func EncodeBencode(msg map[string]interface{}) ([]byte, error) {
	return nil, errors.New("this is a new error")
}

// DecodeBencode bencode 解码
func DecodeBencode(msg []byte) (map[string]interface{}, error) {
	return make(map[string]interface{}), errors.New("this is a new error")
}
