/**
 * @Time    :2022/6/3 10:32
 * @Author  :MELF晓宇
 * @Email   :xyzh.melf@petalmail.com
 * @FileName:main.go
 * @Project :go-tlv
 * @Blog    :https://blog.csdn.net/qq_29537269
 * @Guide   :https://guide.melf.space
 * @Information:
 *
 */

package main

import (
	"encoding/binary"
	"fmt"
	mtlv "go-tlv/tlv"
)

func main() {
	tlv := mtlv.InitTlv(mtlv.SIZE2, mtlv.SIZE2, binary.BigEndian)
	tlvMap := make(map[uint64]interface{})
	tlvMap[0x100+1] = "车载机"
	tlvMap[0x100+1] = 20
	tlvMap[0x100+2] = 3.14
	tlvMap[0x100+3] = -15
	tlvMap[0x100+4] = true
	tlvMap[0x100+5] = false
	//tlvMap[0x100+6] = nil
	a := tlv.WriteData(tlvMap)
	//fmt.Println(a)

	//x := tlv.Read(a)
	//fmt.Println(x.Tag)
	//fmt.Println(x.Length)
	//fmt.Println(string(x.ValueByte))

	xs := tlv.ReadData(a)
	for _, x := range xs {
		fmt.Println(x.Tag, x.Length, string(x.ValueByte))
	}
}
