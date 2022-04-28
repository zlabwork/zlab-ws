package mysql

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
)

type SessionRepository struct {
	Conn *sql.DB
}

func NewSessionRepository() (*SessionRepository, error) {

	h, err := getHandle()
	if err != nil {
		return nil, err
	}
	return &SessionRepository{Conn: h.Conn}, nil
}

func (sr *SessionRepository) NewUID(sid int64, userIds []int64) error {

	t := time.Now()

	// 1.
	var s []string
	for _, id := range userIds {
		s = append(s, strconv.FormatInt(id, 10))
	}
	data := strings.Join(s, ",")

	// 2.
	stmt, err := sr.Conn.Prepare("INSERT INTO `im_session` (`id`, `title`, `uids`, `ctime`, `mtime`) VALUES (?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(sid, "", data, t, t)
	if err != nil {
		return err
	}

	return nil
}

func (sr *SessionRepository) SetUID(sid int64, userIds []int64) error {

	t := time.Now()

	// 1.
	var s []string
	for _, id := range userIds {
		s = append(s, strconv.FormatInt(id, 10))
	}
	data := strings.Join(s, ",")

	// 2.
	stmt, err := sr.Conn.Prepare("UPDATE `im_session` SET `uids`=?, `mtime`=? WHERE `id`=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(data, t, sid)
	if err != nil {
		return err
	}

	return nil
}

func (sr *SessionRepository) GetUID(sid int64) ([]int64, error) {

	s := ""
	row := sr.Conn.QueryRow("SELECT `uids` FROM `im_session` WHERE `id` = ? LIMIT 1", sid)
	err := row.Scan(&s)
	if err != nil {
		return nil, err
	}
	arr := strings.Split(s, ",")

	var u []int64
	for _, id := range arr {
		i, _ := strconv.ParseInt(id, 10, 64)
		u = append(u, i)
	}

	return u, nil
}
