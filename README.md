## 消息协议规范

```
消息头 (32 Bit)
+-------------------------------------------+
| 8 Bit Unused | 8 Bit Type | 16 Bit Length ｜
+-------------------------------------------+

消息体
+----------------------------------------------------------------+
| 64 Bit SequenceID | 64 Bit SenderID | 64 Bit ReceiverID | Data ｜
+----------------------------------------------------------------+
```

```golang
// Data 部分使用单元分隔符分隔
const (
	nul uint8 = 0x00 // 空字符
	lf  uint8 = 0x0A // 换行
	cr  uint8 = 0x0D // 回车键
	fs  uint8 = 0x1C // 文件分隔符
	gs  uint8 = 0x1D // 组分隔符
	rs  uint8 = 0x1E // 记录分隔符
	us  uint8 = 0x1F // 单元分隔符
)
```

```golang
// bytes
bs[0] = 0x1F
copy(bs, bs1)
bytes.Split(bs, []byte{0x1F})
```

## TODO LIST
1. int64 js 不能正确解析问题
