package dht

import (
	"errors"
	"strconv"
	"strings"
)

// encodeBencode bencode 编码
func encodeBencode(msg map[string]interface{}) ([]byte, error) {
	if msg == nil || len(msg) == 0 {
		return nil, errors.New("msg is null")
	}
	var err error
	str, err := encodeDict(msg)
	if err != nil {
		return nil, errors.New("EncodeBencode")
	}
	return []byte(str), nil
}

// decodeBencode bencode 解码
func decodeBencode(msg []byte) (map[string]interface{}, error) {

	if msg == nil || len(msg) == 0 {
		return nil, errors.New("msg is null")
	}
	var result map[string]interface{}
	var err error
	result, _, err = decodeDict(msg, 0)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// encodeInterface
func encodeInterface(data interface{}) (string, error) {
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
		return "error", errors.New("invalid type when encode")
	}
}

// 将int类型转化为bencode编码
func encodeInt(data int) (string, error) {
	return strings.Join([]string{"i", "e"}, strconv.Itoa(data)), nil
}

// 将string类型转化为bencode编码
func encodeString(data string) (string, error) {
	return strings.Join([]string{strconv.Itoa(len(data)), data}, ":"), nil
}

// 将list类型转化为bencode编码
func encodeList(data []interface{}) (string, error) {
	result := ""
	for _, item := range data {
		var err error
		str, err := encodeInterface(item)
		if err != nil {
			return "error", errors.New("encodeList")
		}
		result = strings.Join([]string{result, str}, "")
	}
	return strings.Join([]string{"l", result, "e"}, ""), nil
}

// 将map类型转化为bencode编码
func encodeDict(data map[string]interface{}) (string, error) {
	result := ""
	for key, item := range data {
		str1, _ := encodeString(key)
		var err error
		str2, err := encodeInterface(item)
		if err != nil {
			return "error", errors.New("encodeDict")
		}
		result = strings.Join([]string{result, str1, str2}, "")
	}
	return strings.Join([]string{"d", result, "e"}, ""), nil
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
		return nil, -1, errors.New("invalid type when encode")
	}
}

// 解码int类型
func decodeInt(data []byte, start int) (int, int, error) {
	end := start + 1
	for data[end] != 'e' {
		end++
		if end == len(data) {
			return -1, -1, errors.New("e not found when decode int")
		}
	}
	num, _ := strconv.Atoi(string(data[start+1 : end]))
	return num, end + 1, nil
}

// 解码string类型
func decodeString(data []byte, start int) (string, int, error) {
	middle := start + 1
	for data[middle] != ':' {
		middle++
		if middle == len(data) {
			return "", -1, errors.New("： not found when decode string")
		}
	}
	len, _ := strconv.Atoi(string(data[start:middle]))
	end := middle + len + 1
	return string(data[middle+1 : end]), end, nil
}

// 解码List类型
func decodeList(data []byte, start int) ([]interface{}, int, error) {
	list := make([]interface{}, 0)
	end := start + 1
	for end < len(data) {
		var item interface{}
		var err error
		item, end, err = decodeInterface(data, end)
		if err != nil {
			return nil, -1, errors.New("decodeList")
		}
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
		var err1 error
		key, end, err1 = decodeInterface(data, end)
		if err1 != nil {
			return nil, -1, errors.New("invalid type of dictionary")
		}
		switch key.(type) {
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
		var err2 error
		item, end, err2 = decodeInterface(data, end)
		if err2 != nil {
			return nil, -1, errors.New("invalid type of dictionary")
		}
		dict[key.(string)] = item
		if end == len(data) {
			return nil, -1, errors.New("e not found when decode list")
		}
		if data[end] == 'e' {
			break
		}
	}
	return dict, end + 1, nil
}
