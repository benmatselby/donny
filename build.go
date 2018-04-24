package main

import (
	"fmt"
	"os"
	"regexp"
	"text/tabwriter"
	"time"

	"github.com/urfave/cli"
)

// ListBuilds will call the VSTS API and get a list of builds
func ListBuilds(c *cli.Context) {
	count := c.Int("count")
	filterBranch := c.String("branch")

	builds, error := client.Builds.List()
	if error != nil {
		fmt.Println(error)
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
		name := builds[index].Definition.Name
		result := "✅ "
		buildNo := builds[index].BuildNumber
		branch := builds[index].Branch

		// Deal with date formatting for the finish time
		finish, error := time.Parse(time.RFC3339, builds[index].FinishTime)
		finishAt := finish.Format("2006-01-02 15:04:05")
		if error != nil {
			finishAt = builds[index].FinishTime
		}

		// Filter on branches
		matched, _ := regexp.MatchString(".*"+filterBranch+".*", branch)
		if matched == false {
			continue
		}

		// Provide some UI mechanism to show good/bad builds
		// Wanted to use faith/color, but it doesn't work too well with tabwriter
		if builds[index].Result == "failed" {
			result = "❌ "
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", result, name, branch, buildNo, finishAt)
	}

	w.Flush()
}
