package zlabws

type Message struct {
    Id    string
    Type  uint8
    From  int64
    To    int64
    Data  []byte
    Ctime int64
}

type MessageService interface {
    Message(id string) (*Message, error)
    CreateMsg(msg *Message) error
    DeleteMsg(id string) error
}
