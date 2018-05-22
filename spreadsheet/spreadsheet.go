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


func getClient(secretJson string, tokFile string) (*http.Client, error) {
	secret :=filepath.ToSlash(secretJson)

	b, err := ioutil.ReadFile(secret)
	if err != nil {
		b = []byte(secret)
	}
	conf, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		return nil, err
	}
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok, err = tokenFromVar(conf)
		if err != nil {

			tok = getTokenFromWeb(conf)
			saveToken(tokFile, tok)
		}
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

	tok, err := conf.Exchange(context.Background(), authCode)
	if err != nil {
fmt.Println("Unable to retrieve token from web")
}
	return tok
}

type tokenJson struct {
	accessToken string `json:"access_token"`
	expiry int `json:"expires_in"`
	refreshToken string `json:"refresh_token"`
	tokenType string `json:"token_type"`
}

func tokenFromVar(conf *oauth2.Config) (*oauth2.Token, error) {
	tokenRaw := os.Getenv("ACCESS_TOKEN")
	tokJson := new(tokenJson)
	fmt.Printf("%v\n", tokJson)
	err := json.Unmarshal(([]byte)(tokenRaw), tokJson)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(tokJson.accessToken)
	tok := &oauth2.Token{
		AccessToken:tokJson.accessToken,
		TokenType: tokJson.tokenType,
		RefreshToken: tokJson.refreshToken,
		Expiry: time.Now().Add(time.Duration(tokJson.expiry) * time.Second),
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

func GetSheet(secret string, tokFile string) (*sheets.SpreadsheetsService, error) {
	client, err := getClient(secret, tokFile)
	if err != nil {
		return nil, err
	}

	s, err := sheets.New(client)
	if err != nil {
		return nil, err
	}
	return s.Spreadsheets, nil
}


func SetMapping(ss *sheets.SpreadsheetsService, sheetID string, sheetName string) (map[string]string, error) {
	var e2s map[string]string
	e2s = map[string]string{}
	//	readRange := + sheetName + "!A2:B"
	resp, err := ss.Values.Get("1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms", "Class Data!A1:B").Do()
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
