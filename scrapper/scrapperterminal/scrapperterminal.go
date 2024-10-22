package scrapperterminal

import (
	"BotNewsScrapper/hotnews"
	"BotNewsScrapper/htmlgetter"
	"log"
	"regexp"
	"strings"
	"time"
)

var (
	terminalURL = "https://blackterminal.com/news?hl=ru"
)

type ScrapperTerminal struct {
	HTMLGetter htmlgetter.HTMLGetter
}

func (s ScrapperTerminal) Scrape(channel chan<- hotnews.WebNews, url string, duration time.Duration) {
	if url == "" {
		url = terminalURL
	}
	go func() {
		for {
			body, err := s.HTMLGetter.GetHTML(url)
			if err != nil {
				log.Println("Terminal GetHTML err ", err.Error())
				continue
			}
			for _, n := range s.AnalysisHTML(body, url, time.Now().Format("2006-01-02")) {
				channel <- n
			}
			time.Sleep(duration)
		}
	}()
}

func (s ScrapperTerminal) AnalysisHTML(html string, url string, timeNow string) []hotnews.WebNews {
	answer := make([]hotnews.WebNews, 0)
	html = fix(html)
	strs := regexp.MustCompile("<div class=\"align-self-center news-date-badge\">(.*?)<a href=\"(.*?)\">").FindAllStringSubmatch(html, -1)
	for _, str := range strs {
		url2 := url + str[2] + "?hl=ru"
		html2, err := s.HTMLGetter.GetHTML(url2)
		if err != nil {
			continue
		}
		html2 = fix(html2)
		datetime := regexp.MustCompile("datetime=\""+timeNow+"T").FindAllString(html2, -1)
		if datetime == nil {
			continue
		}
		title := regexp.MustCompile("<div class=\"news-header\" itemprop=\"headline\">(.*?)<strong>(.*?)</strong>").FindAllStringSubmatch(html2, -1)
		subtitle := regexp.MustCompile("style=\"max-width: 700px; display: block; width: 100%;\" class=\"mb-2\"/>(.*?)</article>").FindAllStringSubmatch(html2, -1)
		stock := regexp.MustCompile("<div class=\"ticker grey d-none d-sm-block\"> (.*?) </div>").FindAllStringSubmatch(html2, -1)

		answer = append(answer, hotnews.WebNews{
			From:     hotnews.Terminal,
			URL:      url2,
			Title:    title[0][2],
			SubTitle: fix(subtitle[0][1]),
			Stocks:   []hotnews.Stock{{stock[0][1], hotnews.URLTBankStock + stock[0][1]}},
			Time:     timeNow,
		})

	}
	return answer
}

func fix(html string) string {
	html = strings.Replace(html, "  ", "", strings.Count(html, "  "))
	html = strings.Replace(html, "\t", "", strings.Count(html, "\t"))
	html = strings.Replace(html, "\n", "", strings.Count(html, "\n"))
	html = strings.Replace(html, "\r", "", strings.Count(html, "\r"))
	return html
}
