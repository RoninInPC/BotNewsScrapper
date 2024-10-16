package hotnews

type News interface {
	GetNews() string
}

const (
	BKS   = "BKS"
	Finam = "Финам"
	TBank = "T_Bank"
)

type WebNews struct {
	From     string
	URL      string
	Title    string
	SubTitle string
	Time     string
}

func (n WebNews) GetNews() string {
	return n.Title
}
