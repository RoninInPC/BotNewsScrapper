package hotnews

type News interface {
	GetNews() string
}

const (
	BKS   = "BKS"
	Finam = "Finam"
	TBank = "T-Bank"
)

type WebNews struct {
	From  string
	URL   string
	Title string
	Time  string
}

func (n WebNews) GetNews() string {
	return n.Title + "\n\n" + n.URL + "\n\n"
}
