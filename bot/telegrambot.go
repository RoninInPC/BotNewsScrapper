package bot

import (
	"BotNewsScrapper/channelsstorage"
	"BotNewsScrapper/channelsstorage/redischannels"
	"BotNewsScrapper/hotnews"
	"BotNewsScrapper/htmlgetter/withbrowser"
	"BotNewsScrapper/makeimagefromweb/warmmap"
	"BotNewsScrapper/newsstorage"
	"BotNewsScrapper/newsstorage/redisstorage"
	"BotNewsScrapper/scrapper"
	"BotNewsScrapper/scrapper/scrapperterminal"
	"BotNewsScrapper/sender"
	"github.com/and3rson/telemux/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron"
	"gopkg.in/ini.v1"
	"log"
	"time"
)

type TelegramBot struct {
	ChannelNews      chan hotnews.WebNews
	TelegramChannels channelsstorage.ChannelsStorage
	Scrappers        []scrapper.Scrapper[hotnews.WebNews]
	NewsStorage      newsstorage.NewsStorage[hotnews.WebNews]
	Senders          []sender.Sender

	Commands Commands
	BotApi   *tgbotapi.BotAPI
	AdminId  int64
}

func InitBot(fileConfig string) TelegramBot {
	inidata, err1 := ini.Load(fileConfig)
	if err1 != nil {
		panic(err1)
	}
	token := inidata.Section("telegram_bot").Key("token").String()

	secRedisNewsStorage := inidata.Section("redis_storage")
	secRedisChannelStorage := inidata.Section("redis_channels")

	tb := TelegramBot{}

	tb.ChannelNews = make(chan hotnews.WebNews)
	tb.AdminId, _ = inidata.Section("telegram_bot").Key("admin").Int64()
	tb.Commands = make(Commands, 0)

	tb.Scrappers = []scrapper.Scrapper[hotnews.WebNews]{
		//scrapperbks.ScrapperBKS{HTMLGetter: withbrowser.Init()},
		//scrapperfinam.ScrapperFinam{HTMLGetter: withbrowser.Init()},
		//scrappertbank.ScrapperTBank{HTMLGetter: simple.Simple{}},
		scrapperterminal.ScrapperTerminal{HTMLGetter: withbrowser.Init()},
	}
	db, _ := secRedisNewsStorage.Key("db").Int()
	tb.NewsStorage = redisstorage.Init[hotnews.WebNews](secRedisNewsStorage.Key("addr").String(),
		secRedisNewsStorage.Key("password").String(), db)

	db, _ = secRedisChannelStorage.Key("db").Int()
	tb.TelegramChannels = redischannels.Init(secRedisChannelStorage.Key("addr").String(),
		secRedisChannelStorage.Key("password").String(), db)

	b, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}
	tb.BotApi = b

	loc, _ := time.LoadLocation("Europe/Moscow")
	cronMoscow := cron.NewWithLocation(loc)

	tb.Senders = []sender.Sender{
		TelegramSenderImage{
			Telegram:         &tb,
			Cron:             cronMoscow,
			MakeImageFromWeb: warmmap.Init(),
			CronSetup:        "0 43 10 * * 1-5",
			Emoji:            "🗺️",
			Caption:          "Тепловая карта открытия биржи."},
		TelegramSenderImage{
			Telegram:         &tb,
			Cron:             cronMoscow,
			MakeImageFromWeb: warmmap.Init(),
			CronSetup:        "0 43 18 * * 1-5",
			Emoji:            "🗺️",
			Caption:          "Тепловая карты закрытия биржи."},
	}
	return tb
}

func (telegramBot TelegramBot) InitBotMenu() {
	var sliceArr []tgbotapi.BotCommand
	for _, action := range telegramBot.Commands {
		if len(action.Description) > 0 {
			sliceArr = append(sliceArr, tgbotapi.BotCommand{
				Command:     action.Name,
				Description: action.Description,
			})
		}
	}
	cmdCfg := tgbotapi.NewSetMyCommands(
		sliceArr...,
	)
	_, _ = telegramBot.BotApi.Send(cmdCfg)
}

func (telegramBot TelegramBot) GetUpdates(timeOut int) tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = timeOut
	return telegramBot.BotApi.GetUpdatesChan(u)
}

func (t *TelegramBot) Work(duration time.Duration) {
	for _, s := range t.Scrappers {
		s.Scrape(t.ChannelNews, "", duration)
	}
	for _, s := range t.Senders {
		s.Work()
	}
	loc, _ := time.LoadLocation("Europe/Moscow")
	go func() {
		for {
			t.NewsStorage.Free()
			for news := range t.ChannelNews {

				if time.Now().In(loc).Format("15:04:05") >= "20:00:00" {
					continue
				}
				if t.TelegramChannels.Size() == 0 {
					continue
				}
				if !t.NewsStorage.Add(news) {
					continue
				}
				log.Println("added:", news.Title, " ")

				for _, channelId := range t.TelegramChannels.GetChatsId() {
					subTitle := "\n\n"
					if len(news.SubTitle) > 0 {
						subTitle = "\n" + news.SubTitle + "\n\n"
					}
					msg := tgbotapi.NewMessage(channelId,
						//"#"+news.From+"\n\n"+
						news.MakeTags()+
							"\n*"+news.Title+"*"+subTitle+
							"@"+t.BotApi.Self.UserName)
					msg.ReplyMarkup = news.MakeButtons()
					msg.ParseMode = tgbotapi.ModeMarkdown

					t.BotApi.Send(msg)
				}
			}
			time.Sleep(duration)
		}
	}()
}

func (telegramBot TelegramBot) Dispatch() {
	mux := telemux.NewMux()

	for _, command := range telegramBot.Commands {
		mux.AddHandler(telemux.NewHandler(command.Filter, func(u *telemux.Update) {
			command.Action.Action(u)
		}))
	}

	for update := range telegramBot.GetUpdates(40) {
		mux.Dispatch(telegramBot.BotApi, update)
	}
}
