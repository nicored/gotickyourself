package main

import (
	"fmt"
	"regexp"

	tickspot "github.com/nicored/gotickyourself"

	"strconv"

	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func getTasksCmd(tick *tickspot.Tick) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "tasks",
		Short:  "List tasks",
		Long:   ``,
		PreRun: initConfigFiles,
		Run:    runTasksCmd,
	}

	cmd.Flags().StringP("filter", "f", "", "Filters by task name")
	cmd.Flags().StringP("project", "p", "", "Filters by project name")
	cmd.Flags().StringP("client", "c", "", "Filters by client name")

	cmd.AddCommand(&cobra.Command{
		Use:    "default",
		Short:  "Set as default",
		Long:   "Example: gty tasks default 1234",
		PreRun: initConfigFiles,
		Run:    defaultTaskCmd,
	})

	aliasCmd := &cobra.Command{
		Use:    "alias",
		Short:  "List all tasks with an alias",
		PreRun: initConfigFiles,
		Run:    aliasCmd,
	}

	addAliasCmd := &cobra.Command{
		Use:    "add",
		Short:  "Add a new alias",
		PreRun: initConfigFiles,
		Run:    addAliasCmd,
	}

	rmAliasCmd := &cobra.Command{
		Use:    "rm",
		Short:  "Remove aliases",
		PreRun: initConfigFiles,
		Run:    removeAliasCmd,
	}
	rmAliasCmd.Flags().BoolP("force", "f", false, "Force remove all")

	aliasCmd.AddCommand(addAliasCmd, rmAliasCmd)
	cmd.AddCommand(aliasCmd)

	return cmd
}

func runTasksCmd(cmd *cobra.Command, args []string) {
	var rt *regexp.Regexp
	filterTask := cmd.Flag("filter").Value.String()
	if filterTask != "" {
		rt = regexp.MustCompile("(?i)" + filterTask)
	}

	var rp *regexp.Regexp
	filterProject := cmd.Flag("project").Value.String()
	if filterProject != "" {
		rp = regexp.MustCompile("(?i)" + filterProject)
	}

	var rc *regexp.Regexp
	filterClient := cmd.Flag("client").Value.String()
	if filterClient != "" {
		rc = regexp.MustCompile("(?i)" + filterClient)
	}

	tasks := []*tickspot.Task{}
	for _, p := range projectsConfig.Projects {
		if rp != nil && !rp.MatchString(p.Name) {
			continue
		}
		if rc != nil && !rc.MatchString(p.Client.Name) {
			continue
		}

		for _, t := range p.Tasks {
			if rt == nil || rt.MatchString(t.Name) {
				tasks = append(tasks, t)
			}
		}
	}

	printTasks(tasks)
}

func aliasCmd(cmd *cobra.Command, args []string) {
	tasks := []*tickspot.Task{}
	for _, t := range Alias {
		tasks = append(tasks, t)
	}

	printTasks(tasks)
}

func removeAliasCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		force := cmd.Flag("force").Value.String()
		removeAllAliases(force == "true")
		os.Exit(0)
	}

	taskID, err := strconv.Atoi(args[0])
	errfOnMismatch(err, nil, "%s is not a valid task ID\n", args[0])

	task, ok := Tasks[taskID]
	errfOnMismatch(ok, true, "Task with ID %d does not exist\n", taskID)

	task.Alias = ""
	updateConfigFile(projectsPath, projectsConfig)
}

func removeAllAliases(force bool) {
	if len(Alias) == 0 {
		log.Println("No alias to remove")
		os.Exit(0)
	}

	if force == false {
		log.Println("You must use the force flag (-f) to remove all the following aliases")
		aliasCmd(nil, nil)
		os.Exit(1)
	}

	for _, t := range Alias {
		t.Alias = ""
	}

	updateConfigFile(projectsPath, projectsConfig)
}

func addAliasCmd(cmd *cobra.Command, args []string) {
	errOnMismatch(len(args), 2, "Too few arguments")

	newAlias := strings.ToLower(strings.TrimSpace(args[0]))
	_, err := isValidAlias(newAlias)
	errfOnMismatch(err, nil, "Invalid alias: %s\n", err)

	taskID, err := strconv.Atoi(args[1])
	errfOnMismatch(err, nil, "%d is not a valid task ID\n")

	task, ok := Tasks[taskID]
	errfOnMismatch(ok, true, "Task ID '%d' does not exist\n")

	task.Alias = newAlias
	updateConfigFile(projectsPath, projectsConfig)
	printTasks([]*tickspot.Task{task})
}

func defaultTaskCmd(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		printDefaultTask()
		return
	}

	taskID, err := strconv.Atoi(args[0])
	errfOnMismatch(err, nil, "Task ID must be an integer. %s\n", err)

	setDefaultTask(taskID)
}

func setDefaultTask(taskID int) {
	task, ok := Tasks[taskID]
	errfOnMismatch(ok, true, "Task ID %d does not exist.\n", taskID)
	projectsConfig.DefaultTask = task

	for tID, task := range Tasks {
		task.IsDefault = tID == taskID
	}

	updateConfigFile(projectsPath, projectsConfig)

	fmt.Println("Default task successfully set:")
	printDefaultTask()
}

func printDefaultTask() {
	if projectsConfig.DefaultTask == nil {
		fmt.Println("You have no default task")
		return
	}

	printTasks([]*tickspot.Task{projectsConfig.DefaultTask})
}

func printTasks(tasks []*tickspot.Task) {
	projects := projectsConfig.Projects

	mappedClients := map[int]map[int][]*tickspot.Task{}
	for _, t := range tasks {
		proj := projects[t.ProjectId]
		if _, ok := mappedClients[proj.ClientId]; !ok {
			mappedClients[proj.ClientId] = map[int][]*tickspot.Task{}
		}

		if _, ok := mappedClients[proj.ClientId][proj.ID]; !ok {
			mappedClients[proj.ClientId][proj.ID] = []*tickspot.Task{}
		}

		mappedClients[proj.ClientId][proj.ID] = append(mappedClients[proj.ClientId][proj.ID], t)
	}

	for cID, projects := range mappedClients {
		fmt.Printf("%s:\n", projectsConfig.Clients[cID].Name)

		for pID, tasks := range projects {
			project := projectsConfig.Projects[pID]
			fmt.Printf("\t%s:\n", project.Name)

			for _, t := range tasks {
				alias := t.Alias
				if alias != "" {
					alias = "(" + t.Alias + ") "
				}

				defaultStr := ""
				if t.IsDefault {
					defaultStr = "(Default)"
				}

				fmt.Printf("\t\t%s%d - %s %s\n", alias, t.ID, t.Name, defaultStr)
			}
		}
	}
}

func isValidAlias(alias string) (bool, error) {
	alias = strings.ToLower(strings.TrimSpace(alias))

	if alias == "" {
		return false, errors.New("Alias cannot be empty")
	}

	for _, reserved := range reservedNames {
		if alias == reserved {
			return false, fmt.Errorf("Alias cannot not be any of the following reserved words: %s", strings.Join(reservedNames, ", "))
		}
	}

	r := regexp.MustCompile("^[a-z]+[a-z0-9\\-_.]+")
	rAlias := r.FindString(alias)

	if alias != rAlias {
		return false, errors.New("An alias must starts with a letter; must not have spaces; can have digits; can have a hyphen (-), an underscore (_) or a dot (.)")
	}

	return true, nil
}
