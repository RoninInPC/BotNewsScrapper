package simplestorage

import (
	"BotNewsScrapper/hotnews"
	"crypto/md5"
	"encoding/hex"
	"time"
)

type SimpleStorage[News hotnews.News] struct {
	Map map[string]time.Time
}

func Init[News hotnews.News]() SimpleStorage[News] {
	return SimpleStorage[News]{Map: make(map[string]time.Time)}
}

func (s SimpleStorage[News]) Add(news hotnews.News) bool {
	hash := GetMD5Hash(news.GetNews())
	if _, ok := s.Map[hash]; !ok {
		s.Map[hash] = time.Now()
		return true
	}
	return false
}

func (s SimpleStorage[News]) Free() {
	timeNow := time.Now()
	deletedKeys := make([]string, 0)
	for key, value := range s.Map {
		if value.Add(time.Hour * 24).Before(timeNow) {
			deletedKeys = append(deletedKeys, key)
		}
	}
	for _, key := range deletedKeys {
		delete(s.Map, key)
	}
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
