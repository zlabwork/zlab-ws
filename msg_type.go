package zlabws

type AuthMsg struct {
	Id      string
	From    int64
	Token   string
	Uuid    string
	Version string
	Os      string
	Time    int64
}

type CustomMsg struct {
	Id   string
	From int64
	To   int64
	Type string
	Data string
	Time int64
}

type TextMsg struct {
	Id   string
	From int64
	To   int64
	Text string
	Time int64
}

type ImageMsg struct {
	Id     string
	From   int64
	To     int64
	Width  int32
	Height int32
	Time   int64
}

type VoiceMsg struct {
	Id       string
	From     int64
	To       int64
	Size     int32
	Duration int32
	Time     int64
}

type VideoMsg struct {
	Id       string
	From     int64
	To       int64
	Size     int32
	Duration int32
	Time     int64
}

type FileMsg struct {
	Id   string
	From int64
	To   int64
	Name string
	Mime string
	Size int32
	Time int64
}

type StickerMsg struct {
	Id   string
	From int64
	To   int64
	Name string
	Time int64
}

type LocationMsg struct {
	Id   string
	From int64
	To   int64
	Lat  float64
	Lng  float64
	Desc string
	Time int64
}

type NoticeMsg struct {
	Id   string
	From int64
	To   int64
	Type string
	Data string
	Time int64
}

type CommandMsg struct {
}
