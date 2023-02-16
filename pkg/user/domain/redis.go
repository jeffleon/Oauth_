package domain

type RedisRepository interface {
	HSet(hashKey, field, value string) error
	HGet(hashKey, field string) (string, error)
}
