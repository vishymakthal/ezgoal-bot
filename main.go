package main

import (
   "bytes"
   "fmt"
   "github.com/turnage/graw"
   "github.com/turnage/graw/reddit"
   "net/http"
   "os"
   "regexp"
   "time"
)

const (
   apiURL string = "https://api.ciscospark.com/v1/messages/"
)

var roomID string = os.Getenv("WEBEX_ROOM")

type goalBot struct {
   bot reddit.Bot
}

func sendWebExMessage(title string, contentLink string, commentsLink string){
   msg := fmt.Sprintf(`{"roomId": "%v", "markdown": "%v | [Link](%v) | [Comments / Mirror Links](http://reddit.com%v)"}`,roomID, title, contentLink, commentsLink)
   fmt.Printf(msg)
   var jsonStr = []byte(msg)
   req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonStr))
   req.Header.Add("Authorization", "Bearer " + os.Getenv("WEBEX_BOT_TOKEN"))
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

   if b,err := regexp.MatchString(`\[[0-9]+\]|[0-9]{1,2}'`,p.Title); b && (err==nil) {
      <-time.After(5 * time.Second)
      fmt.Printf("Title: %s | URL: %s | Permalink: %s\n", p.Title, p.URL, p.Permalink)
      sendWebExMessage(p.Title, p.URL, p.Permalink)
   }
   return nil
}

func main() {
   if bot, err := reddit.NewBotFromAgentFile("creds.txt", 0); err != nil {
      fmt.Println("Failed to create bot handle: ", err)
   } else {
      cfg := graw.Config{Subreddits: []string{"soccer"}}
      handler := &goalBot{bot: bot}
      if _, wait, err := graw.Run(handler, bot, cfg); err != nil {
         fmt.Println("Failed to start graw run: ", err)
      } else {
         fmt.Println("graw run failed: ", wait())
      }
   }
}
