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
)

// Node TLV节点
type Node struct {
	Tag       uint64           // T
	Length    uint64           // L
	ValueByte []byte           // V
	order     binary.ByteOrder // 字节序
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
