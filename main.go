package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

const (
	apiURL string = "https://api.ciscospark.com/v1/messages/"
)

var roomID string = os.Getenv("WEBEX_ROOM")

type goalBot struct {
	bot reddit.Bot
}

func sendWebExMessage(title string, contentLink string, commentsLink string) {
	log.Info("Posting to WebEx Room")
	msg := fmt.Sprintf(`{"roomId": "%v", "markdown": "####%v\n| [Link](%v) | [Comments / Mirror Links](http://reddit.com%v)"}`, roomID, title, contentLink, commentsLink)
	fmt.Printf(msg)
	var jsonStr = []byte(msg)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonStr))
	req.Header.Add("Authorization", "Bearer "+os.Getenv("WEBEX_BOT_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	fmt.Println("response status:", resp.Status)
}

func (r *goalBot) Post(p *reddit.Post) error {
	if b, err := regexp.MatchString(`\[[0-9]+\]|[0-9]{1,2}'`, p.Title); b && (err == nil) {
		if b, err := regexp.MatchString(`Arsenal|Manchester|Liverpool|Tottenham|Chelsea|Sheffield|Wolve|Leicester|Norwich|Brighton|Southampton`, p.Title); b && (err == nil) {
			log.Info("Found a goal to send")
			<-time.After(10 * time.Second)
			fmt.Printf("Title: %s | URL: %s | Permalink: %s\n", p.Title, p.URL, p.Permalink)
			sendWebExMessage(p.Title, p.URL, p.Permalink)
		}
	}
	return nil
}

func main() {

	log.Info("Starting bot")
	if bot, err := reddit.NewBotFromAgentFile("creds.txt", 0); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("Failed to create bot handle")
	} else {
		log.Info("Loaded Reddit credentials")
		cfg := graw.Config{Subreddits: []string{"soccer"}}
		handler := &goalBot{bot: bot}
		if _, wait, err := graw.Run(handler, bot, cfg); err != nil {
			fmt.Println("Failed to start graw run: ", err)
			log.WithFields(log.Fields{
				"error": err,
			}).Fatal("failed to start graw run")
		} else {
			log.WithFields(log.Fields{
				"error": wait(),
			}).Fatal("graw run failed")
		}
	}
}
