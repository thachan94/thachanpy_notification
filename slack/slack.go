package slack

import (
	"fmt"
	"net/http"
	"encoding/json"
	"bytes"
	"time"
	"errors"
	"reflect"
	"../configs"
	"../info"
)

type Field struct {
	Title	string 	`json:"title"`
	Value	string	`json:"value"`
	Short	bool	`json:"short"`
}

type Attachment struct {
	Color 		string		`json:"color"`
	Title 		string 		`json:"title"`
	TitleLink 	string 		`json:"title_link"`
	Text 		string 		`json:"text"`
	Fields 		[]*Field	`json:"fields"`
	Footer		string		`json:"footer"`
}

type Payload struct {
	Attachments []Attachment `json:"attachments"`
}

func (attachment *Attachment) AddField(field Field) *Attachment {
	attachment.Fields = append(attachment.Fields, &field)
	return attachment
}


func getPayload() *Payload {
	config := configs.GetConfigs()
	info := info.GetInfo()

	var color string
	if info.Warning {
		color = config.SlackWarningColor
	} else {
		color = config.SlackOKColor
	}
	attachment := Attachment{
		Color: 		color,
		Title: 		config.SlackTitle,
		TitleLink: 	config.SlackTitleLink,
		Text: 		config.SlackText,
		Footer:		config.SlackFoooter,
	}

	v := reflect.ValueOf(info)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	getType := v.Type()

	for i:=0; i< v.NumField(); i++ {
		field := getType.Field(i).Name
		value := v.Field(i).Interface()

		short := true
		if field == "Uptime" {
			short = false
		}
		if field == "Warning" {
			if fmt.Sprintf("%v", value) == "true" {
				value = fmt.Sprintf("[WARNING] - Disk Used is %s", info.UsedDiskPercent)
				short = false
			}
		}
		_, is_bool := value.(bool)
		if !is_bool {
			attachment.AddField(
				Field {
					Title: field,
					Value: fmt.Sprintf("%s", value),
					Short: short,
				},
			)
		}
	}
	payload := Payload {
		Attachments: []Attachment{attachment},
	}
	return &payload
}

func SendSlackMessage() error {
	config := configs.GetConfigs()
	payload:= getPayload()
	slackBody, _ := json.Marshal(payload)

	req, err := http.NewRequest(
		http.MethodPost,
		config.WebhookURL,
		bytes.NewBuffer(slackBody),
	)
	if err != nil {
		return err
	}

	req.Header.Set(
		"Content-Type", "application/json",
	)

	client := http.Client {
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
    if buf.String() != "ok" {
        return errors.New("Non-ok response returned from Slack")
	}
	
	return nil
}