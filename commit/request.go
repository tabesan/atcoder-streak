package commit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	name       string
	repo       string
	URL        *url.URL
	HTTPClient *http.Client
	Data       Commits
}

func NewClient(name, repo string) *Client {
	c := &Client{
		name: name,
		repo: repo,
		HTTPClient: &http.Client{
			Timeout: time.Second * 15,
		},
	}
	c.createURL()
	return c
}

func (c *Client) createURL() {
	c.URL = &url.URL{
		Scheme: "https",
		Host:   "api.github.com",
		Path:   "repos/" + c.name + "/" + c.repo + "/commits/heads/master",
	}
}

func (c *Client) GETRequest() {
	req, err := http.NewRequest("GET", c.URL.String(), nil)
	if err != nil {
		fmt.Println(err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	bytes := []byte(body)
	err = json.Unmarshal(bytes, &c.Data)
	if err != nil {
		fmt.Println(err)
	}
}

func (c Commits) GetDate() time.Time {
	return c.Commit.Author.Date
}
