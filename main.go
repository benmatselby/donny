package main

import (
	"fmt"
	"os"

	"github.com/benmatselby/donny/vsts"
	"github.com/urfave/cli"
)

func main() {
	account := os.Getenv("VSTS_ACCOUNT")
	project := os.Getenv("VSTS_PROJECT")
	team := os.Getenv("VSTS_TEAM")
	token := os.Getenv("VSTS_TOKEN")

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
