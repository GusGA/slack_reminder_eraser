package main

import (
	"os"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"log"
	"fmt"
)
var SLACK_TOKEN string
var SLACK_USER string

func init() {
	SLACK_TOKEN = os.Getenv("SLACK_TOKEN")
	SLACK_USER  = os.Getenv("SLACK_USER")
}

const SLACK_URL_GET_REMINDERS = "https://slack.com/api/reminders.list?token="

type Reminder struct {
	Id string `json:"id"`
	Creator string `json:"creator"`
	User string `json:"user"`
	Text string `json:"text"`
	Recurring bool `json:"recurring"`
	Time uint32 `json:"time"`
	CompleteTS uint32 `josn:"complete_ts"`

}

type slackResponse struct {
	Ok bool `json:"ok"`
	Reminders []*Reminder `json:"reminders"`
}


func GetReminders() map[string]*Reminder {

	reminders := make(map[string]*Reminder)

	resp, err := http.Get(SLACK_URL_GET_REMINDERS + SLACK_TOKEN)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var response slackResponse

	if err = json.Unmarshal(body, &response); err != nil {
		log.Fatal(err)
	}

	for _, rem := range response.Reminders {
		if rem.User == SLACK_USER {
			reminders[rem.Id] = rem
		}
	}
	return reminders
}

func DeleteReminders(reminders map[string]*Reminder){
	type deleteResponse struct {
		Ok bool `json:"ok"`
		Error error `json:",omitempty"`
		Id string `json:",omitempty"`
	}

	deleteCh := make(chan *deleteResponse, 5)

	for _, rem := range reminders {
		go func(id string) {
			url := fmt.Sprintf("https://slack.com/api/reminders.delete?token=%s&reminder=%s",SLACK_TOKEN, id)
			resp, err := http.Get(url)
			defer resp.Body.Close()
			content := &deleteResponse{}
			if err != nil {
				content.Error = err
				content.Ok = false
			} else {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
				if err := json.Unmarshal(body, content); err != nil {
					log.Fatal(err)
				}
			}
			content.Id = id
			deleteCh <- content

		}(rem.Id)

		select {
		case content := <- deleteCh:
			if content.Ok {
				fmt.Printf("Reminder id: %s was successfully deleted\n", content.Id)
			} else {
				fmt.Errorf("Reminder id %s could not be deleted, error %s\n", content.Id, content.Error)
			}

		}
	}
}

func main() {
	reminders := GetReminders()
	DeleteReminders(reminders)

}

