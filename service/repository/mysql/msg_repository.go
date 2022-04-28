package mysql

import (
	"app"
	"database/sql"
	"time"
)

type MessageRepository struct {
	Conn *sql.DB
}

func NewMessageRepository() (*MessageRepository, error) {

	h, err := getHandle()
	if err != nil {
		return nil, err
	}
	return &MessageRepository{Conn: h.Conn}, nil
}

func (rf *MessageRepository) GetTodo(userId int64) ([]*app.RepoMsg, error) {

	// 1. Query
	rows, err := rf.Conn.Query("SELECT `id`, `mid`, `type`, `sender`, `receiver`, `data`, `ctime` FROM `im_todo` WHERE `receiver` = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 2. Scan
	var result []*app.RepoMsg
	for rows.Next() {
		msg := &app.RepoMsg{}
		err := rows.Scan(&msg.Id, &msg.Mid, &msg.Type, &msg.Sender, &msg.Receiver, &msg.Data, &msg.Ctime)
		if err != nil {
			return nil, err
		}
		result = append(result, msg)
	}
	if rows.Err() != nil {
		return nil, err
	}

	// 4.
	return result, nil
}

func (rf *MessageRepository) DeleteTodo(userId int64) error {

	stmt, err := rf.Conn.Prepare("DELETE FROM `im_todo` WHERE `receiver` = ?")
	if err != nil {
		return err
	}
	if _, err = stmt.Exec(userId); err != nil {
		return err
	}
	return nil
}

func (rf *MessageRepository) CreateTodo(msgType uint8, msgId string, sender, receiver int64, body []byte, date time.Time) error {

	stmt, err := rf.Conn.Prepare("INSERT INTO `im_todo` (`mid`, `type`, `sender`, `receiver`, `data`, `ctime`) VALUES (?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(msgId, msgType, sender, receiver, body, date)
	if err != nil {
		return err
	}
	return nil
}

func (rf *MessageRepository) CreateLogs(msgType uint8, msgId string, sender, receiver int64, body []byte, date time.Time) error {

	stmt, err := rf.Conn.Prepare("INSERT INTO `im_msg` (`mid`, `type`, `sender`, `receiver`, `data`, `ctime`) VALUES (?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(msgId, msgType, sender, receiver, body, date)
	if err != nil {
		return err
	}
	return nil
}
