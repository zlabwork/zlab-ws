package mysql

import (
	"database/sql"
	"time"
)

type RepoFace struct {
	Conn *sql.DB
}

func NewRepoFace() (*RepoFace, error) {

	h, err := getHandle()
	if err != nil {
		return nil, err
	}
	return &RepoFace{Conn: h.Conn}, nil
}

func (rf *RepoFace) Create(msgType uint8, msgId string, sender, receiver int64, body []byte, date time.Time) error {

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

func (rf *RepoFace) CreateTodo(msgType uint8, msgId string, sender, receiver int64, body []byte, date time.Time) error {

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
