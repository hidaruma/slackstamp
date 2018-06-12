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
	"regexp"
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


func RemoveEmoji(sm *SlackMessage, st string) error {
	apiURL := slackAPI + "chat.delete"
	vals := url.Values{}
	vals.Set("token", st)
	vals.Add("channel", sm.ChannelID)
	vals.Add("ts", sm.TimeStamp)
	vals.Add("as_user", "true")

	req, err := http.NewRequest("POST", apiURL, nil)
	if err != nil {
		fmt.Println("Remove Emoji Err")
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = vals.Encode()
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
	//Channel string `json:"channel"`
	UserName string `json:"username"`
	IconURL string `json:"icon_url"`
	Text string `json:"text"`
	Attachments []Attachment `json:"attachments"`
//	AsUser bool `json:"as_user"`
//	User string `json:"user"`
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

type MkdwnIn struct {
	Text string `json:"text,omitempty"`
	Pretext string `json:"pretext,omitempty"`
}

type AttachmentHistory struct {
	ID json.Number `json:"id"`
	Fallback string `json:"fallback,omitempty"`
	Color string `json:"color,omitempty"`
	AuthorName string `json:"author_name,omitempty"`
	AuthorSubname string `json:"author_subname,omitempty"`
	FromURL string `json:"from_url,omitempty"`
	ChannelID string `json:"channel_id,omitempty"`
	ChannelName string `json:"channel_name,omitempty"`
	IsMsgUnfURL bool `json:"is_msg_unfurl,omitempty"`
	MrkdwnIn MkdwnIn `json:"mkdwn_in,omitempty"`
	OriginalURL string `json:"original_url,omitempty"`
	MsgSubtype string `json:"msg_subtype,omitempty"`

	VideoHtml string `json:"video_html,omitempty"`
	VideoHtmlWidth json.Number `json:"video_html_width,omitempty"`
	VideoHtmlHeight json.Number `json:"video_html_height,omitempty"`

	ServiceName string `json:"service_name,omitempty"`
	ServiceURL string `json:"service_url,omitempty"`

	AuthorLink string `json:"author_link,omitempty"`
	AuthorID string `json:"author_id,omitempty"`
	AuthorIcon string `json:"author_icon,omitempty"`
	Title string `json:"title,omitempty"`
	TitleLink string `json:"title_link,omitempty"`
	Fields []Field `json:"fields,omitempty"`
	Pretext string `json:"pretext,omitempty"`
	Text string `json:"text,omitempty"`
	ImageURL string `json:"image_url,omitempty"`
	ThumbURL string `json:"thumb_url,omitempty"`
	ThumbWidth json.Number `json:"thumb_width,omitempty"`
	ThumbHeight json.Number `json:"thumb_height,omitempty"`
	Footer string `json:"footer,omitempty"`
	FooterIcon string `json:"footer_icon,omitempty"`
	Ts json.Number `json:"ts,omitempty"`
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
		ImageURL: stampURLDate,
		ThumbURL: stampURLDate,
		Fallback: "Post Stamp",
	}
	ats = append(ats, at)
	iconURL, err := getUserIcon(sm.UserID, st)
	if err != nil {
		fmt.Println("Error get user icon")
	}
	pmp := PostMessageParameters{
	//Channel: sm.ChannelID,
	Text: " ",
	UserName: sm.UserName,
	IconURL: iconURL,
	Attachments: ats,
//	AsUser: true,
//	User: sm.UserID,
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

type SlackPermalink struct {
	Ok bool `json:"ok"`
	Channel string `json:"channel,omitempty"`
	Permalink string `json:"permalink,omitempty"`
	Error string `json:"error,omitempty"`
}

func getPermalinks(chID string, ts string, st string) (string, error) {
	apiURL := slackAPI + "chat.getPermalink"
	vals := url.Values{}
	vals.Set("token", st)
	vals.Add("channel", chID)
	vals.Add("message_ts", ts)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("get slack permalink")
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = vals.Encode()
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	sp := SlackPermalink{}

	spJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(spJson, sp); err != nil {
		return "", err
	}
	fmt.Println(sp)
	return sp.Permalink, nil
}

type SlackHistory struct {
	Ok bool `json:"bool"`
	Latest string `json:"latest,omitempty"`
	Messages []Message `json:"messages,omitempty"`
	HasMore bool `json:"has_more,omitempty"`
	Error string `json:"error,omitempty"`
	IsLimited bool `json:"is_limited,omitempty"`
}

type Message struct {
	Type string `json:"type"`
	Ts string `json:"ts"`
	User string `json:"user,omitempty"`
	DisplayAsBot bool `json:"display_as_bot,omitempty"`
	File File `json:"file,omitempty"`
	Text string `json:text,omitempty`
	IsStarred bool `json:"is_starred,omitempty"`
	Reactions []Reaction `json:"reactions,omitempty"`
	ClientMsgID string `json:"client_msg_id,omitempty"`
	UserName string `json:"username,omitempty"`
	BotID string `json:"bot_id,omitempty"`
	Attachments []AttachmentHistory `json:"attachments,omitempty"`
	Subtype string `json:"subtype,omitempty"`
	Icons Icons `json:"icons,omitempty"`
	Mrkdwn bool `json:"mrkdwn,omitempty"`

}
type File struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Timestamp string `json:"timestamp"`
	Name string `json:"name,omitempty"`
	Filetype string `json:"filetype"`
	Mimetype string `json:"mimetype,omitempty"`
	PrettyType string `json:"pretty_type,omitempty"`
	Editable bool `json:"editable"`
	Mode string `json:"mode,omitempty"`
	Size json.Number `json:"size"`
	IsExternal bool `json:"is_external"`
	ExternalType string `json:"external_type,omitempty"`
	IsPublic bool `json:"is_public"`
	PublicURLShared bool `json:"public_url_shared"`
	URLPrivate string `json:"url_private,omitempty"`
	URLPrivateDownload string `json:"url_private_download,omitempty"`
	Created json.Number `json:"created"`
	Thumb64 string `json:"thumb_64,omitempty"`
	Thumb80 string `json:"thumb_80,omitempty"`
	Thumb360 string `json:"thumb_360,omitempty"`
	Thumb360H json.Number `json:"thumb_360_h,omitempty"`
	Thumb360W json.Number `json:"thumb_360_w,omitempty"`
	Thumb480 string `json:"thumb_480,omitempty"`
	Thumb480H json.Number `json:"thumb_480_h,omitempty"`
	Thumb480W json.Number `json:"thumb_480_w,omitempty"`
	Thumb720 string `json:"thumb_720,omitempty"`
	Thumb800 string `json:"thumb_800,omitempty"`
	Thumb720H json.Number `json:"thumb_720_h,omitempty"`
	Thumb720W json.Number `json:"thumb_720_w,omitempty"`
	Thumb960 string `json:"thumb_960,omitempty"`
	Thumb960H json.Number `json:"thumb_960_h,omitempty"`
	Thumb960W json.Number `json:"thumb_960_w,omitempty"`
	Thumb1024 string `json:"thumb_1024,omitempty"`
	Thumb1024H json.Number `json:"thumb_1024_h,omitempty"`
	Thumb1024W json.Number `json:"thumb_1024_w,omitempty"`
	OriginalH json.Number `json:"original_h,omitempty"`
	OriginalW json.Number `json:"original_w,omitempty"`
	ImageExifRotation json.Number `json:"image_exif_rotation,omitempty"`
	User string `json:"user,omitempty"`
	UserName string `json:"username,omitempty"`
	PermalinkPublic string `json:"permalink_public,omitempty"`
	CommentsCount json.Number `json:"comments_count,omitempty"`
	Permalink string `json:"permalink,omitempty"`
	Channles []string `json:"channels,omitempty"`
	Groups []string `json:"groups,omitempty"`
	Ims []string `json:"ims,omitempty"`
	InitialComment Comment `json:"initial_comment,omitempty"`
}


type Comment struct {
	ID string `json:"id"`
	Created json.Number `json:"created"`
	Timestamp string `json:"timestamp"`
	User string `json:"user,omitempty"`
	IsIntro bool `json:"is_intro"`
	Comment string `json:"comment"`

}

type Icons struct {
	Emoji string `json:"emoji,omitempty"`
	Image36 string `json:"image_36,omitempty"`
	Image48 string `json:"image_48,omitempty"`
	Image64 string `json:"image_64,omitempty"`
	Image72 string `json:"image_72,omitempty"`
	Image96 string `json:"image_96,omitempty"`
	Image128 string `json:"image_128,omitempty"`
}

type Reaction struct {
	Name string `json:"name"`
	Count int `json:"count"`
	Users []string `json:"users"`
}


func getHistory(channelID string, st string) (*SlackHistory, error) {
	apiURL := slackAPI + "channels.history"
	vals := url.Values{}
	vals.Set("token", st)
	vals.Add("channel", channelID)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("can't get slack history")
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = vals.Encode()
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	shJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var sh SlackHistory
	if err = json.Unmarshal(shJson, &sh); err != nil {
		return nil, err
	}
	return &sh, nil
}

func RemoveStamp(sm *SlackMessage, st string) error {
	apiURL := slackAPI + "chat.delete"
	re := regexp.MustCompile(`^rmstamp\s<(.+)>$`)
	msURL := re.FindStringSubmatch(sm.Text)
	if msURL[1] == "" {
		return errors.New("Invalid URL\n")
	}
	fmt.Println(msURL[1])
	reMs := regexp.MustCompile(`^https://.+\.slack\.com/archives/(.+)/.+$`)
	m := reMs.FindStringSubmatch(msURL[1])
	channelID := m[1]
	fmt.Println(channelID)
	sh, err := getHistory(channelID, st)
	if err != nil {
		return err
	}
	if !sh.Ok {
		return errors.New(sh.Error)
	}
	var ts string
	for _, ms := range sh.Messages {
		if ms.Type != "message" {
			continue
		}
		plink, err := getPermalinks(channelID, ms.Ts, st)
		fmt.Println(plink)
		if err != nil {
			fmt.Printf("Got error %v\n", err)
		}
		fmt.Println(plink)
		if plink == msURL[1] {
			ts = ms.Ts
		}
	}
	if ts == "" {
		return errors.New("Can't get TS\n")
	}

	vals := url.Values{}
	vals.Set("token", st)
	vals.Add("channel", channelID)
	vals.Add("ts", ts)
	vals.Add("as_user", "true")

	req, err := http.NewRequest("POST", apiURL, nil)
	if err != nil {
		fmt.Println("Remove Stamp Err")
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.URL.RawQuery = vals.Encode()
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


func IsRmStamp(word string) bool {
	re := regexp.MustCompile(`^rmstamp\s.+`)
	if re.MatchString(word) {
		return true
	}
	return false
}

func IsEmoji(emoji string) bool {
	re := regexp.MustCompile(`^:.+:$`)
	if re.MatchString(emoji) {
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
