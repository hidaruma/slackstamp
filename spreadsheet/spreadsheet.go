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
	"time"
)

func GetConfig(secretJson string) (*oauth2.Config, error){
	secret :=filepath.ToSlash(secretJson)

	b, err := ioutil.ReadFile(secret)
	if err != nil {
		b = []byte(secret)
	}
	conf, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return nil, err
	}
	return conf, err
}

func GetToken(ctx context.Context, conf *oauth2.Config, tokFile string) *oauth2.Token {

	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok, err = tokenFromVar(conf)
		if err != nil {

			tok = getTokenFromWeb(conf)
			saveToken(tokFile, tok)
		}
	}
	return tok
}

func RefreshToken(ctx context.Context, conf *oauth2.Config, tk *oauth2.Token) (*oauth2.Token, error) {
	ts := conf.TokenSource(ctx, tk)
	ntk := oauth2.ReuseTokenSource(tk, ts)
	return ntk.Token()
}

func GetClient(ctx context.Context, conf *oauth2.Config, token *oauth2.Token) *http.Client {
	return conf.Client(ctx, token)
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
	fmt.Printf("Go to the following link in your browser then type the auth code: %v\n", authURL)
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		os.Exit(1)
		}

	tok, err := conf.Exchange(context.Background(), authCode)
	if err != nil {
fmt.Println("Unable to retrieve token from web")
}
	return tok
}

func tokenFromVar(conf *oauth2.Config) (*oauth2.Token, error) {
	tokenRaw := os.Getenv("ACCESS_TOKEN")
	var tokJson interface{}
	err := json.Unmarshal(([]byte)(tokenRaw), &tokJson)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tokMap := tokJson.(map[string]interface{})
	at := tokMap["access_token"].(string)
	tt := tokMap["token_type"].(string)
	rt := tokMap["refresh_token"].(string)
	ep := tokMap["expires_in"].(float64)
	tok := &oauth2.Token{
		AccessToken: at,
		TokenType: tt,
		RefreshToken: rt,
		Expiry: time.Now().Add(time.Duration(ep) * time.Second),
	}
	fmt.Println(tok)
	return tok, err
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

func GetSheet(client *http.Client) (*sheets.SpreadsheetsService, error) {

	s, err := sheets.New(client)
	if err != nil {
		return nil, err
	}
	return s.Spreadsheets, nil
}


func SetMapping(ss *sheets.SpreadsheetsService, sheetID string, sheetName string) (map[string]string, error) {
	var e2s map[string]string
	e2s = map[string]string{}
	readRange := sheetName + "!A2:B"
	resp, err := ss.Values.Get(sheetID, readRange).Do()
	if err != nil {
		fmt.Println("Error get val")
		return nil, err
	}
	for _, row := range resp.Values {
		emoji := row[0].(string)
		url := row[1].(string)
		e2s[emoji] = url
		fmt.Printf("%s:%s\n", emoji, url)
	}
	return e2s, nil
}
