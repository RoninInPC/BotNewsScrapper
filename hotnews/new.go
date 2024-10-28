package hotnews

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

type News interface {
	GetNews() string
}

const (
	BKS      = "BKS"
	Finam    = "Финам"
	TBank    = "T_Bank"
	Terminal = "BlackTerminal"
)

const (
	URLTBankStock = "https://www.tbank.ru/invest/stocks/"
)

type Stock struct {
	Stock string
	URL   string
}

type WebNews struct {
	From     string
	URL      string
	Title    string
	SubTitle string
	Stocks   []Stock
	Time     string
}

func (n WebNews) GetNews() string {
	return n.Title
}

func (n WebNews) MakeTags() string {
	answer := ""
	for _, stock := range n.Stocks {
		answer += "#" + strings.ToUpper(stock.Stock) + " "
	}
	return answer + "\n"
}

func (n WebNews) MakeButtons() tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, stock := range n.Stocks {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL(
				strings.ToUpper(stock.Stock), stock.URL)))
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
