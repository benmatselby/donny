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

	envVars := `
In order for donny to integrate with VSTS, you need to define the following environment variables:

* VSTS_ACCOUNT = %s
* VSTS_PROJECT = %s
* VSTS_TEAM    = %s
* VSTS_TOKEN   = %s
`
	if account == "" || project == "" || team == "" || token == "" {
		return false, fmt.Errorf(envVars, account, project, team, token)
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
			Name:  "iteration:cards",
			Usage: "List the work items in a given iteration",
			Action: func(c *cli.Context) {
				args := c.Args()
				if len(args) < 1 {
					fmt.Printf("Please specify an iteration\n")
					cli.ShowSubcommandHelp(c)
					return
				}
				iteration := args[0]

				list := v.GetWorkItemsForIterationByState(iteration)
				fmt.Println(list)
			},
		},
		{
			Name:  "iteration:list",
			Usage: "List all the iterations",
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
