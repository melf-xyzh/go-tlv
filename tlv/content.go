/**
 * @Time    :2022/6/3 11:01
 * @Author  :MELF晓宇
 * @Email   :xyzh.melf@petalmail.com
 * @FileName:contant.go
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
)

type TLVContent struct {
	Type   []byte
	Length uint64
	Value  []byte
}

func (t *TLVContent) WriteToBytes() (b []byte) {
	newBuffer := bytes.NewBuffer([]byte{})
	//写入tag(T)
	binary.Write(newBuffer, binary.BigEndian, t.Type)
	//写入长度(L)
	binary.Write(newBuffer, binary.BigEndian, t.Length)
	//写入数据(V)
	binary.Write(newBuffer, binary.BigEndian, t.Value)
	return newBuffer.Bytes()
}
