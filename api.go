package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
	// Parent string
	Parent struct {
		Database_id string `json:"database_id"`
	} `json:"parent"`
	Properties struct {
		Name struct {
			Title struct {
				Text []struct {
					Content string `json:"content"`
				} `json:"text"`
			} `json:"title"`
			DoDate struct {
				Date struct {
					Start string `json:"start"`
					End string `json:"end"`
				} `json:"date"`
			} `json:"Do date"`
			Link struct {
				Url string `json:"url"`
			} `json:"Link"`
		} `json:"Name"`
	} `json:"properties"`
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

func (c *Client) PostItem() (error) {
	var jsonData = []byte(`{
		"parent": {
			"database_id": "f50193cd93f2488d8b1dd1c5d3a8cb7d"
		},
		"properties": {
			"Name": {
				"title": [
					{	
						"text": {
							"content": "Notion API"
						}
					}
				]
			}
		}
	}`)
	req, err := c.newRequest(http.MethodPost, "/pages", bytes.NewBuffer(jsonData)) 
	if err != nil {
		return err
	}
	res, err := c.HTTPClient.Do(req)
	if (err != nil) {
		return err
	}
	defer res.Body.Close()
	fmt.Printf("Response: %v", res)
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

