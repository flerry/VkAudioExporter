package vk

import (
	"io/ioutil"
	"strings"
	"net/http"
	"encoding/json"
	"strconv"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"github.com/cydev/zero"
	"net/url"
)

const (
	AUDIO_GET = "https://api.vk.com/method/audio.get?access_token=" +
		"b00e0925535b7f5beaf7fe6f01f3fc545d2144820b2c639a62e04c376ca55ebe408737b58962ee037b6f7&owner_id="
	AUDIO_COUNT = "https://api.vk.com/method/audio.getCount?" +
		"access_token=71b1ca0aca0fb70b4c485efa182d75854c1c3f9f27826253b1ec377a487a2afbfa42bc4d4d98bc22b362c&owner_id="
	USER_AGENT = "VKAndroidApp/4.9-1118 (Android 5.1; SDK 22; armeabi-v7a; UMI IRON; ru)"
	PROXY      = "http://46.8.29.52:65233" //"http://46.8.29.52:65233" "http://77.73.65.155:8080"
)

type AudioCount struct {
	CountAudio int `json:"response"`
}

type AudioList struct {
	Response []struct {
		Artist string `json:"artist"`
		Title  string `json:"title"`
		URL    string `json:"url,omitempty"`
	} `json:"response"`
}

func downloadFile(URL string, name string) (file tgbotapi.FileReader, response *http.Response) {
	resp, err := http.Get(URL)
	if err != nil {
		log.Printf(err.Error())
	}

	contentLength, _ := strconv.Atoi(resp.Header.Get("Content-Length"))

	file = tgbotapi.FileReader{Name: name, Reader: resp.Body, Size: int64(contentLength)}

	return file, resp
}

func exportFileToChat(URL string, name string, message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	file, response := downloadFile(URL, name)

	defer response.Body.Close()

	audioConfig := tgbotapi.NewAudioUpload(message.Chat.ID, file)

	bot.Send(audioConfig)
}

func GetAudio(fromEnd bool, id string, message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	audioList := AudioList{}

	proxyUrl, _ := url.Parse(PROXY)
	clientWithProxy := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}

	request, _ := http.NewRequest("GET", AUDIO_GET+id, nil)
	request.Header.Set("User-Agent", USER_AGENT)

	response, err := clientWithProxy.Do(request)

	if err != nil {
		log.Printf(err.Error())
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Printf(err.Error())
	}

	newJSON := strings.Replace(string(body), strconv.Itoa(getAudioCount(id))+",", "", 1)

	json.Unmarshal([]byte(newJSON), &audioList)

	if zero.IsZero(audioList) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Мне кажется, у тебя закрыты аудио-записи или их просто нет... Это все массонннннны!!!")
		msg.ReplyToMessageID = message.MessageID

		bot.Send(msg)
	}

	if fromEnd {
		if len(audioList.Response) != 0 {

			for elem := len(audioList.Response) - 1; elem != 0; elem-- {
				audioEntity := audioList.Response[elem]

				name := audioEntity.Title
				URL := audioEntity.URL

				if URL != "" {
					exportFileToChat(URL, name, message, bot)
				}
			}
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Слушай, а ты не забыл пробел между /getFromEnd и твоим id??? Или лишний может поставил... Все в этой жизни важно!")
			msg.ReplyToMessageID = message.MessageID

			bot.Send(msg)
			return
		}
	} else {
		if len(audioList.Response) != 0 {
			for elem := range audioList.Response {
				audioEntity := audioList.Response[elem]

				name := audioEntity.Title
				URL := audioEntity.URL

				if URL != "" {
					exportFileToChat(URL, name, message, bot)
				}
			}
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Слушай, а ты не забыл пробел между /getAll и твоим id??? Или лишний может поставил... Все в этой жизни важно!")
			msg.ReplyToMessageID = message.MessageID

			bot.Send(msg)

			return
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "Я закончил... Думаю, тебе понравится результат :) \nИ если это так... То мой создатель был бы рад, если бы ты оценил меня здесь: https://storebot.me/bot/vksoundbot")
	msg.ReplyToMessageID = message.MessageID

	bot.Send(msg)
}

func getAudioCount(id string) (count int) {
	audioCount := AudioCount{}
	client := &http.Client{}

	request, _ := http.NewRequest("GET", AUDIO_COUNT+id, nil)
	request.Header.Set("User-Agent", USER_AGENT)

	response, err := client.Do(request)

	if err != nil {
		log.Printf(err.Error())
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Printf(err.Error())
	}

	json.Unmarshal(body, &audioCount)

	return audioCount.CountAudio
}
