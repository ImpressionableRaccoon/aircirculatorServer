package telegram

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

var (
	ErrSendMessageFailed = errors.New("send telegram message failed")
)

type sendMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func SendMessage(token string, chatID string, msg string) (err error) {
	data := sendMessage{
		ChatID:    chatID,
		Text:      msg,
		ParseMode: "HTML",
	}

	client := resty.New()

	link := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	r, err := client.R().SetBody(data).Post(link)
	if err != nil {
		return err
	}

	if r.StatusCode() != 200 {
		return fmt.Errorf("%w: %s", ErrSendMessageFailed, string(r.Body()))
	}

	return nil
}
