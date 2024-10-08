package withbrowser

import "github.com/playwright-community/playwright-go"

type WithBrowser struct {
	PlayWright *playwright.Playwright
}

var (
	isInstalled = false
)

func Init() WithBrowser {
	if !isInstalled {
		_ = playwright.Install()
		isInstalled = true
	}
	h := WithBrowser{}
	pl, err := playwright.Run()
	if err == nil {
		h.PlayWright = pl
	}
	return h
}

func (h WithBrowser) GetHTML(url string) (string, error) {
	browser, err := h.PlayWright.Firefox.Launch()
	if err != nil {
		return "", err
	}

	defer browser.Close()
	page, err := browser.NewPage()

	if err != nil {
		return "", err
	}
	defer page.Close()

	response, err := page.Goto(url)
	if err != nil {
		return "", err
	}

	bytes, err := response.Body()
	return string(bytes), err
}
