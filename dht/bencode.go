package dht

import (
	"bytes"
	"errors"
	"reflect"
	"strconv"
)

// encodeBencode bencode 编码
func encodeBencode(msg map[string]interface{}) ([]byte, error) {
	var b bytes.Buffer
	b.WriteByte('t') //buffer第一位用来判断编码类型是否正确
	encodeDictionary(&b, msg)
	var bytes = b.Bytes()
	if bytes[0] == 'f' {
		return []byte{}, errors.New("error: The type of argument does't conform to the bencode specification")
	}

	return bytes[1:len(bytes)], nil
}

// decodeBencode bencode 解码
func decodeBencode(msg []byte) (map[string]interface{}, error) {
	return make(map[string]interface{}), errors.New("this is a new error")
}

/*
*func：int编码
 */
func encodeInt(b *bytes.Buffer, oldInt int64) {
	b.WriteByte('i')
	b.WriteString(strconv.FormatInt(oldInt, 10))
	b.WriteByte('e')
}

/*
*func：uint编码
 */
func encodeUint(b *bytes.Buffer, oldInt uint64) {
	b.WriteByte('i')
	b.WriteString(strconv.FormatUint(oldInt, 10))
	b.WriteByte('e')
}

/*
*func：string编码
 */
func encodeString(b *bytes.Buffer, oldStr string) {
	b.WriteString(strconv.Itoa(len(oldStr)))
	b.WriteByte(':')
	b.WriteString(oldStr)
}

/*
*func：list编码
 */
func encodeList(b *bytes.Buffer, oldList []interface{}) {
	b.WriteByte('l')
	for _, o := range oldList {
		encodeInterface(b, o)
	}
	b.WriteByte('e')

}

/*
*func：dictionary编码
 */
func encodeDictionary(b *bytes.Buffer, oldDir map[string]interface{}) {
	b.WriteByte('d')
	for o := range oldDir {
		encodeString(b, o)
		encodeInterface(b, oldDir[o])
	}
	b.WriteByte('e')
}

/*
*func：interface{}编码
*匹配不同的数据类型
*用于dictionary数据项和list的不同对象的解码
 */
func encodeInterface(b *bytes.Buffer, oldData interface{}) {
	switch oldData := oldData.(type) {
	case int, int8, int16, int32, int64:
		encodeInt(b, reflect.ValueOf(oldData).Int())
	case uint, uint8, uint16, uint32, uint64:
		encodeUint(b, reflect.ValueOf(oldData).Uint())
	case string:
		encodeString(b, oldData)
	case []interface{}:
		encodeList(b, oldData)
	case map[string]interface{}:
		encodeDictionary(b, oldData)
	default:
		var nb bytes.Buffer
		nb.WriteByte('f')
		*b = nb
	}
}

/*
*func：字符串解码
*返回值string：解码结果
*返回值bool：传参oldStr是否符合bencode编码规范
 */
func decodeString(oldStr string) (string, bool) {
	var sLen = ""
	var newStr = ""
	for i := 0; i < len(oldStr); i++ {
		if oldStr[i] == ':' {
			sLen = oldStr[0:i]
			newStr = oldStr[i+1 : len(oldStr)]
			break
		}
	}

	a, error := strconv.Atoi(sLen)
	if error != nil || a != len(newStr) || a < 1 {
		return "", false
	}

	return newStr, true
}
