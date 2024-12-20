package scrapperfinam

import (
	"BotNewsScrapper/hotnews"
	"BotNewsScrapper/htmlgetter"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	BaseURL = "https://www.finam.ru/publications/section/market/"
)

type ScrapperFinam struct {
	HTMLGetter htmlgetter.HTMLGetter
}

func (s ScrapperFinam) Scrape(
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
				log.Println("Finam HTML err ", err.Error())
				time.Sleep(duration)
				continue
			}
			for _, n := range s.AnalysisHTML(html, url, time.Now().Format("2006-01-02")) {
				channel <- n
			}
			time.Sleep(duration)
		}
	}()
}

func (s ScrapperFinam) AnalysisHTML(html string, u string, timeNow string) []hotnews.WebNews {
	answer := make([]hotnews.WebNews, 0)
	html = strings.Replace(html, "\t", "", strings.Count(html, "\t"))
	html = strings.Replace(html, "\n", "", strings.Count(html, "\n"))
	html = strings.Replace(html, "\r", "", strings.Count(html, "\r"))
	strs := regexp.MustCompile(`<a href="([^"]*)" data-chp-url="([^"]*)" class="cl-blue font-l bold"([^<]*)>*</a>`).FindAllStringSubmatch(html, -1)
	strs1 := regexp.MustCompile(`<p class="font-s cl-black">(.*?)</p>(.*?)</div>`).FindAllStringSubmatch(html, -1)
	if strs1 == nil {
		return answer
	}
	for i, str := range strs {
		if len(str) < 4 {
			continue
		}
		str1 := regexp.MustCompile("([^>]*)").FindAllString(str[3], -1)

		parsed, _ := url.Parse(u)
		if len(str1) < 2 {
			continue
		}

		time := regexp.MustCompile(`\d\d\d\d\d\d\d\d`).FindAllStringSubmatch(str[1], -1)

		timeNew := time[0][0][:4] + "-" + time[0][0][4:6] + "-" + time[0][0][6:8]

		if timeNow != timeNew {
			continue
		}
		stocks := make([]hotnews.Stock, 0)

		for _, s1 := range regexp.MustCompile("href=\"(.*?)\"").FindAllStringSubmatch(strs1[i][2], -1) {
			splitted := strings.Split(s1[1], "/")
			stock := splitted[len(splitted)-2]
			stocks = append(stocks, hotnews.Stock{Stock: stock, URL: parsed.Scheme + "://" + parsed.Host + s1[1]})
		}

		answer = append(answer, hotnews.WebNews{
			From:     hotnews.Finam,
			Title:    FixTitle(str1[1]),
			SubTitle: FixTitle(strs1[i][1]),
			Stocks:   stocks,
			URL:      parsed.Scheme + "://" + parsed.Host + str[1],
			Time:     timeNew})
	}
	return answer
}

func FixTitle(str string) string {
	str = strings.Replace(str, "&quot;", "\"", strings.Count(str, "&quot;"))
	str = strings.Replace(str, "&#171;", "\"", strings.Count(str, "&#171;"))
	str = strings.Replace(str, "&#187;", "\"", strings.Count(str, "&#187;"))
	return str
}
