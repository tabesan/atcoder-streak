package notify

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Notify struct {
	accessToken string
	notifyApi   string
}

// Initialize notify
func NewNotify() *Notify {
	n := new(Notify)
	n.accessToken = "***************"
	n.notifyApi = "https://notify-api.line.me/api/notify"

	return n
}

func (n *Notify) SendNotify(msg string) {
	u, err := url.ParseRequestURI(n.notifyApi)
	if err != nil {
		log.Fatal(err)
	}

	c := &http.Client{}
	form := url.Values{}
	form.Add("message", msg)
	body := strings.NewReader(form.Encode())

	req, err := http.NewRequest("POST", u.String(), body)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+n.accessToken)

	_, err = c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
}
