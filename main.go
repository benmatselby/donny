package main

import (
	"fmt"
	"os"

	"github.com/benmatselby/donny/vsts"
	"github.com/urfave/cli"
)

var (
	account string
	project string
	team    string
	token   string
)

// environmentVars can validate if everything is ok before we start
func environmentVars() (bool, error) {
	account = os.Getenv("VSTS_ACCOUNT")
	project = os.Getenv("VSTS_PROJECT")
	team = os.Getenv("VSTS_TEAM")
	token = os.Getenv("VSTS_TOKEN")

	if account == "" {
		return false, fmt.Errorf("Please define VSTS_ACCOUNT as an environment variable")
	}

	if project == "" {
		return false, fmt.Errorf("Please define VSTS_PROJECT as an environment variable")
	}

	if team == "" {
		return false, fmt.Errorf("Please define VSTS_TEAM as an environment variable")
	}

	if token == "" {
		return false, fmt.Errorf("Please define VSTS_TOKEN as an environment variable")
	}

	return true, nil
}

func main() {

	_, err := environmentVars()
	if err != nil {
		fmt.Println(err)
		return
	}

	v := vsts.New(account, project, team, token)

	app := cli.NewApp()
	app.Name = "donny"
	app.Author = "@benmatselby"
	app.Usage = "CLI Application to get sprint related data out of Visual Studio Team Services"
	app.Commands = []cli.Command{
		{
			Name:  "list:cards",
			Usage: "List the user stories in the sprint",
			Action: func(c *cli.Context) {
				args := c.Args()
				if len(args) < 1 {
					fmt.Printf("Please specify a sprint\n")
					cli.ShowSubcommandHelp(c)
					return
				}
				sprint := args[0]

				workItems, error := v.GetWorkItemsForIteration(sprint)
				for index := 0; index < len(workItems); index++ {
					fmt.Println("* " + workItems[index].Fields.Title)
				}
				if error != nil {
					fmt.Println(error)
				}
			},
		},
		{
			Name:  "list:boards",
			Usage: "List the boards",
			Action: func(c *cli.Context) {
				iterations, error := v.GetIterations()
				for index := 0; index < len(iterations); index++ {
					fmt.Println(iterations[index].Name)
				}

				if error != nil {
					fmt.Println(error)
				}
			},
		},
	}

	app.Run(os.Args)
}
