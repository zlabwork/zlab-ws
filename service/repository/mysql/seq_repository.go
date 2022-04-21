package mysql

import (
	"database/sql"
	"strconv"
)

type SeqRepository struct {
	Conn *sql.DB
}

func NewSeqRepository() (*SeqRepository, error) {

	h, err := getHandle()
	if err != nil {
		return nil, err
	}
	return &SeqRepository{Conn: h.Conn}, nil
}

func (sr *SeqRepository) GetAll() (map[int64]uint64, error) {

	// 1. Query
	rows, err := sr.Conn.Query("SELECT `id`, `max` FROM `im_seq`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 2. Scan
	var result = make(map[int64]uint64)
	for rows.Next() {
		var k int64
		var v uint64
		err := rows.Scan(&k, &v)
		if err != nil {
			return nil, err
		}
		result[k] = v
	}
	if rows.Err() != nil {
		return nil, err
	}

	// 3.
	return result, nil
}

func (sr *SeqRepository) SetMax(section int64, value uint64) error {

	stmt, err := sr.Conn.Prepare("UPDATE `im_seq` SET `max` = ? WHERE `id` = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(value, section)
	if err != nil {
		return err
	}

	return nil
}

func (sr *SeqRepository) FillInitData(section int64, value int64) error {

	fix := "INSERT INTO `im_seq` (`id`, `max`) VALUES"
	data := ""
	v := strconv.FormatInt(value, 10)

	for i := int64(0); i < section; i++ {
		sec := strconv.FormatInt(i, 10)
		data += "(" + sec + "," + v + "),"

		// 分段执行
		if (i+1)%10 == 0 || i+1 == section {
			sql := fix + data[0:len(data)-1]
			data = ""
			stmt, err := sr.Conn.Prepare(sql)
			if err != nil {
				return err
			}
			_, err = stmt.Exec()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
