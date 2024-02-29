package repository

type IStore interface {
	Set(key, value string, ttl int)
	Get(key string) (string, bool)
	Delete(key string)
}
