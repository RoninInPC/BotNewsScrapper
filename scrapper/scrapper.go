package scrapper

import (
	"BotNewsScrapper/hotnews"
	"time"
)

type AnalysisHTML func(string) []hotnews.News

type Scrapper[news hotnews.News] interface {
	Scrape(chan<- news, string, time.Duration)
	AnalysisHTML(html string, url string, timeNow string) []news
}
