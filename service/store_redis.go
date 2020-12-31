package adash

import (
	"errors"
	"strconv"
	"strings"

	"github.com/mediocregopher/radix/v3"
)

// RedisStore хранилище данных полученных от трейдинговых алгоритмов.
type RedisStore struct {
	pool *radix.Pool
	db   int
}

func NewRedisStore(pool *radix.Pool, db int) *RedisStore {
	return &RedisStore{pool: pool, db: db}
}

// Scan returns all instance names startings with "pattern".
func (rs *RedisStore) Scan(pattern string, result *[]string) error {
	*result = (*result)[0:0]

	s := radix.NewScanner(rs.pool, radix.ScanOpts{Command: "SCAN", Pattern: pattern})
	var key string
	for s.Next(&key) {
		*result = append(*result, key)
	}

	return s.Close()
}

func (rs *RedisStore) Keys(pattern string, result *[]string) error {
	*result = (*result)[0:0]

	s := radix.NewScanner(rs.pool, radix.ScanOpts{Command: "KEYS", Key: pattern})
	var key string
	for s.Next(&key) {
		*result = append(*result, key)
	}

	return s.Close()
}

func (rs *RedisStore) GetObject(key string, result interface{}) error {
	return rs.pool.Do(radix.Cmd(result, "HGETALL", key))
}

func (rs *RedisStore) Get(key string, val interface{}) error {
	return rs.pool.Do(radix.Cmd(val, "GET", key))
}

func (rs *RedisStore) MGet(keys []string, val *[]string) error {
	return rs.pool.Do(radix.Cmd(val, "MGET", keys...))
}

func (rs *RedisStore) Set(key string, val string) error {
	err := rs.pool.Do(radix.Cmd(nil, "SET", key, val))
	if err != nil {
		return errors.New("redis key " + key + " setting failed. Detais: " + err.Error())
	}
	return err
}

func (rs *RedisStore) Publish(code string) (bool, error) {

	var res int
	err := rs.pool.Do(radix.Cmd(&res, "PUBLISH", CMND+":"+strconv.Itoa(rs.db), code))
	if err != nil {
		return false, errors.New("redis publish " + code + " to channel " + CMND + " failed. Detais: " + err.Error())
	}
	return res > 0, err
}

func (rs *RedisStore) ListRange(key string, from, to int) ([]string, error) {

	var res []string
	err := rs.pool.Do(radix.Cmd(&res, "LRANGE", key, strconv.Itoa(from), strconv.Itoa(to)))
	if err != nil {
		return nil, errors.New("redis lrange " + key + " failed. Detais: " + err.Error())
	}

	return res, nil
}

func (rs *RedisStore) ListLen(key string) (int, error) {

	keyb := []byte(key)

	if strings.HasPrefix(key, "M:") {
		keyb[0] = 'L'
		key = string(keyb)
	} else if strings.HasPrefix(key, "L:") == false {
		key = "L:" + key
	}

	var res int
	err := rs.pool.Do(radix.Cmd(&res, "LLEN", key))
	if err != nil {
		return 0, errors.New("redis lrange " + key + " failed. Detais: " + err.Error())
	}

	return res, nil
}
func (rs *RedisStore) LLen(key string) (int, error) {

	var res int
	err := rs.pool.Do(radix.Cmd(&res, "LLEN", key))
	if err != nil {
		return 0, errors.New("redis lrange " + key + " failed. Detais: " + err.Error())
	}

	return res, nil
}
func (rs *RedisStore) HGet(key, field string) (string, error) {
	var res string
	err := rs.pool.Do(radix.Cmd(&res, "HGET", key, field))
	if err != nil {
		return "", errors.New("redis HGET " + key + " getting failed. Detais: " + err.Error())
	}
	return res, err
}
