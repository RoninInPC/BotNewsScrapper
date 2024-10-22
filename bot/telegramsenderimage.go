package bot

import (
	"BotNewsScrapper/makeimagefromweb"
	"crypto/md5"
	"encoding/hex"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron"
	"log"
	"os"
	"time"
)

type TelegramSenderImage struct {
	Telegram         *TelegramBot
	Cron             *cron.Cron
	MakeImageFromWeb makeimagefromweb.MakeImageFromWeb
	CronSetup        string
	Emoji            string
	Caption          string
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
					photo.Caption = t.Emoji + "*" + t.Caption + "*"
					photo.ParseMode = tgbotapi.ModeMarkdown
					t.Telegram.BotApi.Send(photo)
				}
			}
		} else {
			log.Println("println MakeImageFromMap warmap", err.Error())
		}
	})
	go t.Cron.Start()
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
