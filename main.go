package main

import (
	"BotNewsScrapper/bot"
	"github.com/and3rson/telemux/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func main() {
	tbbot := bot.InitBot("config.ini")

	tbbot.Commands = append(tbbot.Commands, bot.Command{
		"help",
		"помощь",
		func(u *telemux.Update) bool {
			return bot.FilterDefault(u, "help")
		},
		&bot.SimpleActionStruct{
			func(telegramBot bot.TelegramBot, u *telemux.Update) {
				telegramBot.BotApi.Send(tgbotapi.NewMessage(u.FromChat().ID, "Здравствуйте, я новостной бот, я беру новости с https://bcs.ru/news, https://www.finam.ru/publications/section/market/, https://www.tbank.ru, https://www.tbank.ru/about/news/, https://www.tbank.ru/invest/social/profile/Interfax/"))
			},
			tbbot}})
	tbbot.Commands = append(tbbot.Commands, bot.Command{
		"start",
		"Инициализация группы/пользователя",
		func(u *telemux.Update) bool {
			return bot.FilterDefault(u, "start")
		},
		&bot.SimpleActionStruct{
			func(telegramBot bot.TelegramBot, u *telemux.Update) {
				tbbot.TelegramChannels.Add(u.FromChat().ID)
			},
			tbbot}})
	tbbot.InitBotMenu()
	go tbbot.Work(time.Second * 10)
	tbbot.Dispatch()
}
