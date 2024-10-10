package withbrowser

import "github.com/playwright-community/playwright-go"

type WithBrowser struct {
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
	return h
}

func (h WithBrowser) GetHTML(url string) (string, error) {
	pl := &playwright.Playwright{}
	for {
		pl1, err := playwright.Run()
		if err == nil && pl1 != nil {
			pl = pl1
			break
		}
	}

	browser, err := pl.Chromium.Launch()
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

	pl.Stop()
	return string(bytes), err
}
