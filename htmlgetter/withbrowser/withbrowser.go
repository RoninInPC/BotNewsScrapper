package withbrowser

import (
	"github.com/playwright-community/playwright-go"
	"log"
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
		ans, _ := exec.Command("npx", "playwright", "uninstall").Output()
		log.Println(ans)
		ans, _ = exec.Command("npx", "playwright", "install", "--with-deps").Output()
		log.Println(ans)
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
		h.ReInstall()
		return "", err
	}

	defer browser.Close()
	page, err := browser.NewPage()

	if err != nil {
		h.ReInstall()
		return "", err
	}
	defer page.Close()

	page.SetDefaultTimeout(800000)
	response, err := page.Goto(url)
	if err != nil {
		h.ReInstall()
		return "", err
	}

	bytes, err := response.Body()

	pl.Stop()
	return string(bytes), err
}

func (h WithBrowser) GetScreenshot(url string) ([]byte, error) {
	pl, err := playwright.Run()
	if err != nil {
		h.ReInstall()
		return nil, err
	}

	browser, err := pl.Firefox.Launch()
	if err != nil {
		h.ReInstall()
		return nil, err
	}
	defer browser.Close()

	page, err := browser.NewPage()

	if err != nil {
		h.ReInstall()
		return nil, err
	}
	defer page.Close()
	page.SetDefaultTimeout(800000)
	_, err = page.Goto(url)

	if err != nil {
		h.ReInstall()
		return nil, err
	}

	screen, err := page.Screenshot()
	pl.Stop()
	return screen, err
}
