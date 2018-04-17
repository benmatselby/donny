package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// ListBuilds will call the VSTS API and get a list of builds
func ListBuilds(c *cli.Context) {
	count := c.Int("count")

	builds, error := client.Builds.List()
	if error != nil {
		fmt.Println(error)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	for index := 0; index < count; index++ {
		name := builds[index].Definition.Name
		result := builds[index].Result
		buildNo := builds[index].BuildNumber
		finish, error := time.Parse(time.RFC3339, builds[index].FinishTime)
		finishAt := finish.Format("2006-01-02 15:04:05")
		if error != nil {
			finishAt = builds[index].FinishTime
		}

		colour := color.New(color.FgGreen)
		if result == "failed" {
			colour = color.New(color.FgRed)
		}

		colour.Fprintf(w, "%s\t%s\t%s\t\n", name, buildNo, finishAt)
	}

	w.Flush()
}
