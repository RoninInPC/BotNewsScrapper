package main

import (
	"BotNewsScrapper/bot"
	"github.com/and3rson/telemux/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"time"
)

func main() {
	tbbot := bot.InitBot("/etc/project/config.ini")

	tbbot.Commands = append(tbbot.Commands, bot.Command{
		"help",
		"помощь",
		func(u *telemux.Update) bool {
			return bot.FilterDefault(u, "help")
		},
		&bot.SimpleActionStruct{
			func(telegramBot bot.TelegramBot, u *telemux.Update) {
				telegramBot.BotApi.Send(tgbotapi.NewMessage(u.FromChat().ID, "Здравствуйте, я новостной бот, я беру новости с https://www.finam.ru/publications/section/market/"))
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
				if u.Message.From.ID == telegramBot.AdminId {
					tbbot.TelegramChannels.Add(u.FromChat().ID)
				}
			},
			tbbot}})

	/*tbbot.Commands = append(tbbot.Commands, bot.Command{
		"add",
		"Добавление канала/группы для постинга (/add id)",
		func(u *telemux.Update) bool {
			return bot.FilterDefault(u, "add")
		},
		&bot.SimpleActionStruct{
			func(telegramBot bot.TelegramBot, u *telemux.Update) {
				if u.Message.From.ID == telegramBot.AdminId {
					id, _ := strconv.ParseInt(strings.Split(u.Message.Text, " ")[1], 10, 64)
					tbbot.TelegramChannels.Add(id)
				}
			},
			tbbot}})
	tbbot.Commands = append(tbbot.Commands, bot.Command{
		"delete",
		"Удаление канала/группы для постинга (/delete id)",
		func(u *telemux.Update) bool {
			return bot.FilterDefault(u, "delete")
		},
		&bot.SimpleActionStruct{
			func(telegramBot bot.TelegramBot, u *telemux.Update) {
				if u.Message.From.ID == telegramBot.AdminId {
					id, _ := strconv.ParseInt(strings.Split(u.Message.Text, " ")[1], 10, 64)
					tbbot.TelegramChannels.Delete(id)
				}
			},
			tbbot}})*/

	tbbot.InitBotMenu()
	go tbbot.Work(time.Second * 30)
	tbbot.Dispatch()
}
