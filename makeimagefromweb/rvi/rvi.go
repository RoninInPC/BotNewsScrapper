package rvi

import (
	"bytes"
	"github.com/playwright-community/playwright-go"
	"image"
	"image/draw"
	"image/png"
	"os"
	"time"
)

var (
	isInstalled = false
)

const (
	x0       = 50
	y0       = 40
	x1       = 980
	y1       = 615
	FileName = "rvi.png"
	BaseURL  = "https://ru.tradingview.com/chart/?symbol=RUS%3ARVI"
)

type RVI struct {
	PlayWright *playwright.Playwright
}

func Init() RVI {
	rvi := RVI{}
	if !isInstalled {
		_ = playwright.Install()
		isInstalled = true
	}
	pl, err := playwright.Run()
	if err == nil {
		rvi.PlayWright = pl
	}
	return rvi
}

func (r RVI) Get(url string) (string, image.Image, error) {
	if url == "" {
		url = BaseURL
	}
	browser, err := r.PlayWright.Firefox.Launch()
	if err != nil {
		return "", nil, err
	}
	defer browser.Close()

	page, err := browser.NewPage()

	if err != nil {
		return "", nil, err
	}
	defer page.Close()
	_, err = page.Goto(url)

	if err != nil {
		return "", nil, err
	}

	time.Sleep(time.Millisecond * 100)

	screen, err := page.Screenshot()

	if err != nil {
		return "", nil, err
	}

	img, _, _ := image.Decode(bytes.NewReader(screen))

	croppedImg := image.NewRGBA(image.Rect(x0, y0, x1-x0, y1-y0))
	draw.Draw(croppedImg, croppedImg.Bounds(), img, image.Pt(x0, y0), draw.Src)

	out, _ := os.Create("./" + FileName)
	err = png.Encode(out, croppedImg)
	return FileName, croppedImg, err
}
