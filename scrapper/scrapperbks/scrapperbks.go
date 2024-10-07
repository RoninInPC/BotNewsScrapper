package scrapperbks

import (
	"BotNewsScrapper/hotnews"
	"BotNewsScrapper/htmlgetter"
	"regexp"
	"strings"
	"time"
)

const (
	BaseURL = "https://bcs.ru/news"
)

type ScrapperBKS struct {
	HTMLGetter htmlgetter.HTMLGetter
}

func (s ScrapperBKS) Scrape(
	channel chan<- hotnews.WebNews,
	url string,
	duration time.Duration) {
	if url == "" {
		url = BaseURL
	}
	go func() {
		for {
			html, err := s.HTMLGetter.GetHTML(url)
			if err != nil {
				continue
			}

			for _, n := range s.AnalysisHTML(html, url, time.Now().Format("2006-01-02")) {
				channel <- n
			}
			time.Sleep(duration)
		}
	}()
}

func (s ScrapperBKS) AnalysisHTML(html string, url string, timeNow string) []hotnews.WebNews {
	answer := make([]hotnews.WebNews, 0)
	strs := regexp.MustCompile(`"publishDate":\s*"`+timeNow+`",\s*"title":\s*"(.*?)",\s*"slug":\s*"(.*?)"`).FindAllStringSubmatch(html, -1)
	for _, str := range strs {
		if len(str) < 4 {
			continue
		}
		answer = append(answer, hotnews.WebNews{
			From:  hotnews.BKS,
			Title: str[2],
			URL:   url + strings.Replace(str[1], "-", "/", 2) + str[3],
			Time:  str[1]})
	}
	return answer
}
