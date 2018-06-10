package webhook

import (
	"strconv"
	"time"
	"strings"
	"net/url"
	"fmt"
	"errors"
	"encoding/json"
	"net/http"
	"io/ioutil"
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

	ServiceID string `json:"service_id"`

}

type userProfile struct {
	ok bool `json:"ok"`
	profile profile `json:"profile,omitempty"`
	error string `json:"error,omitempty"`
}

type profile struct {
	title string `json:"title"`
	phone string `json:"phone"`
	skype string `json:"skype"`
	realName string `json:"real_name"`
	realNameNormalized string `json:"real_name_normalized"`
	firstName string `json:"first_name,omitempty"`
	lastName string `json:"last_name,omitempty"`
	displayName string `json:"display_name"`
	displayNameNormalized string `json:"display_name_normalized"`
	fields map[string]interface{} `json:"fields,omitempty"`
	statusText string `json:"status_text"`
	statusEmoji string `json:"status_emoji"`
	statusExpiration string `json:"status_expiration"`
	avatarHash string `json:"avatar_hash"`
	email string `json: "email"`
	image24 string `json:"image_24"`
	image32 string `json:"image_32"`
	image48 string `json:"image_48"`
	image72 string `json:"image_72"`
	image192 string `json:"image_192"`
	image512 string `json:"image_512"`
	imageOriginal string `json:"image_original,omitempty"`
	statusTextCanonical string `json:"status_text_canonical"`
	team string `json:"team,omitempty"`
}


const slackAPI = "https://slack.com/api/"

func ParseSlackMessage(r *http.Request) (*SlackMessage, error) {
	r.ParseForm()
	var token, teamID, teamDomain, channelID, channelName, timeStamp, userID, userName, text, triggerWord, serviceID string
	fmt.Printf("%#v\n", r.Form)
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
		case "service_id":
			serviceID = val[0]
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
		ServiceID: serviceID,
	}
	fmt.Printf("%#v", sm)
	return &sm, nil
}

type removeJson struct {
	ok bool `json:"ok"`
	error string `json:"error,omitempty"`
	channel string `json:"channel,omitempty"`
	ts string `json:"ts,omitempty"`
}


func RemoveEmoji(sm *SlackMessage) error {
	apiURL := slackAPI + "chat.delete"
	vals := url.Values{}
	vals.Set("token", sm.Token)
	vals.Add("channel", sm.ChannelID)
	vals.Add("ts", sm.TimeStamp)
	vals.Add("as_user", "true")

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(vals.Encode()))
	if err != nil {
		fmt.Println("Remove Emoji Err")
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	rmJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("remove message JSON parse Error\n")
	}
	rm := removeJson{}
	if err = json.Unmarshal(rmJson, &rm); err != nil {
		fmt.Println("Json Unmarshal Error")
		return err
	}
	if !rm.ok {
		return errors.New(rm.error)
	}
	return nil
}

type postMessageParameters struct {
	channel string `json:"channel"`
	userName string `json:"user_name"`
	iconURL string `json:"icon_url"`
	attachments []attachment `json:"attachments"`
	asUser bool `json:"as_user"`
	user string `json:"user"`
}

type attachment struct {
	fallback string `json:"fallback"`
	color string `json:"color"`
	authorName string `json:"author_name"`
	authorLink string `json:"author_link"`
	authorIcon string `json:"author_icon"`
	title string `json:"title"`
	titleLink string `json:"title_link"`
	fields map[string]interface{} `json:"fields,omitempty"`

	text string `json:"text"`
	imageURL string `json:"image_url"`
	thumbURL string `json:"thumb_url"`
	footer string `json:"footer"`
	footerIcon string `json:"footer_icon"`
	ts float64 `json:"ts"`
}

func EncodeStamp(sm *SlackMessage, st string, stampURL string) ([]byte, error) {
	stampURLDate := addDateString(stampURL)

	var ats []attachment
	at := attachment{
		text: "",
		imageURL: stampURLDate,
	}
	ats = append(ats, at)
	iconURL, err := getUserIcon(sm.UserID, st)
	if err != nil {
		fmt.Println("Error get user icon")
		return nil, err
	}
	pmp := postMessageParameters{
	channel: sm.ChannelID,
	userName: sm.UserName,
	iconURL: iconURL,
	attachments: ats,
	asUser: true,
	user: sm.UserID,
	}

	res, err := json.Marshal(pmp)
	if err != nil {
		return nil, err
	}

	fmt.Println("Encode Message Successfully")
	return res, nil
}


func IsEmoji(emoji string) bool {
	emojiRune := []rune(emoji)
	if emojiRune[0] == ':' && ( emojiRune[len(emoji)-1] == ':' || (emojiRune[len(emoji) -1] == '\n' && emojiRune[len(emoji) -2] == ':' )){
		return true
	} 
	return false

}

func getUserIcon(userID string, st string) (string, error) {
	apiURL := slackAPI +  "users.profile.get"
	vals := url.Values{}
	vals.Set("token", st)
	vals.Add("user_id", userID)
	req, err := http.NewRequest("GET", apiURL, strings.NewReader(vals.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		return "", err
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	userJson, err := ioutil.ReadAll(resp.Body)
	var up userProfile
	if err := json.Unmarshal(userJson, &up); err != nil {
		return "", err
	}
	iconURL := up.profile.image72
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
