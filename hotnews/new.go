package hotnews

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
