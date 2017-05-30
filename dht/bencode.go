package dht

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"
	"strings"
)

func main() {
	s := "123"
	print(strings.ContainsRune(s, 1))
}

// EncodeBencode bencode 编码
func EncodeBencode(msg map[string]interface{}) ([]byte, error) {
	if msg == nil {
		return nil, errors.New("msg is null")
	}
	return []byte(encodeDict(msg)), nil
}

// encodeInterface
func encodeInterface(data interface{}) string {
	switch data.(type) {
	case int:
		return encodeInt(data.(int))
	case string:
		return encodeString(data.(string))
	case []interface{}:
		return encodeList(data.([]interface{}))
	case map[string]interface{}:
		return encodeDict(data.(map[string]interface{}))
	default:
		panic("invalid type when encode")
	}
}

// 将int类型转化为bencode编码
func encodeInt(data int) string {
	return strings.Join([]string{"i", "e"}, strconv.Itoa(data))
}

// 将string类型转化为bencode编码
func encodeString(data string) string {
	return strings.Join([]string{strconv.Itoa(len(data)), data}, ":")
}

// 将list类型转化为bencode编码
func encodeList(data []interface{}) string {
	result := ""
	for _, item := range data {
		result = strings.Join([]string{result, encodeInterface(item)}, "")
	}
	return strings.Join([]string{"l", result, "e"}, "")
}

// 将map类型转化为bencode编码
func encodeDict(data map[string]interface{}) string {
	result := ""
	for key, item := range data {
		result = strings.Join([]string{result, encodeString(key), encodeInterface(item)}, "")
	}
	return strings.Join([]string{"d", result, "e"}, "")
}

// DecodeBencode bencode 解码
func DecodeBencode(msg []byte) (map[string]interface{}, error) {
	if msg == nil {
		return nil, errors.New("msg is null")
	}
	var result map[string]interface{}
	result, _, _ = decodeDict(msg, 0)
	return result, nil
}

// decodeInterface
func decodeInterface(data []byte, start int) (interface{}, int, error) {
	switch data[start] {
	case 'i':
		return decodeInt(data, start)
	case '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return decodeString(data, start)
	case 'l':
		return decodeList(data, start)
	case 'd':
		return decodeDict(data, start)
	default:
		panic("invalid type when decode")
	}
}

// 解码int类型
func decodeInt(data []byte, start int) (int, int, error) {
	end := start + 1
	for data[end] == 'e' {
		end++
		if end == len(data) {
			return -1, -1, errors.New("e not found when decode int")
		}
	}
	return bytesToInt(data[start+1 : end]), end + 1, nil
}

// 解码string类型
func decodeString(data []byte, start int) (string, int, error) {
	middle := start + 1
	for data[middle] == ':' {
		middle++
		if middle == len(data) {
			return "", -1, errors.New("： not found when decode string")
		}
	}
	len := bytesToInt(data[start:middle])
	end := middle + len + 1
	return string(data[middle+1 : end]), end + 1, nil
}

// 解码List类型
func decodeList(data []byte, start int) ([]interface{}, int, error) {
	list := make([]interface{}, 0)
	end := start + 1
	for end < len(data) {
		var item interface{}
		item, end, _ = decodeInterface(data, end)
		list = append(list, item)
		if data[end] == 'e' {
			break
		}
		if end == len(data) {
			return nil, -1, errors.New("e not found when decode list")
		}
	}
	return list, end + 1, nil
}

// 解码字典类型
func decodeDict(data []byte, start int) (map[string]interface{}, int, error) {
	dict := make(map[string]interface{})
	end := start + 1
	for end < len(data) {
		var key interface{}
		var item interface{}
		key, end, _ = decodeInterface(data, end)
		switch item.(type) {
		case string:
		default:
			return nil, -1, errors.New("invalid type of dictionary")
		}
		if data[end] == 'e' {
			return nil, -1, errors.New("invalid type of dictionary")
		}
		if end == len(data) {
			return nil, -1, errors.New("e not found when decode list")
		}
		item, end, _ = decodeInterface(data, end)
		dict[key.(string)] = item
		if data[end] == 'e' {
			break
		}
		if end == len(data) {
			return nil, -1, errors.New("e not found when decode list")
		}
	}
	return dict, end + 1, nil
}

// 字节数组转int
func bytesToInt(b []byte) int {
	byteBuffer := bytes.NewBuffer(b)
	var tmp int
	binary.Read(byteBuffer, binary.BigEndian, &tmp)
	return tmp
}
