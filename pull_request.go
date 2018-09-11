package main

import (
	"fmt"
	"os"
	"regexp"
	"text/tabwriter"
	"time"

	"github.com/benmatselby/go-azuredevops/azuredevops"
	"github.com/urfave/cli"
)

// ListPullRequests will call the API and get a list of iterations
func ListPullRequests(c *cli.Context) {
	state := c.String("state")
	count := c.Int("count")
	filterRepo := c.String("repo")

	options := &azuredevops.PullRequestListOptions{State: state}
	pulls, _, err := client.PullRequests.List(options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get list of pull requests: %v", err)
		os.Exit(2)
	}

	if len(pulls) == 0 {
		return
	}

	if len(pulls) < count {
		count = len(pulls)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.FilterHTML)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", "", "ID", "Title", "Repo", "Created")

	for index := 0; index < count; index++ {
		pull := pulls[index]

		repoName := pull.Repo.Name

		// Filter on branches
		matched, _ := regexp.MatchString(".*"+filterRepo+".*", repoName)
		if matched == false {
			continue
		}

		title := pull.Title
		status := pull.Status

		// Deal with date formatting
		when, err := time.Parse(time.RFC3339, pull.Created)
		created := when.Format(appDateTimeFormat)
		if err != nil {
			created = pull.Created
		}

		var result string
		if status == "completed" {
			result = appSuccess
		} else if status == "abandoned" {
			result = appStale
		} else {
			result = appPending
		}

		fmt.Fprintf(w, "%s \t%d\t%s\t%s\t%s\n", result, pull.ID, title, repoName, created)
	}

	w.Flush()
}
