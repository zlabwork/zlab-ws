package zlabws

type Msg struct {
    Id   string
    Type uint8
    From int64
    To   int64
    Data []byte
}

type MsgService interface {
    Msg(id int64) (*Msg, error)
}
