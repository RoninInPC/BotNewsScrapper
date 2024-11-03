package withbrowser

import (
	"github.com/playwright-community/playwright-go"
	"log"
	"os/exec"
	"sync/atomic"
	"time"
)

type WithBrowser struct {
}

var (
	isInstalled atomic.Bool
)

func Init() WithBrowser {
	h := WithBrowser{}
	isInstalled.Store(true)
	_ = playwright.Install()
	return h
}

func (h WithBrowser) ReInstall() {
	ReInstall()
}

func (h WithBrowser) GetHTML(url string) (string, error) {
	if !isInstalled.Load() {
		return "", nil
	}
	pl, err := playwright.Run()
	if err != nil {
		return "", err
	}
	defer pl.Stop()
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

	page.SetDefaultTimeout(800000)
	response, err := page.Goto(url)
	if err != nil {
		return "", err
	}

	bytes, err := response.Body()

	return string(bytes), err
}

func (h WithBrowser) GetScreenshot(url string) ([]byte, error) {
	pl, err := playwright.Run()
	if err != nil {
		return nil, err
	}
	defer pl.Stop()
	browser, err := pl.Chromium.Launch()
	if err != nil {
		return nil, err
	}
	defer browser.Close()

	page, err := browser.NewPage()

	if err != nil {
		return nil, err
	}
	defer page.Close()
	page.SetDefaultTimeout(800000)
	_, err = page.Goto(url)

	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second * 4)
	screen, err := page.Screenshot()
	return screen, err
}

func ReInstall() {
	if isInstalled.Load() == true {
		isInstalled.Store(false)
		ans, _ := exec.Command("npx", "playwright", "uninstall", "--all").Output()
		log.Println(string(ans))
		time.Sleep(time.Second * 10)
		ans, _ = exec.Command("npx", "playwright", "install", "--with-deps").Output()
		log.Println(string(ans))
		isInstalled.Store(true)
	}
}
