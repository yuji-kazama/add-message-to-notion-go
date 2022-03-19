package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/slack-go/slack"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		verify(r)
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("[ERROR] Failed to read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		payload, err := url.QueryUnescape(string(body))
		if err != nil {
			log.Printf("[ERROR] Failed to unescape: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		payload = strings.Replace(payload, "payload=", "", 1)

		var message slack.InteractionCallback
		if err := json.Unmarshal([]byte(payload), &message); err != nil {
			log.Printf("[ERROR] Failed to unmarshal json: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if message.User.ID != os.Getenv("SLACK_USER_ID") {
			log.Printf("[ERROR] User ID is invalid: %v", message.User.ID);
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		switch message.Type {
		case "message_action":
			modal, err := createModal()
			if err != nil {
				log.Printf("[ERROR] Unable to create modal: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			api := slack.New(os.Getenv("SLACK_ACCESS_TOKEN"))
			if _, err := api.OpenView(message.TriggerID, *modal); err != nil {
				log.Printf("[ERROR] Unable to open view: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return

			}
			w.WriteHeader(http.StatusOK)
			return

		default:
			log.Printf("[ERROR] Unknown request type: %v", message.Type)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})

	log.Println("[INFO] Server listening")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func verify(r *http.Request) error {
	verifier, err := slack.NewSecretsVerifier(r.Header, os.Getenv("SLACK_SIGNING_SECRET"))
	if err != nil {
		log.Printf("[ERROR] Unable to verify secrets: %v", err)
		return err
	}
	return verifier.Ensure()
}

func createModal() (*slack.ModalViewRequest, error) {
	j := `
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
						"initial_value": "hoge"
					},
					"label": {
						"type": "plain_text",
						"text": "Title"
					}
				}
			]
	}`
	var modal slack.ModalViewRequest
	if err := json.Unmarshal([]byte(j), &modal); err!= nil {
		return nil, fmt.Errorf("failed to unmarchal json: %w", err)
	}
	return &modal, nil
}


