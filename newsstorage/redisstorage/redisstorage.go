package redisstorage

import (
	"BotNewsScrapper/hotnews"
	"crypto/md5"
	"encoding/hex"
	red "github.com/go-redis/redis"
	"time"
)

type RedisStorage[News hotnews.News] struct {
	Client *red.Client
}

func Init[News hotnews.News](addr string, password string, db int) RedisStorage[News] {
	return RedisStorage[News]{Client: red.NewClient(&red.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})}
}

func (r RedisStorage[News]) Add(news News) bool {
	if r.Contains(news) {
		return false
	}
	return r.Client.Append(GetMD5Hash(news.GetNews()), time.Now().String()).Err() == nil
}

func (r RedisStorage[News]) Free() {
	for _, key := range r.GetAllKeys() {
		t, _ := time.Parse(time.Now().String(), r.Client.Get(key).String())
		if t.Add(time.Hour * 72).Before(time.Now()) {
			r.Client.Del(key)
		}
	}
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func (r RedisStorage[News]) Contains(news News) bool {
	return r.Client.Get(GetMD5Hash(news.GetNews())).Err() == nil
}

func (r RedisStorage[News]) GetAllKeys() []string {
	res, err := r.Client.Keys("*").Result()
	if err != nil {
		return []string{}
	}
	return res
}
