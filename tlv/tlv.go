/**
 * @Time    :2022/6/3 10:35
 * @Author  :MELF晓宇
 * @Email   :xyzh.melf@petalmail.com
 * @FileName:tlv.go
 * @Project :go-tlv
 * @Blog    :https://blog.csdn.net/qq_29537269
 * @Guide   :https://guide.melf.space
 * @Information:
 *
 */

package gtlv

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

type lengthSize int

// TlvConfig Tlv配置
type TlvConfig struct {
	tagSize     lengthSize       // Tag长度（字节）
	lengthSize  lengthSize       // Length长度（字节）
	minNodeSize lengthSize       // Node最小长度（字节）
	order       binary.ByteOrder // 字节序（大端存储 / 小端存储）
}

// InitTlv
/**
 * @Description: 初始化TLV
 * @param tagSize
 * @param lengthSize
 * @param order
 * @return t
 */
func InitTlv(tagSize, lengthSize lengthSize, order binary.ByteOrder) (t TlvConfig) {
	return TlvConfig{
		tagSize:     tagSize,
		lengthSize:  lengthSize,
		minNodeSize: tagSize + lengthSize,
		order:       order,
	}
}

// TypeToByte
/**
 * @Description: 格式转换（转为byte）
 * @receiver tlvc
 * @param v
 * @return []byte
 */
func (tlvc *TlvConfig) TypeToByte(v interface{}, order binary.ByteOrder) []byte {
	switch v.(type) {
	case []byte:
		return v.([]byte)
	case int16:
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, order, uint16(v.(int16)))
		return buf.Bytes()
	case uint8:
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, order, v.(uint8))
		return buf.Bytes()
	case uint:
		var uintV uint = v.(uint)
		uint32V := uint32(uintV)
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, order, uint32V)
		return buf.Bytes()
	case string:
		return []byte(v.(string))
	case int:
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, order, uint32(v.(int)))
		return buf.Bytes()
	case float64:
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, order, v.(float64))
		return buf.Bytes()
	case bool:
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, order, v.(bool))
		return buf.Bytes()
	default:
		typeOfA := reflect.TypeOf(v)
		fmt.Println(typeOfA.Name(), typeOfA.Kind(), "不支持")
		return nil
	}
}

// WriteNodes
/**
 * @Description: 生成Tlv报文
 * @receiver tlvc
 * @param data
 * @return dataByte
 */
func (tlvc *TlvConfig) WriteNodes(data map[uint64]interface{}) (dataByte []byte) {
	buf := bytes.NewBuffer([]byte{})
	for k, v := range data {
		vByte := tlvc.TypeToByte(v, tlvc.order)
		writeBt := tlvc.Write(k, vByte)
		buf.Write(writeBt)
	}
	return buf.Bytes()
}

// Write
/**
 * @Description: 数字作为Tag
 * @receiver tlvc
 * @param tag
 * @param value
 * @return returnByte
 */
func (tlvc *TlvConfig) Write(tag uint64, value []byte) (returnByte []byte) {
	buf := bytes.NewBuffer([]byte{})
	// 写入tag(T)
	switch tlvc.tagSize {
	case SIZE1:
		binary.Write(buf, tlvc.order, uint8(tag))
	case SIZE2:
		binary.Write(buf, tlvc.order, uint16(tag))
	case SIZE4:
		binary.Write(buf, tlvc.order, uint32(tag))
	case SIZE8:
		binary.Write(buf, tlvc.order, tag)
	}
	// 写入长度(L)
	switch tlvc.tagSize {
	case SIZE1:
		binary.Write(buf, tlvc.order, uint8(len(value)))
	case SIZE2:
		binary.Write(buf, tlvc.order, uint16(len(value)))
	case SIZE4:
		binary.Write(buf, tlvc.order, uint32(len(value)))
	case SIZE8:
		binary.Write(buf, tlvc.order, uint64(len(value)))
	}
	//写入数据(V)
	binary.Write(buf, tlvc.order, value)
	return buf.Bytes()
}

// Read
/**
 * @Description: 读取Tlv内容
 * @receiver tlvc
 * @param tlvByte
 * @return TlvContent
 */
func (tlvc *TlvConfig) Read(tlvByte []byte) (TlvContent Node) {
	//接收数据缓冲区
	bio := bytes.NewBuffer([]byte{})
	// 取到的数据写入缓冲区
	binary.Write(bio, tlvc.order, tlvByte)
	// 读取tag(T)
	switch tlvc.tagSize {
	case SIZE1:
		var tag uint8
		binary.Read(bio, tlvc.order, &tag)
		TlvContent.Tag = uint64(tag)
	case SIZE2:
		var tag uint16
		binary.Read(bio, tlvc.order, &tag)
		TlvContent.Tag = uint64(tag)
	case SIZE4:
		var tag uint32
		binary.Read(bio, tlvc.order, &tag)
		TlvContent.Tag = uint64(tag)
	case SIZE8:
		binary.Read(bio, tlvc.order, &TlvContent.Tag)
	}
	// 读取长度(T)
	switch tlvc.lengthSize {
	case SIZE1:
		var length uint8
		binary.Read(bio, tlvc.order, &length)
		TlvContent.Length = uint64(length)
	case SIZE2:
		var length uint16
		binary.Read(bio, tlvc.order, &length)
		TlvContent.Length = uint64(length)
	case SIZE4:
		var length uint32
		binary.Read(bio, tlvc.order, &length)
		TlvContent.Length = uint64(length)
	case SIZE8:
		binary.Read(bio, tlvc.order, &TlvContent.Length)
	}
	valueByte := make([]byte, TlvContent.Length, TlvContent.Length)
	binary.Read(bio, tlvc.order, &valueByte)
	TlvContent.ValueByte = valueByte
	return
}

// ReadNodes
/**
 * @Description: 读取Tlv内容
 * @receiver tlvc
 * @param tlvByte
 * @return Nodes
 */
func (tlvc *TlvConfig) ReadNodes(tlvByte []byte) (Nodes []Node) {
	//接收数据缓冲区
	bio := bytes.NewBuffer([]byte{})
	// 取到的数据写入缓冲区
	binary.Write(bio, tlvc.order, tlvByte)
	for {
		if len(bio.Bytes()) < int(tlvc.tagSize+tlvc.lengthSize) {
			break
		}
		node := Node{}
		// 读取tag(T)
		switch tlvc.tagSize {
		case SIZE1:
			var tag uint8
			binary.Read(bio, tlvc.order, &tag)
			node.Tag = uint64(tag)
		case SIZE2:
			var tag uint16
			binary.Read(bio, tlvc.order, &tag)
			node.Tag = uint64(tag)
		case SIZE4:
			var tag uint32
			binary.Read(bio, tlvc.order, &tag)
			node.Tag = uint64(tag)
		case SIZE8:
			binary.Read(bio, tlvc.order, &node.Tag)
		}
		// 读取长度(T)
		switch tlvc.lengthSize {
		case SIZE1:
			var length uint8
			binary.Read(bio, tlvc.order, &length)
			node.Length = uint64(length)
		case SIZE2:
			var length uint16
			binary.Read(bio, tlvc.order, &length)
			node.Length = uint64(length)
		case SIZE4:
			var length uint32
			binary.Read(bio, tlvc.order, &length)
			node.Length = uint64(length)
		case SIZE8:
			binary.Read(bio, tlvc.order, &node.Length)
		}
		valueByte := make([]byte, node.Length, node.Length)
		binary.Read(bio, tlvc.order, &valueByte)
		node.ValueByte = valueByte
		node.order = tlvc.order
		Nodes = append(Nodes, node)
	}
	return
}
