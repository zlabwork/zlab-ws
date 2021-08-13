package zlabws

type User struct {
    Id int64
}

type UserService interface {
    User(id int64) (*User, error)
}
