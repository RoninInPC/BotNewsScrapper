package bot

import (
	"BotNewsScrapper/makeimagefromweb"
	"crypto/md5"
	"encoding/hex"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron"
	"os"
	"time"
)

type TelegramSenderImage struct {
	Telegram         *TelegramBot
	Cron             *cron.Cron
	MakeImageFromWeb makeimagefromweb.MakeImageFromWeb
	CronSetup        string
}

func (t TelegramSenderImage) Work() {
	if t.Cron == nil {
		t.Cron = cron.New()
	}

	t.Cron.AddFunc(t.CronSetup, func() {
		name, _, err := t.MakeImageFromWeb.Get("")
		if err == nil {
			file, err := os.Open(name)
			if err == nil {
				reader := tgbotapi.FileReader{Name: GetMD5Hash(name + time.Now().String()), Reader: file}
				for _, id := range t.Telegram.TelegramChannels.GetChatsId() {
					photo := tgbotapi.NewPhoto(id, reader)
					t.Telegram.BotApi.Send(photo)
				}
			}
		}
	})
	go t.Cron.Start()
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
