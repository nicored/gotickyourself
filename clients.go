package tickspot

import (
	"bytes"
	"encoding/json"
	"strconv"
)

type Client struct {
	ID        int    `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Archive   bool   `json:"archive,omitempty"`
	Url       string `json:"url,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func (c Client) GetID() int {
	return c.ID
}

// GET /clients.json
func (t *Tick) GetClients() ([]*Client, error) {
	path := "/clients.json"

	page := 1
	clients := []*Client{}

	for {
		var b *bytes.Buffer
		resp, err := t.SendRequest(MethodGET, path+"?page="+strconv.Itoa(page), b)
		if err != nil {
			return nil, err
		}

		newClients := []*Client{}
		err = json.Unmarshal(resp.Body, &newClients)
		if err != nil {
			return nil, err
		}

		if len(newClients) == 0 {
			break
		}

		clients = append(clients, newClients...)
		page++
	}

	return clients, nil
}

// GET /clients/12.json
func (t *Tick) GetClient() (*Client, error) {

	return nil, nil
}

// PUT /clients/12.json
//{
//	"name":"The New Republic",
//	"archive":false
//}
func (t *Tick) CreateClient(name string, archive bool) error {

	return nil
}

func IndexClients(clients []*Client) map[int]*Client {
	indexedClients := map[int]*Client{}

	for _, client := range clients {
		if _, ok := indexedClients[client.ID]; !ok {
			indexedClients[client.ID] = client
		}
	}

	return indexedClients
}
