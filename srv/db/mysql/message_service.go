package mysql

import (
	"fmt"
	"zlabws"
)

type MessageService struct {
	h *handle
}

func (m *MessageService) Message(id string) (*zlabws.Message, error) {
	row := m.h.Conn.QueryRow("SELECT `id`,`type`,`from`,`to`,`data`,`ctime` FROM `im_msg` WHERE `id` = ? LIMIT 1", id)
	msg := zlabws.Message{}
	row.Scan(&msg.Id, &msg.Type, &msg.From, &msg.To, &msg.Data, &msg.Ctime)
	if len(msg.Id) == 0 {
		return nil, fmt.Errorf("no data")
	}
	return &msg, nil
}

func (m *MessageService) CreateMsg(msg *zlabws.Message) error {
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

func (m *MessageService) DeleteMsg(id string) error {
	stmt, err := m.h.Conn.Prepare("DELETE FROM im_msg WHERE id = ?")
	if err != nil {
		return err
	}
	if _, err = stmt.Exec(id); err != nil {
		return err
	}
	return nil
}
