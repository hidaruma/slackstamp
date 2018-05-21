package webhook

import (
	"strconv"
	"time"
	"strings"
	"github.com/nlopes/slack"
	"fmt"
	"errors"
	"encoding/json"
	"net/http"
)

type SlackMessage struct {
	Token string `json:"token"`
	TeamID string `json:"team_id"`
	TeamDomain string `json:"team_domain"`
	ChannelID string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	TimeStamp string `json:"timestamp"`
	UserID string `json:"user_id"`
	UserName string `json:"user_name"`
	Text string `json:"text"`
	TriggerWord string `json:"trigger_word"`

}

func ParseSlackMessage(r *http.Request) (*SlackMessage, error) {
	r.ParseForm()
	var token, teamID, teamDomain, channelID, channelName, timeStamp, userID, userName, text, triggerWord string
	for key, val := range r.Form {
		switch key {
		case "token":
			token = val[0]
		case "team_id":
			teamID = val[0]
		case "team_domain":
			teamDomain = val[0]
		case "channel_id":
			channelID = val[0]
		case "channel_name":
			channelName = val[0]
		case "timestamp":
			timeStamp = val[0]
		case "user_id":
			userID = val[0]
		case "user_name":
			userName = val[0]
		case "text":
			text = val[0]
		case "trigger_word":
			triggerWord = val[0]
		default:
			return nil, errors.New("Invalid request\n")
		}

	}
	sm := SlackMessage{
		Token: token,
		TeamID: teamID,
		TeamDomain: teamDomain,
		ChannelID: channelID,
		ChannelName: channelName,
		TimeStamp: timeStamp,
		UserID: userID,
		UserName: userName,
		Text: text,
		TriggerWord: triggerWord,
	}

	fmt.Print(sm)
	return &sm, nil
}

func RemoveEmoji(sm *SlackMessage, slackToken string) error {
	api := slack.New(slackToken)
	_, _, err  := api.DeleteMessage(sm.ChannelID, sm.TimeStamp)
	if err != nil {
		return err
	}
	return nil
}

func EncodeStamp(sm *SlackMessage, slackToken string, stampURL string) ([]byte, error) {
	stampURLDate := addDateString(stampURL)

	var ats []slack.Attachment
	at := slack.Attachment{
		Text: "",
		ImageURL: stampURLDate,
	}
	ats = append(ats, at)
	iconURL, err := getUserIcon(sm.UserID, slackToken)
	if err != nil {
		return nil, err
	}
	pmp := slack.PostMessageParameters{
	Channel: sm.ChannelID,
	Username: sm.UserName,
	IconURL: iconURL,
	Attachments: ats,
	AsUser: true,
	User: sm.UserID,
	}

	res, err := json.Marshal(pmp)
	if err != nil {
		return nil, err
	}

	fmt.Println("Encode Message Successfully")
	return res, nil
}


func IsEmoji(emoji string) bool {
	if []rune(emoji)[0] == ':' && []rune(emoji)[len(emoji)-1] == ':' {
		return true
	} else {
		return false
	}
}

func getUserIcon(userID string, st string) (string, error) {
	api := slack.New(st)
	user, err := api.GetUserInfo(userID)
	if err != nil {
		return "", err
	}
	iconURL := user.Profile.Image72
	return iconURL, nil
}

func GetStampURL(emoji string, mapping map[string]string) string {
	imageURL := mapping[emoji]
	return imageURL
}

func addDateString(url string) string {
	now := strconv.Itoa(int(time.Now().Unix()))
	if strings.ContainsRune(url, '?') {
		return url + "&" + now
	} else {
		return url + "?" + now
	}
}
