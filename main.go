package main

import (
	"fmt"
	"os"

	"github.com/benmatselby/donny/version"
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

const (
	appDateFormat     string = "02-01-2006"
	appDateTimeFormat string = "02-01-2006 15:04"
	appSuccess        string = "✅"
	appFailure        string = "❌"
	appPending        string = "🗂"
	appProgress       string = "🏗"
	appStale          string = "🕳"
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

CLI Application to get data out of Visual Studio Team Services into the terminal, where we belong...
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

	app := cli.NewApp()
	app.Name = "donny"
	app.Author = "@benmatselby"
	app.Usage = getUsage(false)
	app.Version = version.GITCOMMIT
	app.Commands = []cli.Command{
		{
			Name:   "build:list",
			Usage:  "        List all the builds",
			Action: ListBuilds,
			Flags: []cli.Flag{
				cli.IntFlag{Name: "count", Value: 10, Usage: "How many builds to display"},
				cli.StringFlag{Name: "branch", Value: ".*", Usage: "Filter by branch name"},
			},
			Category: "build",
		},
		{
			Name:   "iteration:burndown",
			Usage:  "Show column based data for the iteration",
			Action: ShowIterationBurndown,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "board", Value: "Stories", Usage: "Display board by type"},
			},
			Category: "iteration",
		},
		{
			Name:   "iteration:items",
			Usage:  "List the work items in a given iteration",
			Action: ListItemsInIteration,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "board", Value: "Stories", Usage: "Display board by type"},
				cli.StringFlag{Name: "filter-tag", Value: "", Usage: "Filter by a given tag"},
				cli.StringFlag{Name: "hide-tag", Value: "", Usage: "Hide items with a given tag"},
			},
			Category: "iteration",
		},
		{
			Name:     "iteration:list",
			Usage:    "List all the iterations",
			Action:   ListIterations,
			Category: "iteration",
		},
		{
			Name:     "iteration:people",
			Usage:    "Show people based data for the iteration",
			Action:   ShowIterationPeopleBreakdown,
			Category: "iteration",
		},
		{
			Name:     "plan:list",
			Usage:    "     List all the delivery plans",
			Action:   ListDeliveryPlans,
			Category: "plans",
		},
		{
			Name:     "plan:timeline",
			Usage:    "     Show the delivery plan timeline",
			Action:   GetDeliveryPlanTimeLine,
			Category: "plans",
		},
		{
			Name:   "pr:list",
			Usage:  "           List all the pull requests",
			Action: ListPullRequests,
			Flags: []cli.Flag{
				cli.StringFlag{Name: "state", Value: "active", Usage: "Filter by pull request state"},
				cli.StringFlag{Name: "repo", Value: ".*", Usage: "Filter by repo name"},
				cli.IntFlag{Name: "count", Value: 10, Usage: "How many pull requests to display"},
			},
			Category: "pull requests",
		},
	}

	app.Run(os.Args)
}
