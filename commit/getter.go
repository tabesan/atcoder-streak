package commit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Getter interface {
	GetCommit(req *http.Request) []Commits
}

type ComGetter struct {
	HTTPClient *http.Client
}

func NewComGetter() *ComGetter {
	cg := &ComGetter{
		HTTPClient: &http.Client{
			Timeout: time.Second * 15,
		},
	}
	return cg
}

func (c *ComGetter) GetCommit(req *http.Request) []Commits {
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
	var data []Commits
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		fmt.Println(err)
	}

	return data
}
