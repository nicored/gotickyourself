package main

import (
	"tickspot"

	"fmt"

	"log"

	"time"

	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

func getInitCmd(tick *tickspot.Tick) *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialises gty.",
		Long:  ``,
		Run:   initConfig,
	}
}

func initConfig(cmd *cobra.Command, args []string) {
	rolesConfig = &Roles{}
	tick := &tickspot.Tick{
		Client:  &tickspot.TickClient{Username: rolesConfig.Username, Password: ""},
		BaseUrl: baseUrl,
	}

	var err error
	for {
		fmt.Printf("\nEmail [%s]: ", rolesConfig.Username)
		fmt.Scanln(&rolesConfig.Username)

		fmt.Print("Password: ")
		password, err := terminal.ReadPassword(0)
		fmt.Println()

		tick.Client = &tickspot.TickClient{Username: rolesConfig.Username, Password: string(password)}
		err = initRoles(tick)
		if err != nil {
			log.Println("An error occurred. Make sure your username and password are correct.")
			continue
		}

		break
	}

	rolesPath = filepath.Join(configPath, cnfRolesName+".yml")

	err = initUser(tick)
	errfOnMismatch(err, nil, "", err)

	rolesConfig.User = tick.User
	updateConfigFile(rolesPath, rolesConfig)

	initConfigFiles(nil, nil)
	updateProjects()
	updateConfigFile(projectsPath, projectsConfig)

	fmt.Printf("Thanks %s for using Go Tick Yourself\n", tick.User.FirstName)
}

func initRoles(tick *tickspot.Tick) error {
	var err error

	roles := []*tickspot.Role{}

	fmt.Println("Loading roles...")
	roles, err = tick.GetRoles()
	if err != nil {
		return err
	}

	for r, role := range roles {
		fmt.Printf("%d - Subsc ID: %d  Company: %s - Token: %s\n", r+1, role.SubscriptionID, role.Company, role.APIToken)
	}

	for {
		var roleSelect int
		fmt.Print("Choose a role: ")
		fmt.Scanln(&roleSelect)
		rolesConfig.UpdatedAt = time.Now()

		if roleSelect <= 0 || roleSelect > len(roles) {
			fmt.Println("Invalid Selection")
			continue
		}

		rolesConfig.Role = roles[roleSelect-1]
		tick.Role = rolesConfig.Role

		break
	}

	return nil
}

func initUser(tick *tickspot.Tick) error {
	user, err := tick.GetUserByEmail(tick.Client.Username)
	if err != nil {
		return err
	}

	tick.User = user

	return nil
}
