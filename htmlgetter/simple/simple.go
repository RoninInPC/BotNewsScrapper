package simple

import (
	"io"
	"net/http"
	"time"
)

type Simple struct {
}

func (s Simple) GetHTML(url string) (string, error) {
	myClient := &http.Client{Timeout: 10 * time.Second}
	r, err := myClient.Get(url)
	if err != nil {
		return "", err
	}
	bytes, _ := io.ReadAll(r.Body)
	return string(bytes), nil
}

func (s Simple) GetScreenshot(url string) ([]byte, error) {
	return []byte{}, nil
}
