package scrappertbank

import (
	"BotNewsScrapper/hotnews"
	"BotNewsScrapper/htmlgetter"
	"regexp"
	"strings"
	"time"
)

const (
	stringResponseFormatMainNews = "https://cfg.tinkoff.ru/about/public/api/news/platform/v1/getArticles?lang=ru-RU&dateTimeFrom=2000-01-01T00%3A00%3A00%2B00%3A00&pageOffset=0&pageSize=100&sortBy=publishedAt"
	stringResponseNewsAnalitics  = "https://acdn.tinkoff.ru/resources?name=invest_research_resource"
	stringURLInterfax            = "https://www.tbank.ru/invest/social/profile/Interfax/"
	stringResponseNewsInterfax   = "https://www.tbank.ru/api/invest-gw/social/v1/profile/9a28aaa0-3d42-4758-b87c-d9af2e19b0fa/post?limit=30&sessionId=E2BREud5CVIEkYgo64ULDwq136jgPJfY.ix-prod-api67&appName=invest&appVersion=1.468.0&origin=web&platform=web"
	tbankURL                     = "https://www.tbank.ru"
	tbankNewsURL                 = "https://www.tbank.ru/about/news/"

	MainNews      = "MainNews"
	InterfaxNews  = "InterfaxNews"
	AnalyticsNews = "AnalyticsNews"
)

type ScrapperTBank struct {
	HTMLGetter htmlgetter.HTMLGetter
}

func appendInChannel[Info any](channel chan<- Info, infos []Info) {
	for _, i := range infos {
		channel <- i
	}
}

func (s ScrapperTBank) Scrape(channel chan<- hotnews.WebNews, url string, duration time.Duration) {
	go func() {
		for {
			timeNow := time.Now().Format("2006-01-02")

			//mainNews, _ := s.HTMLGetter.GetHTML(stringResponseFormatMainNews)
			interfaxNews, _ := s.HTMLGetter.GetHTML(stringResponseNewsInterfax)
			//analyticsNews, _ := s.HTMLGetter.GetHTML(stringResponseNewsAnalitics)

			//appendInChannel(channel, s.AnalysisHTML(mainNews, MainNews, timeNow))
			appendInChannel(channel, s.AnalysisHTML(interfaxNews, InterfaxNews, timeNow))
			//appendInChannel(channel, s.AnalysisHTML(analyticsNews, AnalyticsNews, timeNow))
			time.Sleep(duration)
		}
	}()
}

func (s ScrapperTBank) AnalysisHTML(code string, typeNews string, timeNow string) []hotnews.WebNews {
	answer := make([]hotnews.WebNews, 0)
	if typeNews == MainNews {
		strs := regexp.MustCompile(`"slug":\s*"(.*?)",\s*"title":\s*"(.*?)"(.*?)"publishedAt":\s*"(.*?)T`).FindAllStringSubmatch(code, -1)

		if strs != nil {

			for _, s := range strs {
				if len(s) < 3 {
					continue
				}
				if s[4] != timeNow {
					continue
				}
				answer = append(answer,
					hotnews.WebNews{
						From:  FixStringMarkdown(hotnews.TBank + "_Новости"),
						URL:   tbankNewsURL + s[1],
						Title: FixTitle(s[2]),
						Time:  timeNow})
			}

		}
	}

	if typeNews == InterfaxNews {
		strs := regexp.MustCompile(`"id":\s*"(.*?)",\s*"text":\s*"(.*?)\\n(.*?)"likesCount"(.*?)"owner"(.*?)"reactions"`).FindAllStringSubmatch(code, -1)
		strs1 := regexp.MustCompile(`"inserted":\s*"(.*?)T`).FindAllStringSubmatch(code, -1)
		if strs1 != nil && strs1 != nil {
			for i, s1 := range strs {
				if len(s1) < 4 {
					continue
				}
				if strs1[i][1] != timeNow {
					continue
				}

				answer = append(answer,
					hotnews.WebNews{
						From:  FixStringMarkdown("Интерфакс"),
						URL:   stringURLInterfax + s1[1],
						Title: FixTitle(s1[2]),
						Time:  timeNow})
			}
		}

	}

	if typeNews == AnalyticsNews {
		strs := regexp.MustCompile(`"date":\s*"(.*?)T(.*?)"title":\s*"(.*?)"(.*?)"itemUrl":\s*"(.*?)"`).FindAllStringSubmatch(code, -1)
		for _, s := range strs {
			if len(s) < 5 {
				continue
			}
			if s[1] != timeNow {
				continue
			}
			answer = append(answer,
				hotnews.WebNews{
					From:  FixStringMarkdown(hotnews.TBank + "_Аналитика"),
					URL:   tbankURL + s[5],
					Title: FixTitle(s[3]),
					Time:  timeNow})
		}
	}
	return answer
}

func FixStringMarkdown(str string) string {
	str = strings.Replace(str, "_", "\\_", strings.Count(str, "_"))
	str = strings.Replace(str, "`", "\\`", strings.Count(str, "`"))
	str = strings.Replace(str, "[", "\\[", strings.Count(str, "["))
	return strings.Replace(str, "*", "\\*", strings.Count(str, "*"))
}

func FixTitle(str string) string {
	str = strings.Replace(str, "\\\"", "\"", strings.Count(str, "\\\""))
	return str
}
