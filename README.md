# go-tlv

### 介绍

一个TLV（Tag、Length、Value）报文的简易封装库。

### 安装

```
go get github.com/melf-xyzh/go-tlv
```

### 例子

```go
package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/melf-xyzh/go-tlv/tlv"
)

func main() {
	// 初始化TLV配置
	tlv := gtlv.InitTlv(gtlv.SIZE2, gtlv.SIZE2, binary.BigEndian)
	// 创建TLV Map
	tlvMap := make(map[uint64]interface{})
	// Char
	tlvMap[0x100] = byte(5)
	// HexBytes
	deviceId := "12233221"
	tlvMap[0x100+1], _ = hex.DecodeString(deviceId)
	// Short
	tlvMap[0x100+2] = int16(22)
	// Uint
	tlvMap[0x100+3] = int16(16)

	// 生成TLV报文
	a := tlv.WriteNodes(tlvMap)
	// 读取TLV报文
	xs := tlv.ReadNodes(a)
	// 打印报文结果
	for _, x := range xs {
		switch x.Tag {
		case 0x100:
			fmt.Println(x.Tag, x.Length, x.GetChar())
		case 0x100 + 1:
			fmt.Println(x.Tag, x.Length, x.GetHexBytesString())
		case 0x100 + 2:
			fmt.Println(x.Tag, x.Length, x.GetShort())
		case 0x100 + 3:
			fmt.Println(x.Tag, x.Length, x.GetUint())
		}
	}
}
```