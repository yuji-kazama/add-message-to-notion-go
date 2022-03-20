# Slack bot for adding a Slack message to Notion
This is a Slack bot for adding a message in Slack to a database in Notion.

## Requirement
* Go 1.16 or higher
* direnv

## Usag 
Set the following as environment valuables in ".env" file.
* SLACK_SIGNING_SECRET
* SLACK_ACCESS_TOKEN
* SLACK_USER_ID
* NOTION_INTEGRATION_TOKEN
* NOTION_DATABASE_ID

## Setup in local
TODO

## Run in local
* Install ngrok
* Run ngrok
* Modfy the request URL in Slack App
    - Slack api > Features > Interactivity & Shortcut
* Run the following command
    - $ go run cmd/main.go

## Deploy in Clound Functions

$ gcloud functions deploy addMessageToNotionGo --entry-point Function --trigger-http --runtime go116 --region asia-northeast1 --set-env-vars NOTION_INTEGRATION_TOKEN=<ANY_VALUE>,NOTION_DATABASE_ID=<ANY_VALUE>,SLACK_SIGNING_SECRET=<ANY_VALUE>,SLACK_ACCESS_TOKEN=<ANY_VALUE>,SLACK_USER_ID=<ANY_VALUE>

//