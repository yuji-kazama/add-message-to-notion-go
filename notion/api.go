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
	BaseURL	string
	HTTPClient *http.Client
}

type Page struct {
	Object string `json:"object"`
	Id string `json:"id"`
}

type Item struct {
	Title string
	DoDate string
	URL string
}

type Parent struct {
	DatabaseID string `json:"database_id"`
}

type Text struct {
	Content string `json:"content"`
}

func NewClient(baseURL string) (*Client, error) {
	c := new(Client)
	c.BaseURL = baseURL
	c.HTTPClient = new(http.Client)
	return c, nil
}

func (c *Client) newRequest(method, spath string, body io.Reader) (*http.Request, error) {
	url := c.BaseURL + spath
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
							"content": "%ITEM_TITLE%"
						}
					}
				]
			},
			"Do date": {
				"date": {
					"start": "%ITEM_DODATE%",
					"end": null
				}
			},
			"Link": {
				"url": "%ITEM_URL%"
			}
		}
	}`

	itemJson = strings.Replace(string(itemJson), "%ITEM_TITLE%", item.Title, -1)
	itemJson = strings.Replace(string(itemJson), "%ITEM_DODATE%", item.DoDate, -1)
	itemJson = strings.Replace(string(itemJson), "%ITEM_URL%", item.URL, -1)
	
	req, err := c.newRequest(http.MethodPost, "/pages", bytes.NewBuffer([]byte(itemJson))) 
	if err != nil {
		return err
	}
	res, err := c.HTTPClient.Do(req)
	if (err != nil) {
		return err
	}
	defer res.Body.Close()
	// fmt.Printf("Response: %v", res)
	return nil
}

func (c *Client) GetPage(pageId string) (*Page, error) {
	req, err := c.newRequest(http.MethodGet, "/pages/" + pageId, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
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
	// fmt.Printf("Page: %v", page)

	return &page, nil
}

