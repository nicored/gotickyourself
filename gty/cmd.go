package main

import (
	"fmt"
	"log"
	"os"

	tickspot "github.com/nicored/gotickyourself"

	"path/filepath"

	"time"

	"io/ioutil"

	"strings"

	"regexp"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	gtyDir              = ".gty"
	cnfRolesName        = "roles"
	cnfProjectsName     = "projects"
	cnfSettingsName     = "settings"
	baseUrl             = "https://www.tickspot.com"
	updateProjectsAfter = 2 * 24 * time.Hour
)

var (
	homeDir      string
	configPath   string
	rolesPath    string
	projectsPath string
)

var (
	tick           *tickspot.Tick
	rolesConfig    *Roles
	projectsConfig *Projects
	settingsConfig *Settings

	reservedNames = []string{"today", "yesterday", "week", "fortnight", "month"}
)

type Updatable interface {
	LastUpdate() time.Time
	SetUpdate(time.Time)
}

type Settings struct {
	UpdatedAt      time.Time `yaml:"updated_at"`
	HoursPerWeek   float64   `yaml:"hours_per_week"`
	NonWorkingDays []string  `yaml:"non_working_days"`
}

type Roles struct {
	Username  string         `yaml:"username"`
	User      *tickspot.User `yaml:"user"`
	UpdatedAt time.Time      `yaml:"updated_at"`
	Role      *tickspot.Role `yaml:"role"`
}

type Projects struct {
	UpdatedAt   time.Time                 `yaml:"updated_at"`
	Clients     map[int]*tickspot.Client  `yaml:"clients"`
	Projects    map[int]*tickspot.Project `yaml:"projects"`
	DefaultTask *tickspot.Task            `yaml:"-"`
}

var (
	Alias map[string]*tickspot.Task
)

func init() {
	homeDir = getHome()
	configPath = checkConfigDir(homeDir)
}

func main() {
	tick = &tickspot.Tick{
		BaseUrl: baseUrl,
	}

	rootCmd := &cobra.Command{Use: "gty"}

	rootCmd.AddCommand(getInitCmd(tick))
	rootCmd.AddCommand(getSettingsCmd(tick))
	rootCmd.AddCommand(getResetCmd(tick))
	rootCmd.AddCommand(getRolesCmd(tick))
	rootCmd.AddCommand(getUpdateCmd(tick))
	rootCmd.AddCommand(getProjectsCmd(tick))
	rootCmd.AddCommand(getLogCmd(tick))
	rootCmd.AddCommand(getListCmd(tick))
	rootCmd.AddCommand(getSumCmd(tick))
	rootCmd.AddCommand(getTasksCmd(tick))

	rootCmd.Execute()
}

func initConfigFiles(cmd *cobra.Command, args []string) {
	loadSettings()
	rolesPath = loadRoleConfig()
	projectsPath = loadProjects()

	tick.Projects = projectsConfig.Projects
	tick.Clients = projectsConfig.Clients
}

func loadSettings() string {
	settingsConfig = &Settings{
		HoursPerWeek:   40,
		NonWorkingDays: []string{"saturday", "sunday"},
	}

	settingsFile := filepath.Join(configPath, cnfSettingsName+".yml")
	exists := checkConfigFile(settingsFile, settingsConfig)
	if !exists {
		return settingsFile
	}

	fc, err := ioutil.ReadFile(settingsFile)
	errOnMismatch(err, nil, "Could not read settings file")

	err = yaml.Unmarshal(fc, settingsConfig)
	errfOnMismatch(err, nil, "Could not read config file for %s. %s", settingsFile, err)

	return settingsFile
}

func loadRoleConfig() string {
	rolesConfig = &Roles{
		Role: &tickspot.Role{},
	}

	rolesFile := filepath.Join(configPath, cnfRolesName+".yml")
	exists := checkConfigFile(rolesFile, rolesConfig)
	if !exists {
		return rolesFile
	}

	fc, err := ioutil.ReadFile(rolesFile)
	errOnMismatch(err, nil, "Could not read file")

	err = yaml.Unmarshal(fc, rolesConfig)
	errOnMismatch(err, nil, "Could not load roles")

	tick.User = rolesConfig.User
	tick.Role = rolesConfig.Role
	tick.Client = &tickspot.TickClient{
		Username: rolesConfig.Username,
	}

	return rolesFile
}

