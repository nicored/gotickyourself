package main

import (
	tickspot "github.com/nicored/gotickyourself"

	"fmt"

	"regexp"

	"github.com/spf13/cobra"
)

func getProjectsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "projects",
		Short:  "Resets all settings",
		Long:   ``,
		PreRun: initConfigFiles,
		Run:    runProjectsCmd,
	}

	cmd.Flags().StringP("filter", "f", "", "Filters by project name")
	cmd.Flags().StringP("client", "c", "", "Filters by client name")
	return cmd
}

func runProjectsCmd(cmd *cobra.Command, args []string) {
	var rp *regexp.Regexp
	filterProject := cmd.Flag("filter").Value.String()
	if filterProject != "" {
		rp = regexp.MustCompile("(?i)" + filterProject)
	}

	var rc *regexp.Regexp
	filterClient := cmd.Flag("client").Value.String()
	if filterClient != "" {
		rc = regexp.MustCompile("(?i)" + filterClient)
	}

	filtered := map[int][]*tickspot.Project{}

	for _, p := range projectsConfig.Projects {
		if (rp == nil || rp.MatchString(p.Name)) && (rc == nil || rc.MatchString(p.Client.Name)) {
			if _, ok := filtered[p.ClientId]; !ok {
				filtered[p.ClientId] = []*tickspot.Project{}
			}

			filtered[p.ClientId] = append(filtered[p.ClientId], p)
		}
	}

	for cID, projects := range filtered {
		fmt.Printf("%s:\n", projectsConfig.Clients[cID].Name)

		for _, project := range projects {
			fmt.Printf("\t%d -  %s\n", project.ID, project.Name)
		}
		fmt.Println()
	}
}

func updateClients() {
	clients, err := tick.GetClients()
	errfOnMismatch(err, nil, "Could not load clients for update", err)

	projectsConfig.Clients = tickspot.IndexClients(clients)
}

func updateProjects() {
	fmt.Println("Updating projects")
	updateClients()

	projects, err := tick.GetProjects(false, true)
	errfOnMismatch(err, nil, "Could not load projects for update.", err)
	projectsMap := tickspot.IndexProjects(projects)

	if projectsConfig.Projects == nil {
		projectsConfig.Projects = map[int]*tickspot.Project{}
	}

	for pId, project := range projectsMap {
		project.Client = projectsConfig.Clients[project.ClientId]

		oldProject, ok := projectsConfig.Projects[pId]
		if !ok {
			fmt.Printf("New project recently added: (%s) %s\n", project.Client.Name, project.Name)
		}

		if oldProject == nil {
			continue
		}

		addedTasks := updateProjectTasks(project.Tasks, oldProject.Tasks)
		if len(addedTasks) > 0 {
			fmt.Println("New tasks added for project ", project.Name)
			for _, task := range addedTasks {
				fmt.Println("\t -> ", task.Name)
			}
		}
	}

	projectsConfig.Projects = projectsMap
}

func updateProjectTasks(updatedTasks, oldTasks map[int]*tickspot.Task) []*tickspot.Task {
	addedTasks := []*tickspot.Task{}

	for tId, task := range updatedTasks {
		oldTask, ok := oldTasks[tId]
		if !ok {
			addedTasks = append(addedTasks, task)
			continue
		}

		task.Alias = oldTask.Alias
		task.Keywords = oldTask.Keywords
		task.IsDefault = oldTask.IsDefault
	}

	return addedTasks
}
