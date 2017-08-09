package tickspot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"net/http"

	"github.com/snabb/isoweek"
	"regexp"
	"sort"
)

type Pathable interface {
	Path() string
}

type EntryPoint struct {
	path   string
	Method string
}

type Entry struct {
	ID        int     `json:"id,omitempty"`
	Date      string  `json:"date,omitempty"`
	Hours     float64 `json:"hours,omitempty"`
	Notes     string  `json:"notes,omitempty"`
	TaskId    int     `json:"task_id,omitempty"`
	UserId    int     `json:"user_id,omitempty"`
	Url       string  `json:"url,omitempty"`
	CreatedAt string  `json:"created_at,omitempty"`
	UpdatedAt string  `json:"updated_at,omitempty"`
	Task      *Task   `json:"task,omitempty"`
	Billed    bool    `json:"billed"`
}

func (e Entry) GetID() int {
	return e.ID
}

type DateRange struct {
	StartDate string
	EndDate   string
}

func (dr DateRange) getQuery() string {
	return fmt.Sprintf("start_date='%s'&end_date='%s'", dr.StartDate, dr.EndDate)
}

//GET /entries.json
// Either a start_date and end_date have to be provided or an updated_at time
// Each of the following optional parameters can be used to filter the response:
// 	billable (true/false)
//	project_id
//	task_id
//	user_id
//	billed (true/false)
// start_date='2014-09-01'&end_date='2014-09-02'&billable=true"
func (t *Tick) GetEntries(dateRange DateRange) ([]*Entry, error) {
	path := fmt.Sprintf("/entries.json")
	path += "?" + dateRange.getQuery()
	path += "&user_id=" + strconv.Itoa(t.User.ID)

	var b *bytes.Buffer
	resp, err := t.SendRequest(MethodGET, path, b)
	if err != nil {
		return nil, err
	}

	entries := []*Entry{}
	err = json.Unmarshal(resp.Body, &entries)
	if err != nil {
		return nil, err
	}

	return entries, nil
}

//GET /users/4/entries.json
// Either a start_date and end_date have to be provided or an updated_at time
func (t *Tick) GetUserEntries(user *User) ([]*Entry, error) {

	return []*Entry{}, nil
}

//GET /users/4/entries.json
// Either a start_date and end_date have to be provided or an updated_at time
func (t *Tick) GetProjectEntries(project *Project) ([]*Entry, error) {

	return []*Entry{}, nil
}

//GET /tasks/24/entries.json
// Either a start_date and end_date have to be provided or an updated_at time
func (t *Tick) GetTaskEntries(project *Task) ([]*Entry, error) {

	return []*Entry{}, nil
}

//POST /entries.json
func (t *Tick) CreateEntry(date string, hours float64, notes string, task *Task, billed bool) (*Entry, error) {
	path := fmt.Sprintf("/entries.json")

	entry := &Entry{
		Hours:  hours,
		Billed: billed,
		Notes:  notes,
		TaskId: task.ID,
		Date:   date,
	}

	resp, err := t.SendRequest(MethodPOST, path, entry)
	if err != nil {
		return nil, err
	}

	if resp.Code != http.StatusCreated {
		return nil, fmt.Errorf("Request errored and could not create new entry. %s", string(resp.Body))
	}

	err = json.Unmarshal(resp.Body, entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

//PUT /entries/235.json
func (t *Tick) UpdateEntry(entry *Entry) error {
	return nil
}

//DELETE /entries/235.json
func (t *Tick) DeleteEntry(entry *Entry) error {
	return nil
}

func (t *Tick) GetWeeklyTimeSummary() (float64, error) {
	now := time.Now()
	year, week := now.ISOWeek()

	monday := isoweek.StartTime(year, week, time.Local)
	dr := DateRange{
		StartDate: fmt.Sprintf("%d-%d-%d", monday.Year(), monday.Month(), monday.Day()),
		EndDate:   fmt.Sprintf("%d-%d-%d", now.Year(), now.Month(), now.Day()),
	}

	entries, err := t.GetEntries(dr)
	if err != nil {
		return -1.0, err
	}

	total := 0.0
	for _, entry := range entries {
		total += entry.Hours
	}

	return total, nil
}

func (t *Tick) PrintEntries(entries []*Entry) {
	// Sort entries by date
	byDate := func(e1, e2 *Entry) bool {
		return e1.Date < e2.Date
	}
	SortEntriesBy(byDate).Sort(entries)

	dreg := regexp.MustCompile("^\\d{4}-\\d{1,2}-\\d{1,2}")
	currentDate := ""

	for ie, e := range entries {
		eDate := dreg.FindString(e.Date)
		if currentDate != eDate {
			if ie > 0 {
				fmt.Println()
			}

			fmt.Printf("==== Entries for %s ====\n", e.Date)
			currentDate = eDate
		}

		e.Print(t)
	}
}

func (e Entry) Print(t *Tick) {
	task := t.Tasks[e.TaskId]
	if task == nil {
		fmt.Print("-> Task: Unknown")
		fmt.Printf("\t-> Time (hours): %.2f\n", e.Hours)
		return
	}
	fmt.Printf("-> Task: %s\n", task.Name)
	fmt.Printf("\t-> Time (hours): %.2f\n", e.Hours)

	project := t.Projects[task.ProjectId]
	if project == nil {
		fmt.Print("\t-> Project: Unknown")
		return
	}
	fmt.Printf("\t-> Project: %s\n", project.Name)

	client := t.Clients[project.ClientId]
	if client == nil {
		fmt.Print("\t-> Client: Unknown")
		return
	}
	fmt.Printf("\t-> Client: %s\n", client.Name)
}

// SortEntriesBy is the type of a "less" function that defines the ordering of its Entries arguments.
type SortEntriesBy func(e1, e2 *Entry) bool

// entrySorter joins a SortEntriesBy function and a slice of Entries to be sorted.
type entrySorter struct {
	entries []*Entry
	by      func(p1, p2 *Entry) bool // Closure used in the Less method.
}

// Sort is a method on the function type, SortEntriesBy, that sorts the argument slice according to the function.
func (by SortEntriesBy) Sort(entries []*Entry) {
	es := &entrySorter{
		entries: entries,
		by:      by, // The Sort method'e receiver is the function (closure) that defines the sort order.
	}

	sort.Sort(es)
}

// Len is part of sort.Interface.
func (e *entrySorter) Len() int {
	return len(e.entries)
}

// Swap is part of sort.Interface.
func (e *entrySorter) Swap(i, j int) {
	e.entries[i], e.entries[j] = e.entries[j], e.entries[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (e *entrySorter) Less(i, j int) bool {
	return e.by(e.entries[i], e.entries[j])
}
