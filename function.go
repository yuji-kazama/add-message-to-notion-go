package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/leokite/add-message-to-notion-go/notion"

	"github.com/slack-go/slack"
)


func Function(w http.ResponseWriter, r*http.Request) {
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		verify(r)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("[ERROR] Failed to read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		payload, err := url.QueryUnescape(string(body))
		if err != nil {
			log.Printf("[ERROR] Failed to unescape request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		payload = strings.Replace(payload, "payload=", "", 1)

		var slackMessage slack.InteractionCallback
		if err := json.Unmarshal([]byte(payload), &slackMessage); err != nil {
			log.Printf("[ERROR] Failed to unmarshal json: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if slackMessage.User.ID != os.Getenv("SLACK_USER_ID") {
			log.Printf("[ERROR] User ID is invalid: %v", slackMessage.User.ID);
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		switch slackMessage.Type {
		case "message_action":
			j := getModalJson()
			j = strings.Replace(j, "%INITIAL_DATE%", getTodayDateString(), 1)
			j = strings.Replace(j, "%INITIAL_URL%", getMessageURLString(&slackMessage), 1)
			
			modal, err := createModal(j)
			if err != nil {
				log.Printf("[ERROR] Unable to create modal: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return

			}
			api := slack.New(os.Getenv("SLACK_ACCESS_TOKEN"))
			if _, err := api.OpenView(slackMessage.TriggerID, *modal); err != nil {
				log.Printf("[ERROR] Unable to open view: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return

			}
			w.WriteHeader(http.StatusOK)
			return
		case "view_submission":
			values := slackMessage.View.State.Values
			item := notion.Item{
				Title: values["message"]["message_id"].Value,
				DoDate: values["date"]["date_id"].SelectedDate,
				URL: values["link"]["link_id"].Value,
			}
			c := &notion.Client{
				BaseURL: "https://api.notion.com/v1",
				HTTPClient: new(http.Client),
			}

			if err := c.PostItem(&item); err != nil {
				log.Printf("[ERROR] Failed to call Notion API : %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			return 

		default:
			log.Printf("[ERROR] Unknown request type: %v", slackMessage.Type)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	// })

	// log.Println("[INFO] Server listening")
	// if err := http.ListenAndServe(":8080", nil); err != nil {
	// 	log.Fatal(err)
	// }
}

func verify(r *http.Request) error {
	verifier, err := slack.NewSecretsVerifier(r.Header, os.Getenv("SLACK_SIGNING_SECRET"))
	if err != nil {
		log.Printf("[ERROR] Unable to verify secrets: %v", err)
		return err
	}
	return verifier.Ensure()
}

func createModal(j string) (*slack.ModalViewRequest, error) {
	var modal slack.ModalViewRequest
	if err := json.Unmarshal([]byte(j), &modal); err!= nil {
		return nil, fmt.Errorf("failed to unmarchal json: %w", err)
	}
	return &modal, nil
}

func getTodayDateString() string {
	return time.Now().String()[0:10]
}

func getMessageURLString(m *slack.InteractionCallback) string {
	return "https://" + m.Team.Domain + ".slack.com/archives/" + m.Channel.ID + "/p" + m.MessageTs
}

func getModalJson() string {
	return `
	{
			"type": "modal",
			"title": {
				"type": "plain_text",
				"text": "Add to Notion"
			},
			"submit": {
				"type": "plain_text",
				"text": "Submit"
			},
			"blocks": [
				{
					"block_id": "message",
					"type": "input",
					"element": {
						"action_id": "message_id",
						"type": "plain_text_input",
						"initial_value": ""
					},
					"label": {
						"type": "plain_text",
						"text": "Title"
					}
				},
				{
					"block_id": "date",
					"type": "input",
					"element": {
					  "action_id": "date_id",
					  "type": "datepicker",
					  "initial_date": "%INITIAL_DATE%",
					  "placeholder": {
						"type": "plain_text",
						"text": "Select a date"
					  }
					},
					"label": {
					  "type": "plain_text",
					  "text": "Do Date"
					}
				},
				{
					"block_id": "link",
					"type": "input",
					"element": {
					  "action_id": "link_id",
					  "type": "plain_text_input",
					  "initial_value": "%INITIAL_URL%"
					},
					"label": {
					  "type": "plain_text",
					  "text": "Link"
					}
				}
			]
	}`
}




