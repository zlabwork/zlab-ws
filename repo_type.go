package app

import "time"

type RepoMsg struct {
	Id       int64
	Mid      string
	Type     uint8
	Sender   int64
	Receiver int64
	Data     []byte
	Ctime    time.Time
}
