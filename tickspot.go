package tickspot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	MethodPOST   = "POST"
	MethodGET    = "GET"
	MethodDELETE = "DELETE"
)

type TickClient struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Tick struct {
	Role    *Role
	User    *User
	Client  *TickClient
	BaseUrl string
}

type Response struct {
	Body []byte
	Code int
}

type WithID interface {
	GetID() int
}

func (t *Tick) SendRequest(method, path string, data interface{}) (*Response, error) {
	req, err := t.prepareRequest(method, path, data)
	if err != nil {
		return nil, err
	}

	return t.sendRequest(req)
}

func (t *Tick) prepareRequest(method, path string, data interface{}) (*http.Request, error) {
	url := t.ReqUrl() + path

	var body *bytes.Buffer
	if data != nil {
		j, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}

		body = bytes.NewBuffer(j)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Token token=%s", t.Role.APIToken))
	req.Header.Add("User-Agent", fmt.Sprintf("GoTickYourself (%s)", t.Client.Username))
	req.Header.Add("Content-Type", "application/json")

	return req, nil
}

func (t *Tick) sendRequest(req *http.Request) (*Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Body: body,
		Code: resp.StatusCode,
	}, nil
}

func (t *Tick) ReqUrl() string {
	return fmt.Sprintf("%s/%d/api/v2", t.BaseUrl, t.Role.SubscriptionID)
}
