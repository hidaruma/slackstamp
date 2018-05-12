package spreadsheet

import (
	"google.golang.org/api/sheets/v4"
	"golang.org/x/oauth2/google"
	"net/http"
	"io/ioutil"
	"context"
	"path/filepath"
	"fmt"
	"golang.org/x/oauth2"
	"encoding/json"
	"os"
	"log"
)

type Emoji2Stamp  map[string]string

func getClient(credential string, tokFile string) (*http.Client, error) {
	cred :=filepath.ToSlash(credential)

	b, err := ioutil.ReadFile(cred)
	if err != nil {
		b = []byte(cred)
	}
	conf, err := google.ConfigFromJSON(b, "https:/www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return nil, err
	}
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(conf)
		saveToken(tokFile, tok)
	}

	return conf.Client(context.Background(), tok), nil
}

func saveToken(tokFile string, tok *oauth2.Token) {
	fp := filepath.ToSlash(tokFile)
	f, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
		json.NewEncoder(f).Encode(tok)
	}
}

func getTokenFromWeb(conf *oauth2.Config) *oauth2.Token {
	authURL := conf.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println("Go to the following link in your browser then type the auth code: %v\n", authURL)
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		os.Exit(1)
	}
	tok, err := conf.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		fmt.Println("Unable to retrieve token from web")
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	fp := filepath.ToSlash(file)
	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func GetSheet(sheetID string, credFP string, tokFile string) (*sheets.SpreadsheetsService, error) {
	client, err := getClient(credFP, tokFile)
	if err != nil {
		return nil, err
	}
	fmt.Println(client)
	s, err := sheets.New(client)
	if err != nil {
		return nil, err
	}
	return s.Spreadsheets, nil
}


func (e2s Emoji2Stamp) SetMapping(ss *sheets.SpreadsheetsService, sheetID string, sheetName string) error {
	//	readRange := "Class " + sheetName + "!A2:B"
	resp, err := ss.Values.Get(sheetID, "Data!A2:B").Do()
	if err != nil {
		return err
	}
	for _, row := range resp.Values {
		emoji := row[0].(string)
		url := row[1].(string)
	Emoji2Stamp(e2s)[emoji] = url
	}
	return nil
}
