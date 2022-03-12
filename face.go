package app

type CacheFace interface {
	Close() error
	GetToken(id string) (string, error)
	SetToken(id string, token string) error
}

type RepoFace interface {
}
