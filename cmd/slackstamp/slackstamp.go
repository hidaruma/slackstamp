package main

import (
	"net/http"
	"github.com/hidaruma/slackstamp/webhook"
	"github.com/hidaruma/slackstamp/spreadsheet"
	"fmt"
	"github.com/hidaruma/slackstamp/config"
	"os"
	"google.golang.org/api/sheets/v4"
	"time"
	"flag"
)

func main() {
	port := flag.String("port", "", "Port for Heroku")
	conf, err := config.LoadToml("conf.toml")
	if err != nil {
		fmt.Println("No config")
	}
	sheetChan := make(chan *sheets.SpreadsheetsService, 1)
	errChan := make(chan error, 1)

	go func() {
		t := time.NewTicker(5 * time.Minute)
		for {
			select{
			case <-t.C:
				sheet, err := spreadsheet.GetSheet(conf.SpreadSheet.Secret, conf.SpreadSheet.Token)
				sheetChan <- sheet
				if err != nil{
				errChan <- err
				}
			}
		}
		t.Stop()
	}()

	select {
	case sheet := <- sheetChan:
		mapping, err := spreadsheet.SetMapping(sheet, conf.SpreadSheet.ID, conf.SpreadSheet.Name)
		if err != nil {
			fmt.Println("Invalid Sheet Schema or etc.")

		}
		http.HandleFunc(conf.Server.EndPoint, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			sm, err := webhook.ParseSlackMessage(r)
			if err != nil {
				fmt.Println(err)
			}

			if webhook.IsEmoji(sm.Text) {
				stampURL := webhook.GetStampURL(sm.Text, mapping)
				if stampURL == "" {
					fmt.Printf("No match stampURL")
				} else {
					webhook.RemoveEmoji(sm)
					res, err := webhook.EncodeStamp(sm, conf.Slack.Token, stampURL)
					if err != nil {
						fmt.Println(err)
					}
					w.Write(res)
				}

			}

		})

	case err := <- errChan:
		fmt.Println(err)
		os.Exit(1)
	}
	if port == "" {
		port = conf.Server.Port
	} else {
		port = fmt.Sprintf(":%s", port)
	}
	
	err = http.ListenAndServe(conf.Server.Addr + port, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}