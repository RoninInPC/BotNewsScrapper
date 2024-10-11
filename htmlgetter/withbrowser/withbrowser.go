package withbrowser

import (
	"github.com/playwright-community/playwright-go"
	"os/exec"
	"sync/atomic"
)

type WithBrowser struct {
}

var (
	isInstalled *atomic.Bool = nil
)

func Init() WithBrowser {
	if isInstalled == nil {
		isInstalled = &atomic.Bool{}
		isInstalled.Store(false)
		_ = playwright.Install()
		isInstalled.Store(true)
	} else {
		if !isInstalled.Load() {
			_ = playwright.Install()
			isInstalled.Store(true)
		}
	}
	h := WithBrowser{}
	return h
}

func (h WithBrowser) ReInstall() {
	if isInstalled.Load() == true {
		isInstalled.Store(false)
		exec.Command("npx", "playwright", "uninstall").Run()
		exec.Command("npx", "playwright", "install", "--with-deps").Run()
		isInstalled.Store(true)
	}
}

func (h WithBrowser) GetHTML(url string) (string, error) {
	if !isInstalled.Load() {
		return "", nil
	}
	pl, err := playwright.Run()
	if err != nil {
		h.ReInstall()
		return "", err
	}

	browser, err := pl.Firefox.Launch()
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
