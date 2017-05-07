package main

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"vkAudioBot/vk"
)

const (
	START_MSG = "Привет, человек :) Ты наверняка тот, кого задели жадные правообладатели vk." +
		" Позволь, я помогу тебе скачать твои аудио в замечательное место - Telegram..." +
		" Знаешь, я уверен, что сюда эти ужасные люди не дойдут. Ведь Дуров не позволит же, правда?\n" +
		"Кстати, чтобы экспортировать все свои треки, напиши /getAll твой_vk_id. " +
		"Например: /getAll 151665536\n\n[ВНИМАНИЕ]: ВАШИ АУДИОЗАПИСИ ДОЛЖНЫ БЫТЬ ОТКРЫТЫМИ!\n\n" +
		"ID - это набор чисел, который можно узнать так: http://imgur.com/yO2JeNu"
)

func init() {
	formatDate := time.Now().Local().Format("2006-01-02")
	formatTime := time.Now().Local().Format("15:00:00")

	f, err := os.Create("log" + formatDate + ".txt")

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	f.Write([]byte(formatTime))

	log.SetOutput(f)
}

func main() {
	bannedUsers := hashset.New()

	bot, err := tgbotapi.NewBotAPI("300531701:AAE93B5-9jUozeru6w-oLGTGW1raIc75hnE") //314564211:AAH59sKgMcht-F_sVevp3jGXLo9j2VELRqg CASLBOT 300531701:AAE93B5-9jUozeru6w-oLGTGW1raIc75hnE REALBOT
	if err != nil {
		panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s %d] %s", update.Message.From.UserName, update.Message.From.ID, update.Message.Text)

		switch {
		case bannedUsers.Contains(update.Message.Chat.ID):
			{
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Извини, на твой аккаунт временно наложены ограничения."+
					" Если есть какие-то вопросы, пиши сюда - @Flerry")
				msg.ReplyToMessageID = update.Message.MessageID

				bot.Send(msg)
			}
		case update.Message.Text == "/start":
			{
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, START_MSG)
				msg.ReplyToMessageID = update.Message.MessageID

				bot.Send(msg)
			}

		case strings.Contains(update.Message.Text, "/getAll"):
			{
				re := regexp.MustCompile("[0-9]+")
				if re.FindAllString(update.Message.Text, -1) == nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Извините, но скорее всего Вы не ввели свой ID, либо ввели буквенный ID, что не прокатит ^_^.  Мне кажется, стоит "+
						"ввести свой числовой ID... \nНайти его можно так: http://imgur.com/yO2JeNu")
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я уже ушел за твоими аудио... Если я долго не появляюсь - все в порядке, просто у меня очередь :)")
					bot.Send(msg)

					go vk.GetAudio(false, strings.Replace(update.Message.Text, "/getAll ", "", -1), update.Message, bot)
				}
			}

		case strings.Contains(update.Message.Text, "/getFromEnd"):
			{
				re := regexp.MustCompile("[0-9]+")
				if re.FindAllString(update.Message.Text, -1) == nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Извините, но скорее всего Вы не ввели свой ID, либо ввели буквенный ID, что не прокатит ^_^.  Мне кажется, стоит "+
						"ввести свой числовой ID... \nНайти его можно так: http://imgur.com/yO2JeNu")
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я уже ушел за твоими аудио... Если я долго не появляюсь - все в порядке, просто у меня очередь :)")
					bot.Send(msg)

					go vk.GetAudio(true, strings.Replace(update.Message.Text, "/getFromEnd ", "", -1), update.Message, bot)
				}
			}

		case update.Message.Text == "/help":
			{
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ну, значит, пока я умею это:\n"+
					"/getAll - скачать все аудио, начиная с самых новых\n"+
					"Пример: /getAll 123252561\n\n"+
					"/getFromEnd - скачать все аудио, начиная с самых старых\n"+
					"Пример: - /getFromEnd 123252561")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		case strings.Contains(update.Message.Text, "/ban"):
			{
				banId := strings.Replace(update.Message.Text, "/ban ", "", -1)

				bannedUsers.Add(strconv.ParseInt(banId, 10, 64))

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Забанил :)")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		case update.Message.Text == "/clearBanList" && update.Message.Chat.ID == 203110206:
			{
				bannedUsers.Clear()
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Отчистил :)")
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		default:
			{
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Хм, даже не знаю, что тебе ответить."+
					" Если хочешь, чтобы я научился отвечать на такие команды, напиши моему создателю - @Flerry. Можно написать еще отзыв: https://storebot.me/bot/vksoundbot")
				msg.ReplyToMessageID = update.Message.MessageID

				bot.Send(msg)
			}
		}
	}
}
