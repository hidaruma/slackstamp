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

type UserProfile struct {
	Ok bool `json:"ok"`
	Profile Profile `json:"profile,omitempty"`
	Error string `json:"error,omitempty"`
}

type Profile struct {
	Title string `json:"title"`
	Phone string `json:"phone"`
	Skype string `json:"skype"`
	RealName string `json:"real_name"`
	RealNameNormalized string `json:"real_name_normalized"`
	FirstName string `json:"first_name,omitempty"`
	LastName string `json:"last_name,omitempty"`
	DisplayName string `json:"display_name"`
	DisplayNameNormalized string `json:"display_name_normalized"`
	Fields []Field `json:"fields,omitempty"`
	StatusText string `json:"status_text"`
	StatusEmoji string `json:"status_emoji"`
	StatusExpiration  int `json:"status_expiration"`
	AvatarHash string `json:"avatar_hash"`
	Email string `json: "email"`
	Image24 string `json:"image_24"`
	Image32 string `json:"image_32"`
	Image48 string `json:"image_48"`
	Image72 string `json:"image_72"`
	Image192 string `json:"image_192"`
	Image512 string `json:"image_512"`
	ImageOriginal string `json:"image_original,omitempty"`
	StatusTextCanonical string `json:"status_text_canonical"`
	Team string `json:"team,omitempty"`
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

type RemoveJson struct {
	Ok bool `json:"ok"`
	Error string `json:"error,omitempty"`
	Channel string `json:"channel,omitempty"`
	Ts string `json:"ts,omitempty"`
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
	rm := RemoveJson{}
	if err = json.Unmarshal(rmJson, &rm); err != nil {
		fmt.Println("Json Unmarshal Error")
		return err
	}
	if !rm.Ok {
		return errors.New(rm.Error)
	}
	return nil
}

type PostMessageParameters struct {
	Channel string `json:"channel"`
	UserName string `json:"user_name"`
	IconURL string `json:"icon_url"`
	Attachments []Attachment `json:"attachments"`
	AsUser bool `json:"as_user"`
	User string `json:"user"`
}

type Attachment struct {
	Fallback string `json:"fallback,omitempty"`
	Color string `json:"color,omitempty"`
	AuthorName string `json:"author_name,omitempty"`
	AuthorLink string `json:"author_link,omitempty"`
	AuthorIcon string `json:"author_icon,omitempty"`
	Title string `json:"title,omitempty"`
	TitleLink string `json:"title_link,omitempty"`
	Fields []Field `json:"fields,omitempty"`
	Pretext string `json:"pretext,omitempty"`
	Text string `json:"text,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
	ThumbURL string `json:"thumb_url,omitempty"`
	Footer string `json:"footer,omitempty"`
	FooterIcon string `json:"footer_icon,omitempty"`
	Ts float64 `json:"ts,omitempty"`
}

type Field struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short bool `json:"short,omitempty"`
}

func EncodeStamp(sm *SlackMessage, st string, stampURL string) ([]byte, error) {
	stampURLDate := addDateString(stampURL)

	var ats []Attachment
	at := Attachment{
		Text: "",
		ImageURL: stampURLDate,
	}
	ats = append(ats, at)
	iconURL, err := getUserIcon(sm.UserID, st)
	if err != nil {
		fmt.Println("Error get user icon")
		return nil, err
	}
	pmp := PostMessageParameters{
	Channel: sm.ChannelID,
	UserName: sm.UserName,
	IconURL: iconURL,
	Attachments: ats,
	AsUser: true,
	User: sm.UserID,
	}
	fmt.Printf("%#v\n", pmp)
	res, err := json.Marshal(pmp)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(res))
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
	vals.Add("user", userID)
	req, err := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		return "", err
	}
	req.URL.RawQuery = vals.Encode()
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	userJson, err := ioutil.ReadAll(resp.Body)
	var up UserProfile
	fmt.Println(string(userJson))
	if err := json.Unmarshal(userJson, &up); err != nil {
		return "", err
	}
	iconURL := up.Profile.Image72
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
