package zlabws

type Msg struct {
    Id    string
    Type  uint8
    From  int64
    To    int64
    Data  []byte
    Ctime int64
}

type MsgService interface {
    Msg(id string) (*Msg, error)
    CreateMsg(msg *Msg) error
    DeleteMsg(id string) error
}
