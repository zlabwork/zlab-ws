package app

const (
	TypeNull  uint8 = 0x00
	TypeBeats uint8 = 0x01
	TypeAck   uint8 = 0x02
	TypeAuth  uint8 = 0x03
)

const (
	TypeText        uint8 = 0x20 // 文本
	TypeImage       uint8 = 0x21 // 图片
	TypeVoice       uint8 = 0x22 // 语音消息
	TypeVideo       uint8 = 0x23 // 视频消息
	TypeFile        uint8 = 0x24 // 文件
	TypeMedia       uint8 = 0x25 // 媒体
	TypeMusic       uint8 = 0x26 // 音乐媒体
	TypeLive        uint8 = 0x27 // 直播
	TypeLocation    uint8 = 0x28 // 地理位置
	TypeSticker     uint8 = 0x29 // 贴纸
	TypeInteractive uint8 = 0x2A // 互动
	TypeGift        uint8 = 0x2B // 礼物
	TypeNotice      uint8 = 0xA0 // 系统通知
	TypeCommand     uint8 = 0xA1 // 系统指令
	TypeCustom      uint8 = 0xF0 // 自定义消息
)

type MsgAuth struct {
	Token   string
	Uuid    string
	Version string
	Os      string
	Time    int64
}

type MsgCustom struct {
	Type string
	Data string
}

type MsgImage struct {
	Width  int32
	Height int32
	Name   string
}

type MsgVoice struct {
	Size     int32
	Duration int32
	Name     string
}

type MsgVideo struct {
	Size     int32
	Duration int32
	Name     string
}

type MsgFile struct {
	Size int32
	Type string
	Name string
	Desc string
}

type MsgSticker struct {
	Name string
}

type MsgLocation struct {
	Lat  float64
	Lng  float64
	Desc string
}
