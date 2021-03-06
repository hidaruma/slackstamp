package main

import (
	"net/http"
	"github.com/hidaruma/slackstamp/webhook"
	"github.com/hidaruma/slackstamp/spreadsheet"
	"fmt"
	"github.com/hidaruma/slackstamp/config"
	"os"
	"time"
	"flag"
	"context"
)

func main() {
	ctx := context.Background()
	p := flag.String("port", "", "Port for Heroku")
	port := *p
	conf, err := config.LoadToml("conf.toml")
	if err != nil {
		fmt.Println("No config")
	}

	sconf, err := spreadsheet.GetConfig(conf.SpreadSheet.Secret)
	if err != nil {
		fmt.Println("Can't get oauth2 client")
		os.Exit(1)
	}

	tok := spreadsheet.GetToken(ctx, sconf, conf.SpreadSheet.Token)
	errChan1 := make(chan error, 1)
	go func () {
		t := time.NewTicker(50 * time.Minute)
		for {
			select {
			case <-t.C:
				tok, err = spreadsheet.RefreshToken(ctx, sconf, tok)
				if err != nil {
					errChan1 <- err
				}
			}
		}
	t.Stop()
	}()

	client := spreadsheet.GetClient(ctx, sconf, tok)
	if err != nil {
		fmt.Println("Can't get http.client")
		os.Exit(1)
	}
	srv, err := spreadsheet.GetSheet(client)
	if err != nil {
		fmt.Println("Can't get sheet")
		os.Exit(1)
	}
	errChan2 := make(chan error, 1)
	emojiURL := map[string]string{}
	go func() {
		t := time.NewTicker(5 * time.Minute)
		for {
			select{
			case <-t.C:
				mapping, err := spreadsheet.SetMapping(srv, conf.SpreadSheet.ID, conf.SpreadSheet.Name)

				if err != nil{
				errChan2 <- err
				}
				emojiURL = mapping
			}
		}
		t.Stop()
	}()
	go func() {
		select {
		case <-errChan1:
			fmt.Println("Error to Get token")
			os.Exit(1)
		case <-errChan2:
			fmt.Println("Error to get Sheet")
			os.Exit(1)
		}
	}()
	http.HandleFunc(conf.Server.EndPoint, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				sm, err := webhook.ParseSlackMessage(r)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("%#v\n", sm)
				fmt.Println("Parse OK")
				if webhook.IsRmStamp(sm.Text) {
					if err = webhook.RemoveStamp(sm, conf.Slack.Token); err != nil {
						fmt.Println(err)
					}
				}

				if webhook.IsEmoji(sm.Text) {
					stampURL := webhook.GetStampURL(sm.Text, emojiURL)
					if stampURL == "" {
						fmt.Println("No match stampURL")
					} else {
						err = webhook.RemoveEmoji(sm, conf.Slack.Token)
						if err != nil {
							fmt.Println(err)
						}
						res, err := webhook.EncodeStamp(sm, conf.Slack.Token, stampURL)
						if err != nil {
							fmt.Printf("EncodeStampError:%v\n", err)
						}
						w.Write(res)
						fmt.Printf("Written:%v", res)
					}
	
				}
	
			})
	
	if port == "" {
		port = conf.Server.Port
	} else {
		port = fmt.Sprintf(":%s", port)
	}
	fmt.Println(port)
	fmt.Println(os.Getenv("PORT"))
	if err = http.ListenAndServe(port, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}