func loadProjects() string {
	projectsConfig = &Projects{}

	projectsFile := filepath.Join(configPath, cnfProjectsName+".yml")
	exists := checkConfigFile(projectsFile, &Projects{})
	if !exists {
		return projectsFile
	}

	fc, err := ioutil.ReadFile(projectsFile)
	errOnMismatch(err, nil, "Could not read file")

	yaml.Unmarshal(fc, projectsConfig)

	if len(projectsConfig.Projects) == 0 || time.Now().Sub(projectsConfig.UpdatedAt) > updateProjectsAfter {
		updateProjects()
		updateConfigFile(projectsFile, projectsConfig)
	}

	projectsConfig.DefaultTask = getDefaultTask()
	for _, project := range projectsConfig.Projects {
		project.Client = projectsConfig.Clients[project.ClientId]
	}

	indexTasks(projectsConfig.Projects)
	return projectsFile
}

func getDefaultTask() *tickspot.Task {
	for _, project := range projectsConfig.Projects {
		for _, task := range project.Tasks {
			if task.IsDefault == true {
				return task
			}
		}
	}

	return nil
}

func dirExists(dir string) error {
	stat, err := os.Stat(dir)
	if _, ok := err.(*os.PathError); ok {
		err = os.MkdirAll(dir, 0766)
		stat, _ = os.Stat(dir)
	}

	if err != nil {
		return err
	}

	if stat.IsDir() == false {
		return fmt.Errorf("%s is a file, not a directory", dir)
	}

	return nil
}

func errOnMismatch(value interface{}, otherValue interface{}, args ...interface{}) {
	if value != otherValue {
		log.Println(args...)
		os.Exit(1)
	}
}

func errfOnMismatch(value interface{}, otherValue interface{}, msg string, args ...interface{}) {
	if value != otherValue {
		log.Printf(msg, args...)
		os.Exit(1)
	}
}

func getHome() string {
	home, exists := os.LookupEnv("HOME")
	errOnMismatch(exists, true, "Env HOME does not exist")
	return home
}

func checkConfigDir(home string) string {
	configDir := filepath.Join(home, gtyDir)

	// Create .gty directory if it does not exist
	err := dirExists(configDir)
	errOnMismatch(err, nil, "Could not create %s. %s", configDir, err)

	return configDir
}

func checkConfigFile(cnfPath string, dest Updatable) bool {
	_, err := os.Stat(cnfPath)
	if err != nil {
		updateConfigFile(cnfPath, dest)
		return false
	}

	return true
}

func updateConfigFile(cnfPath string, target Updatable) {
	target.SetUpdate(time.Now())

	ymlData, err := yaml.Marshal(target)
	errfOnMismatch(err, nil, "Could not unmarshal config for %s. %s", cnfPath, err)

	f, err := os.Create(cnfPath)
	errfOnMismatch(err, nil, "Could not create config file at %s. %s", cnfPath, err)

	_, err = f.Write(ymlData)
	errfOnMismatch(err, nil, "Could not write to config file at %s. %s", cnfPath, err)
	f.Close()
}

func indexTasks(projects map[int]*tickspot.Project) {
	Alias = map[string]*tickspot.Task{}
	tick.Tasks = map[int]*tickspot.Task{}

	for _, p := range projects {
		for tID, t := range p.Tasks {
			tick.Tasks[tID] = t

			alias := strings.ToLower(strings.TrimSpace(t.Alias))
			if alias != "" {
				Alias[alias] = t
			}
		}
	}
}

func (r *Roles) LastUpdate() time.Time {
	return r.UpdatedAt
}

func (r *Roles) SetUpdate(t time.Time) {
	r.UpdatedAt = t
}

func (p *Projects) LastUpdate() time.Time {
	return p.UpdatedAt
}

func (p *Projects) SetUpdate(t time.Time) {
	p.UpdatedAt = t
}

func (s *Settings) LastUpdate() time.Time {
	return s.UpdatedAt
}

func (s *Settings) SetUpdate(t time.Time) {
	s.UpdatedAt = t
}

func getDateRange(from string) tickspot.DateRange {
	nowTime := time.Now()
	dr := tickspot.DateRange{
		EndDate: fmt.Sprintf("%d-%d-%d", nowTime.Year(), nowTime.Month(), nowTime.Day()),
	}

	rDate := regexp.MustCompile("\\d{4}-\\d{1,2}-\\d{1,2}")
	dateStr := rDate.FindString(from)
	if dateStr != "" {
		dr.StartDate = dateStr
		return dr
	}

	t, isTime := getTimePeriodStart(from)
	errfOnMismatch(isTime, true, "Could not determine time for %s", from)

	return tickspot.DateRange{
		StartDate: fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day()),
		EndDate:   fmt.Sprintf("%d-%d-%d", nowTime.Year(), nowTime.Month(), nowTime.Day()),
	}
}
