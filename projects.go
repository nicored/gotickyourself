package tickspot

import (
	"bytes"
	"encoding/json"
	"strconv"
)

type Project struct {
	ID            int           `json:"id,omitempty" yaml:"id"`
	Name          string        `json:"name,omitempty" yaml:"name"`
	Budget        float64       `json:"budget,omitempty" yaml:"budget"`
	DateClosed    string        `json:"date_closed,omitempty" yaml:"date_closed"`
	Notifications bool          `json:"notifications,omitempty" yaml:"notifications"`
	Billable      bool          `json:"billable,omitempty" yaml:"billable"`
	Recurring     bool          `json:"recurring,omitempty" yaml:"recurring"`
	ClientId      int           `json:"client_id,omitempty" yaml:"client_id"`
	OwnerID       int           `json:"owner_id,omitempty" yaml:"owner_id"`
	Url           string        `json:"url,omitempty" yaml:"url"`
	CreatedAt     string        `json:"created_at,omitempty" yaml:"created_at"`
	UpdatedAt     string        `json:"updated_at,omitempty" yaml:"updated_at"`
	Tasks         map[int]*Task `json:"tasks,omitempty" yaml:"tasks"`
	Client        *Client       `yaml:"-"`
}

func (p Project) GetID() int {
	return p.ID
}

type ProjectOptions struct {
	Billable      bool
	Recurring     bool
	Notifications bool
	Budget        float64
}

//GET /projects.json
//GET /projects/closed.json
func (t *Tick) GetProjects(closed, loadTasks bool) ([]*Project, error) {
	path := "/projects.json"
	if closed == true {
		path = "/projects/closed.json"
	}

	page := 1
	projects := []*Project{}

	for {
		var b *bytes.Buffer
		resp, err := t.SendRequest(MethodGET, path+"?page="+strconv.Itoa(page), b)
		if err != nil {
			return nil, err
		}

		newProjects := []*Project{}
		err = json.Unmarshal(resp.Body, &newProjects)
		if err != nil {
			return nil, err
		}

		if len(newProjects) == 0 {
			break
		}

		projects = append(projects, newProjects...)
		page++
	}

	if loadTasks {
		t.LoadProjectsTasks(projects)
	}

	return projects, nil
}

func (t *Tick) LoadProjectsTasks(projects []*Project) error {
	tasks, err := t.GetTasks()
	if err != nil {
		return err
	}

	projectsMap := IndexProjects(projects)
	for _, task := range tasks {
		if project, ok := projectsMap[task.ProjectId]; ok {
			if project.Tasks == nil {
				project.Tasks = map[int]*Task{}
			}

			project.Tasks[task.ID] = task
		}
	}

	return nil
}

func IndexProjects(projects []*Project) map[int]*Project {
	indexedProjects := map[int]*Project{}

	for _, project := range projects {
		if _, ok := indexedProjects[project.ID]; !ok {
			indexedProjects[project.ID] = project
		}
	}

	return indexedProjects
}

//GET /projects/16.json
func (t *Tick) GetProject(id int) (*Project, error) {

	return nil, nil
}

// POST /projects.json
//{
//	"project":
//	{
//		"name":"Prepare Star Destroyer",
//		"budget":50.0,
//		"notifications":false,
//		"billable":true,
//		"recurring":false,
//		"client_id":12,
//		"owner_id":3
//	}
//}
func (t *Tick) CreateProject(name string, options *ProjectOptions) error {

	return nil
}

//PUT /projects/16.json
//{
//	"project":
//	{
//	"budget":300,
//	"billable":true
//	}
//}
func (t *Tick) UpdateProject(project *Project) error {

	return nil
}

//DELETE /projects/16.json
//204 No Content if success
func (t *Tick) DeleteProject(id int) error {

	return nil
}
