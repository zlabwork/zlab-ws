package mysql

import (
    "fmt"
    "zlabws"
)

type MsgService struct {
    h *handle
}

func (m *MsgService) Msg(id string) (*zlabws.Msg, error) {
    row := m.h.Conn.QueryRow("SELECT `id`,`type`,`from`,`to`,`data`,`ctime` FROM `im_msg` WHERE `id` = ? LIMIT 1", id)
    msg := zlabws.Msg{}
    row.Scan(&msg.Id, &msg.Type, &msg.From, &msg.To, &msg.Data, &msg.Ctime)
    if len(msg.Id) == 0 {
        return nil, fmt.Errorf("no data")
    }
    return &msg, nil
}

func (m *MsgService) CreateMsg(msg *zlabws.Msg) error {
    stmt, err := m.h.Conn.Prepare("INSERT INTO im_msg (`id`,`type`,`from`,`to`,`data`,`ctime`) VALUES (?,?,?,?,?,?)")
    if err != nil {
        return err
    }
    _, err = stmt.Exec(msg.Id, msg.Type, msg.From, msg.To, msg.Data, msg.Ctime)
    if err != nil {
        return err
    }
    return nil
}

func (m *MsgService) DeleteMsg(id string) error {
    stmt, err := m.h.Conn.Prepare("DELETE FROM im_msg WHERE id = ?")
    if err != nil {
        return err
    }
    if _, err = stmt.Exec(id); err != nil {
        return err
    }
    return nil
}
