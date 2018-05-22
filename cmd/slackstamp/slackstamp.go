package main

import (
	"net/http"
	"github.com/hidaruma/slackstamp/webhook"
	"github.com/hidaruma/slackstamp/spreadsheet"
	"fmt"
	"github.com/hidaruma/slackstamp/config"
	"os"
	"github.com/nlopes/slack"
)

func main() {
	conf, err := config.LoadToml("conf.toml")
	if err != nil {
		fmt.Println("No config")
	}
	sheet, err := spreadsheet.GetSheet(conf.SpreadSheet.Secret, conf.SpreadSheet.Token)
	
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}


	mapping, err := spreadsheet.SetMapping(sheet, conf.SpreadSheet.ID, conf.SpreadSheet.Name)
	if err != nil {
		fmt.Println("Invalid Sheet Schema or etc.")

	}
	api := slack.New(conf.Slack.Token)
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
			webhook.RemoveEmoji(sm, api)
			res, err := webhook.EncodeStamp(sm, api, stampURL)
			if err != nil {
				fmt.Println(err)
			}
			w.Write(res)
		}

	}

	})

	http.Handle("/policy", http.StripPrefix("/policy", http.FileServer(http.Dir("../../public/policy"))))
	http.Handle("/terms", http.StripPrefix("/terms", http.FileServer(http.Dir("../../public/terms/"))))
	err = http.ListenAndServe(conf.Server.Port, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}