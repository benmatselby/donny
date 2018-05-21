package main

import (
	"fmt"
	"os"
	"regexp"
	"text/tabwriter"
	"time"

	"github.com/benmatselby/go-vsts/vsts"
	"github.com/urfave/cli"
)

// ListBuilds will call the VSTS API and get a list of builds
func ListBuilds(c *cli.Context) {
	count := c.Int("count")
	filterBranch := c.String("branch")

	options := &vsts.BuildsListOptions{}
	builds, err := client.Builds.List(options)
	if err != nil {
		fmt.Println(err)
	}

	if len(builds) == 0 {
		return
	}

	if len(builds) < count {
		count = len(builds)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.FilterHTML)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", "", "Name", "Branch", "Build", "Finished")
	for index := 0; index < count; index++ {
		build := builds[index]
		name := build.Definition.Name
		result := build.Result
		status := build.Status
		buildNo := build.BuildNumber
		branch := build.Branch

		// Deal with date formatting for the finish time
		finish, err := time.Parse(time.RFC3339, builds[index].FinishTime)
		finishAt := finish.Format(appDateTimeFormat)
		if err != nil {
			finishAt = builds[index].FinishTime
		}

		// Filter on branches
		matched, _ := regexp.MatchString(".*"+filterBranch+".*", branch)
		if matched == false {
			continue
		}

		if status == "inProgress" {
			result = appProgress
		} else if status == "notStarted" {
			result = appPending
		} else {
			if result == "failed" {
				result = appFailure
			} else {
				result = appSuccess
			}
		}

		fmt.Fprintf(w, "%s \t%s\t%s\t%s\t%s\n", result, name, branch, buildNo, finishAt)
	}

	w.Flush()
}
