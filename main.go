package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

type config struct {
	LineBotChannelToken  string
	LineBotChannelSecret string
	Port                 string
}

type gateway struct {
	botClient *linebot.Client
	config
}

func (g *gateway) loadConfig() {
	g.LineBotChannelToken = os.Getenv("LINE_BOT_CHANNEL_TOKEN")
	g.LineBotChannelSecret = os.Getenv("LINE_BOT_CHANNEL_SECRET")
	g.Port = os.Getenv("PORT")
}

func (g *gateway) newClient() {
	g.loadConfig()
	g.botClient, _ = linebot.New(g.LineBotChannelSecret, g.LineBotChannelToken)
}

func (g gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch r.URL.Path {
	case "/callback":
		events, err := g.botClient.ParseRequest(r)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				if _, err = g.botClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(time.Now().String())).Do(); err != nil {
					log.Print(err)
				}
			}
		}
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func main() {
	g := gateway{}
	g.newClient()
	if err := http.ListenAndServe(":"+g.Port, g); err != nil {
		log.Fatal(err)
	}
}
