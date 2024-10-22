package htmlgetter

type HTMLGetter interface {
	GetHTML(url string) (string, error)
	GetScreenshot(url string) ([]byte, error)
}
