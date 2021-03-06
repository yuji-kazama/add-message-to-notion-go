package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Client struct {
	baseURL	string
	httpClient *http.Client
}

type Item struct {
	Title string
	DoDate string
	URL string
}

type Page struct {
	Object string `json:"object"`
	Id string `json:"id"`
}


func NewClient() (*Client) {
	c := new(Client)
	c.baseURL = "https://api.notion.com/v1"
	c.httpClient = new(http.Client)
	return c
}

func (c *Client) newRequest(method, spath string, body io.Reader) (*http.Request, error) {
	url := c.baseURL + spath
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer " + os.Getenv("NOTION_INTEGRATION_TOKEN"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", "2021-05-13")

	return req, nil
}

func (c *Client) PostItem(item *Item) (error) {
	var itemJson = `{
		"parent": {
			"database_id": "f50193cd93f2488d8b1dd1c5d3a8cb7d"
		},
		"properties": {
			"Name": {
				"title": [
					{	
						"text": {
							"content": "${TITLE}"
						}
					}
				]
			},
			"Do date": {
				"date": {
					"start": "${DODATE}",
					"end": null
				}
			},
			"Link": {
				"url": "${URL}"
			}
		}
	}`

	itemJson = strings.Replace(string(itemJson), "${TITLE}", item.Title, -1)
	itemJson = strings.Replace(string(itemJson), "${DODATE}", item.DoDate, -1)
	itemJson = strings.Replace(string(itemJson), "${URL}", item.URL, -1)
	
	req, err := c.newRequest(http.MethodPost, "/pages", bytes.NewBuffer([]byte(itemJson))) 
	if err != nil {
		return err
	}
	res, err := c.httpClient.Do(req)
	if (err != nil) {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (c *Client) GetPage(pageId string) (*Page, error) {
	req, err := c.newRequest(http.MethodGet, "/pages/" + pageId, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Read Error:", err)
		return nil, err
	}

	var page Page
	if err := json.Unmarshal(body, &page); err != nil {
		fmt.Printf("Can not unmarshal JSON: %v", err)
		return nil, err
	}

	return &page, nil
}

