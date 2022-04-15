## 消息协议规范

```
消息头 (32 Bit)
+-------------------------------------------+
| 8 Bit Unused | 8 Bit Type | 16 Bit Length ｜
+-------------------------------------------+

消息体
+--------------------------------------------------------------------+
| 64 Bit SequenceID | 64 Bit SenderID | 64 Bit ReceiverID | Data Part｜
+--------------------------------------------------------------------+

Data Part [text] 
+------+
| text ｜
+------+

Data Part [image] 
+------------------------------------+
| 32 Bit Width | 32 Bit Height | URI ｜
+------------------------------------+

Data Part [voice] 
+-------------------------------------+
| 32 Bit Size | 32 Bit Duration | URI ｜
+-------------------------------------+

Data Part [video] 
+-------------------------------------+
| 32 Bit Size | 32 Bit Duration | URI ｜
+-------------------------------------+

Data Part [sticker] 
+-----+
| URI ｜
+-----+

Data Part [file]
+----------------------------------+
| 32 Bit Size | Type,URI,Name,Desc ｜
+----------------------------------+

Data Part [location]
+-------------------------------------------------+
| 64 Bit Longitude | 64 Bit Latitude | Title,Desc ｜
+-------------------------------------------------+

Data Part [custom]
+------------------------+
| 8 Bit Type | Json Data ｜
+------------------------+

All types above, Use 0x1F instead of commas.
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


## Tools
[在线流程图思维导图](https://www.processon.com)  


## Docs
[从0到1：微信后台系统的演进之路](https://mp.weixin.qq.com/s/fMF_FjcdLiXc_JVmf4fl0w)  
[京东客服即时通讯系统的技术架构演进历程](http://www.52im.net/thread-152-1-1.html)  
[移动社交的即时通讯软件是基于怎样的技术架构](https://www.zhihu.com/question/20458376)  
[浅谈聊天系统架构设计](https://www.isolves.com/it/cxkf/jiagou/2020-12-22/34719.html)  
[如何设计一个亿级消息量的 IM 系统](https://xie.infoq.cn/article/19e95a78e2f5389588debfb1c)  
[现代IM系统中的消息系统架构 - 架构篇](https://developer.aliyun.com/article/698301)  
[现代IM系统中的消息系统架构 - 实现篇](https://developer.aliyun.com/article/710363)  
[一套亿级用户的IM架构技术干货(上篇)：整体架构、服务拆分等](https://zhuanlan.zhihu.com/p/357302917)  
[MQTT协议](https://www.runoob.com/w3cnote/mqtt-intro.html)  
