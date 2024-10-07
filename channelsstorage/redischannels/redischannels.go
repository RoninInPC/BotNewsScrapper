package redischannels

import (
	red "github.com/go-redis/redis"
	"strconv"
)

type RedisChannels struct {
	Client *red.Client
}

func Init(addr string, password string, db int) RedisChannels {
	return RedisChannels{Client: red.NewClient(&red.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})}
}

func (r RedisChannels) Add(id int64) {
	converted := strconv.FormatInt(id, 10)
	if r.Client.Get(converted).Err() != nil {
		r.Client.Append(converted, "Added")
	}
}

func (r RedisChannels) GetChatsId() []int64 {
	answer := make([]int64, 0)
	for _, key := range r.Client.Keys("*").Val() {
		i, _ := strconv.ParseInt(key, 10, 64)
		answer = append(answer, i)
	}
	return answer
}

func (r RedisChannels) Size() int {
	return len(r.Client.Keys("*").Val())
}

func (r RedisChannels) Delete(id int64) {
	r.Client.Del(strconv.FormatInt(id, 10))
}
