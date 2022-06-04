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
	"github.com/melf-xyzh/go-tlv/tlv"
)

func main() {
	// 初始化TLV配置
	tlv := gtlv.InitTlv(gtlv.SIZE2, gtlv.SIZE2, binary.BigEndian)
	// 创建TLV Map
	tlvMap := make(map[uint64]interface{})
	tlvMap[0x100] = "车载机"
	tlvMap[0x100+1] = 20
	tlvMap[0x100+2] = -15
	// 生成TLV报文
	a := tlv.WriteNodes(tlvMap)
	// 读取TLV报文
	xs := tlv.ReadNodes(a)
	// 打印报文结果
	for _, x := range xs {
		if x.Tag == 0x100 {
			fmt.Println(x.Tag, x.Length, x.GetString())
		} else {
			fmt.Println(x.Tag, x.Length, x.GetInt())
		}
	}
}
