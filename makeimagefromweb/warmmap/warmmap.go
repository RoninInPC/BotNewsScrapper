package warmmap

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
	x0       = 3
	y0       = 70
	x1       = 1279
	y1       = 765
	FileName = "warmmap.png"
	BaseURL  = "https://ru.tradingview.com/heatmap/stock/#%7B%22dataSource%22%3A%22MOEXRUSSIA%22%2C%22blockColor%22%3A%22change%22%2C%22blockSize%22%3A%22market_cap_basic%22%2C%22grouping%22%3A%22sector%22%7D"
)

type WarmMap struct {
}

func Init() WarmMap {
	rvi := WarmMap{}
	if !isInstalled {
		_ = playwright.Install()
		isInstalled = true
	}

	return rvi
}

func (w WarmMap) Get(url string) (string, image.Image, error) {
	if url == "" {
		url = BaseURL
	}

	pl, err := playwright.Run()

	if err != nil {
		return "", nil, err
	}

	browser, err := pl.Firefox.Launch()
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

	time.Sleep(time.Millisecond * 110)

	screen, err := page.Screenshot()

	if err != nil {
		return "", nil, err
	}

	img, _, _ := image.Decode(bytes.NewReader(screen))

	croppedImg := image.NewRGBA(image.Rect(x0, y0, x1-x0, y1-y0))
	draw.Draw(croppedImg, croppedImg.Bounds(), img, image.Pt(x0, y0), draw.Src)

	out, _ := os.Create("" + FileName)
	err = png.Encode(out, croppedImg)
	return FileName, croppedImg, err
}
