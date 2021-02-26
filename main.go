package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sogko/go-wordpress"
	"log"
	"os"
	"time"
)

const (
	wpAPI           = "http://ibigdan.com/wp-json/wp/v2"
	tgChannel       = "@ibigdan_tg"
	messageTemplate = "<a href='https://t.me/iv?url=%s&rhash=97b4a4b92bcfca'>%s</a>"
)

func main() {
	wp := wordpress.NewClient(&wordpress.Options{
		BaseAPIURL: wpAPI,
		Username:   "",
		Password:   "",
	})

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	lastPosted := 0

	for {
		posts, _, _, err := wp.Posts().List(nil)
		if err != nil {
			log.Printf("Error receiving posts: %v", err)
			continue
		}
		for i := range posts {
			post := posts[len(posts)-i-1]
			if post.ID <= lastPosted {
				break
			}
			text := fmt.Sprintf(messageTemplate, post.Link, post.Title.Rendered)
			msg := tgbotapi.NewMessageToChannel(tgChannel, text)
			msg.ParseMode = "html"
			_, err = bot.Send(msg)
			if err != nil {
				log.Printf("Error sending to the channel: %v", err)
				continue
			}
			log.Printf("Sent to channel: %s", post.Link)
		}
		if len(posts) > 0 {
			lastPosted = posts[0].ID
		}
		time.Sleep(5 * time.Minute)
	}
}
