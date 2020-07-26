package daemon

import (
	"fmt"
	"net/http"
	"net/url"
)

func sendTelegramMessage(botId string, chatId string, msg string) {
	if botId == "" || chatId == "" || msg == "" {
		return
	}

	endPoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botId)
	formData := url.Values{
		"chat_id":    {chatId},
		"parse_mode": {"html"},
		"text":       {msg},
	}
	_, err := http.PostForm(endPoint, formData)
	if err != nil {
		fmt.Println(fmt.Sprintf("send telegram message error, bot_id=%s, chat_id=%s, msg=%s, err=%s", botId, chatId, msg, err.Error()))
		return
	}
}
