package zlabws

type User struct {
	Id     int64
	Name   string
	Gender uint8
	Avatar string
	Desc   string
}

type UserService interface {
	User(id int64) (*User, error)
}
