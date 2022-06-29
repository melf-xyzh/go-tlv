/**
 * @Time    :2022/6/4 6:55
 * @Author  :MELF晓宇
 * @Email   :xyzh.melf@petalmail.com
 * @FileName:node.go
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
	"encoding/hex"
)

// Node TLV节点
type Node struct {
	Tag       uint64           // T
	Length    uint64           // L
	ValueByte []byte           // V
	order     binary.ByteOrder // 字节序
}

// GetChar
/**
 *  @Description: 获取Char字符
 *  @receiver n
 *  @return value
 */
func (n *Node) GetChar() (value byte) {
	return n.ValueByte[0]
}

// GetShort
/**
 *  @Description: 获取short类型值
 *  @receiver n
 *  @return value
 */
func (n *Node) GetUChar() (value int16) {
	var uintV uint8
	buf := bytes.NewBuffer([]byte{})
	// 取到的数据写入缓冲区
	binary.Write(buf, n.order, n.ValueByte)
	binary.Read(buf, n.order, &uintV)
	return int16(uintV)
}

// GetHexBytesString
/**
 *  @Description: 获取HexBytes对应的字符串
 *  @receiver n
 *  @return value
 */
func (n *Node) GetHexBytesString() (value string) {
	return hex.EncodeToString(n.ValueByte)
}

// GetShort
/**
 *  @Description: 获取short类型值
 *  @receiver n
 *  @return value
 */
func (n *Node) GetShort() (value int16) {
	var uintV uint16
	buf := bytes.NewBuffer([]byte{})
	// 取到的数据写入缓冲区
	binary.Write(buf, n.order, n.ValueByte)
	binary.Read(buf, n.order, &uintV)
	return int16(uintV)
}

// GetUint
/**
 *  @Description: 获取uint类型值
 *  @receiver n
 *  @return value
 */
func (n *Node) GetUint() (value uint) {
	var uintV32 uint32
	buf := bytes.NewBuffer([]byte{})
	// 取到的数据写入缓冲区
	binary.Write(buf, n.order, n.ValueByte)
	binary.Read(buf, n.order, &uintV32)
	return uint(uintV32)
}

// GetString
/**
 * @Description: 获取String
 * @receiver n
 * @return value
 */
func (n *Node) GetString() (value string) {
	return string(n.ValueByte)
}

// GetInt
/**
 * @Description: 获取Int类型
 * @receiver n
 * @return v
 */
func (n Node) GetInt() (v int) {
	var uintV uint32
	buf := bytes.NewBuffer([]byte{})
	// 取到的数据写入缓冲区
	binary.Write(buf, n.order, n.ValueByte)
	binary.Read(buf, n.order, &uintV)
	if uintV > 2147483627 {
		return int(4294967295-uintV)*-1 - 1
	}
	return int(uintV)
}
