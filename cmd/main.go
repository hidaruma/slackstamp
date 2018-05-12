package main

import (
	"net/http"
	"github.com/hidaruma/slackstamp/webhook"
	"github.com/hidaruma/slackstamp/spreadsheet"
	"fmt"
	"github.com/hidaruma/slackstamp/config"
	"os"
)

func main() {
	conf, err := config.LoadToml("conf.toml")
	if err != nil {
		fmt.Println("No config")
		os.Exit(1)
	}
	fmt.Println(conf)
	sheet, err := spreadsheet.GetSheet(conf.SpreadSheet.ID, conf.SpreadSheet.Secret, conf.SpreadSheet.Token)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var mapping spreadsheet.Emoji2Stamp
	err = mapping.SetMapping(sheet, conf.SpreadSheet.ID, conf.SpreadSheet.Name)
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
			webhook.RemoveEmoji(sm, conf.Slack.Token)
			res, err := webhook.EncodeStamp(sm, conf.Slack.Token, stampURL)
			if err != nil {
				fmt.Println(err)
			}
			w.Write(res)
		}

	}

	})
	err = http.ListenAndServe(conf.Server.Port, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}