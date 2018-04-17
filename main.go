package main

import (
	"fmt"
	"os"

	"github.com/benmatselby/go-vsts/vsts"
	"github.com/urfave/cli"
)

var (
	account string
	project string
	team    string
	token   string
	client  *vsts.Client
)

func loadEnvironmentVars() (bool, error) {
	account = os.Getenv("VSTS_ACCOUNT")
	project = os.Getenv("VSTS_PROJECT")
	team = os.Getenv("VSTS_TEAM")
	token = os.Getenv("VSTS_TOKEN")

	if account == "" || project == "" || team == "" || token == "" {
		return false, fmt.Errorf("Env not all set")
	}

	return true, nil
}

func getUsage(withError bool) string {
	usage := `
_______   ______   .__   __. .__   __. ____    ____
|       \ /  __  \  |  \ |  | |  \ |  | \   \  /   /
|  .--.  |  |  |  | |   \|  | |   \|  |  \   \/   /
|  |  |  |  |  |  | |  .    | |  .    |   \_    _/
|  '--'  |   --'  | |  |\   | |  |\   |     |  |
|_______/ \______/  |__| \__| |__| \__|     |__|

CLI Application to get data out of Visual Studio Team Services into the terminal, where we belong
`
	if withError {
		usage = usage + `

In order for donny to integrate with VSTS, you need to define the following environment variables:

* VSTS_ACCOUNT = %s
* VSTS_PROJECT = %s
* VSTS_TEAM    = %s
* VSTS_TOKEN   = %s
`
	}

	return usage
}

func main() {
	_, err := loadEnvironmentVars()
	if err != nil {
		fmt.Println(getUsage(true))
		return
	}

	client = vsts.NewClient(account, project, token)

	usage := getUsage(false)

	app := cli.NewApp()
	app.Name = "donny"
	app.Author = "@benmatselby"
	app.Usage = usage
	app.Commands = []cli.Command{
		{
			Name:   "iteration:cards",
			Usage:  "List the work items in a given iteration",
			Action: ListCardsInIteration,
		},
		{
			Name:   "iteration:list",
			Usage:  "List all the iterations",
			Action: ListIterations,
		},
		{
			Name:   "build:list",
			Usage:  "List all the builds",
			Action: ListBuilds,
			Flags: []cli.Flag{
				cli.IntFlag{Name: "count", Value: 10, Usage: "How many builds to display"},
			},
		},
	}

	app.Run(os.Args)
}
