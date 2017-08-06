package tickspot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

type Task struct {
	ID         int      `json:"id,omitempty" yaml:"id"`
	Name       string   `json:"name,omitempty" yaml:"name"`
	Budget     float64  `json:"budget,omitempty" yaml:"budget"`
	Position   int      `json:"position,omitempty" yaml:"position"`
	ProjectId  int      `json:"project_id,omitempty" yaml:"project_id"`
	DateClosed string   `json:"date_closed,omitempty" yaml:"date_closed"`
	Billable   bool     `json:"billable,omitempty" yaml:"billable"`
	Url        string   `json:"url,omitempty" yaml:"url"`
	CreatedAt  string   `json:"created_at,omitempty" yaml:"created_at"`
	UpdatedAt  string   `json:"updated_at,omitempty" yaml:"updated_at"`
	Alias      string   `yaml:"alias"`
	Keywords   []string `yaml:"keywords"`
	IsDefault  bool     `yaml:"is_default"`
}

func (t Task) GetID() int {
	return t.ID
}

//GET /tasks.json
func (t *Tick) GetTasks() ([]*Task, error) {
	path := "/tasks.json"

	page := 1
	tasks := []*Task{}
	for {
		var b *bytes.Buffer
		resp, err := t.SendRequest(MethodGET, path+"?page="+strconv.Itoa(page), b)
		if err != nil {
			return nil, err
		}

		if resp.Code != 200 {
			return nil, fmt.Errorf("CODE %d: %s", resp.Code, string(resp.Body))
		}

		newTasks := []*Task{}
		err = json.Unmarshal(resp.Body, &newTasks)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, newTasks...)
		if len(newTasks) == 0 {
			break
		}

		page++
	}

	return tasks, nil
}
