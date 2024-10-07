package newsstorage

import (
	"BotNewsScrapper/hotnews"
)

type NewsStorage[News hotnews.News] interface {
	Add(News) bool
	Free()
}
