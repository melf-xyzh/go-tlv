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

package tlv

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

type lengthSize int

const (
	SIZE1 lengthSize = 1 // 1字节
	SIZE2 lengthSize = 2 // 2字节
	SIZE4 lengthSize = 4 // 4字节
	SIZE8 lengthSize = 8 // 8字节
)

// TlvConfig Tlv配置
type TlvConfig struct {
	tagSize     lengthSize       // Tag长度（字节）
	lengthSize  lengthSize       // Length长度（字节）
	minNodeSize lengthSize       // Node最小长度（字节）
	order       binary.ByteOrder // 字节序（大端存储 / 小端存储）
}

// TlvContent
type TlvContent struct {
	Tag       uint64 // T
	Length    uint64 // L
	ValueByte []byte // V
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

func (tlvc *TlvConfig) TypeToByte(v interface{}) []byte {

	typeOfA := reflect.TypeOf(v)
	fmt.Println(typeOfA.Name(), typeOfA.Kind())

	switch v.(type) {
	case string:
		return []byte(v.(string))
	case []byte:
		return v.([]byte)
	case int:
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, v.(int))
		return buf.Bytes()
	case float64:
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, v.(float64))
		return buf.Bytes()
	case bool:
		buf := bytes.NewBuffer([]byte{})
		binary.Write(buf, binary.BigEndian, v.(bool))
		return buf.Bytes()
	default:
		return nil
	}
}

// WriteData
/**
 * @Description: 生成Tlv报文
 * @receiver tlvc
 * @param data
 * @return dataByte
 */
func (tlvc *TlvConfig) WriteData(data map[uint64]interface{}) (dataByte []byte) {
	buf := bytes.NewBuffer([]byte{})
	for k, v := range data {
		vByte := tlvc.TypeToByte(v)
		writeBt := tlvc.Write(k, vByte)
		buf.Write(writeBt)
	}
	return buf.Bytes()
}

// NumToTag
/**
 * @Description: 数字作为Tag
 * @receiver tlv
 * @param num
 * @return tag
 */
func (tlvc *TlvConfig) Write(tag uint64, value []byte) (returnByte []byte) {
	buf := bytes.NewBuffer([]byte{})
	// 写入tag(T)
	switch tlvc.tagSize {
	case 1:
		binary.Write(buf, tlvc.order, uint8(tag))
	case 2:
		binary.Write(buf, tlvc.order, uint16(tag))
	case 4:
		binary.Write(buf, tlvc.order, uint32(tag))
	case 8:
		binary.Write(buf, tlvc.order, tag)
	}
	// 写入长度(L)
	switch tlvc.tagSize {
	case 1:
		binary.Write(buf, tlvc.order, uint8(len(value)))
	case 2:
		binary.Write(buf, tlvc.order, uint16(len(value)))
	case 4:
		binary.Write(buf, tlvc.order, uint32(len(value)))
	case 8:
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
func (tlvc *TlvConfig) Read(tlvByte []byte) (TlvContent TlvContent) {
	//接收数据缓冲区
	bio := bytes.NewBuffer([]byte{})
	// 取到的数据写入缓冲区
	binary.Write(bio, tlvc.order, tlvByte)
	// 读取tag(T)
	switch tlvc.tagSize {
	case 1:
		var tag uint8
		binary.Read(bio, tlvc.order, &tag)
		TlvContent.Tag = uint64(tag)
	case 2:
		var tag uint16
		binary.Read(bio, tlvc.order, &tag)
		TlvContent.Tag = uint64(tag)
	case 4:
		var tag uint32
		binary.Read(bio, tlvc.order, &tag)
		TlvContent.Tag = uint64(tag)
	case 8:
		binary.Read(bio, tlvc.order, &TlvContent.Tag)
	}
	//tlvByte = tlvByte[tlvc.tagSize:]

	//// 取到的数据写入缓冲区
	//binary.Write(bio, tlvc.order, tlvByte[:tlvc.lengthSize])
	// 读取长度(T)
	switch tlvc.lengthSize {
	case 1:
		var length uint8
		binary.Read(bio, tlvc.order, &length)
		TlvContent.Length = uint64(length)
	case 2:
		var length uint16
		binary.Read(bio, tlvc.order, &length)
		TlvContent.Length = uint64(length)
	case 4:
		var length uint32
		binary.Read(bio, tlvc.order, &length)
		TlvContent.Length = uint64(length)
	case 8:
		binary.Read(bio, tlvc.order, &TlvContent.Length)
	}
	valueByte := make([]byte, TlvContent.Length, TlvContent.Length)
	binary.Read(bio, tlvc.order, &valueByte)
	TlvContent.ValueByte = valueByte
	return
}

// Read
/**
 * @Description: 读取Tlv内容
 * @receiver tlvc
 * @param tlvByte
 * @return TlvContent
 */
func (tlvc *TlvConfig) ReadData(tlvByte []byte) (TlvContents []TlvContent) {
	//接收数据缓冲区
	bio := bytes.NewBuffer([]byte{})
	// 取到的数据写入缓冲区
	binary.Write(bio, tlvc.order, tlvByte)
	for {
		if len(bio.Bytes()) < int(tlvc.tagSize+tlvc.lengthSize) {
			break
		}
		TlvContent := TlvContent{}
		// 读取tag(T)
		switch tlvc.tagSize {
		case 1:
			var tag uint8
			binary.Read(bio, tlvc.order, &tag)
			TlvContent.Tag = uint64(tag)
		case 2:
			var tag uint16
			binary.Read(bio, tlvc.order, &tag)
			TlvContent.Tag = uint64(tag)
		case 4:
			var tag uint32
			binary.Read(bio, tlvc.order, &tag)
			TlvContent.Tag = uint64(tag)
		case 8:
			binary.Read(bio, tlvc.order, &TlvContent.Tag)
		}
		//tlvByte = tlvByte[tlvc.tagSize:]

		//// 取到的数据写入缓冲区
		//binary.Write(bio, tlvc.order, tlvByte[:tlvc.lengthSize])
		// 读取长度(T)
		switch tlvc.lengthSize {
		case 1:
			var length uint8
			binary.Read(bio, tlvc.order, &length)
			TlvContent.Length = uint64(length)
		case 2:
			var length uint16
			binary.Read(bio, tlvc.order, &length)
			TlvContent.Length = uint64(length)
		case 4:
			var length uint32
			binary.Read(bio, tlvc.order, &length)
			TlvContent.Length = uint64(length)
		case 8:
			binary.Read(bio, tlvc.order, &TlvContent.Length)
		}
		valueByte := make([]byte, TlvContent.Length, TlvContent.Length)
		binary.Read(bio, tlvc.order, &valueByte)
		TlvContent.ValueByte = valueByte
		TlvContents = append(TlvContents, TlvContent)
	}
	return
}
