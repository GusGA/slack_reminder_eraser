package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var slackToken string
var slackUser string

func init() {
	slackToken = os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		log.Printf("You must set a valid environment variable SLACK_TOKEN")
		os.Exit(1)
	}
	slackUser = os.Getenv("SLACK_USER")
	if slackUser == "" {
		log.Printf("You must set a valid environment variable SLACK_USER")
		os.Exit(1)
	}
}

type reminder struct {
	ID         string `json:"id"`
	Creator    string `json:"creator"`
	User       string `json:"user"`
	Text       string `json:"text"`
	Recurring  bool   `json:"recurring"`
	Time       uint32 `json:"time"`
	CompleteTS uint32 `josn:"complete_ts"`
}

type slackResponse struct {
	Ok        bool        `json:"ok"`
	Reminders []*reminder `json:"reminders"`
}

func getReminders() map[string]*reminder {

	reminders := make(map[string]*reminder)

	resp, err := http.Get("https://slack.com/api/reminders.list?token=" + slackToken)
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
		if rem.User == slackUser {
			reminders[rem.ID] = rem
		}
	}
	return reminders
}

func deleteReminders(reminders map[string]*reminder) {
	type deleteResponse struct {
		Ok    bool   `json:"ok"`
		Error error  `json:",omitempty"`
		ID    string `json:",omitempty"`
	}

	deleteCh := make(chan *deleteResponse, 5)

	for _, rem := range reminders {
		go func(id string) {
			url := fmt.Sprintf("https://slack.com/api/reminders.delete?token=%s&reminder=%s", slackToken, id)
			resp, err := http.Get(url)
			defer resp.Body.Close()
			content := &deleteResponse{Ok: false}
			if err != nil {
				content.Error = err
			} else {
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}
				if err := json.Unmarshal(body, content); err != nil {
					log.Fatal(err)
				}
			}
			content.ID = id
			deleteCh <- content

		}(rem.ID)

		select {
		case content := <-deleteCh:
			if content.Ok {
				fmt.Printf("Reminder id: %s was successfully deleted\n", content.ID)
			} else {
				fmt.Printf("Reminder id %s could not be deleted, error %s\n", content.ID, content.Error)
			}

		}
	}
}

func main() {
	reminders := getReminders()
	deleteReminders(reminders)

}
