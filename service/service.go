package service

type CacheService interface {
	Close() error
	GetToken(id string) (*string, error)
	SetToken(id string, token *string) error
}
