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

// ListPullRequests will call the VSTS API and get a list of iterations
func ListPullRequests(c *cli.Context) {
	state := c.String("state")
	count := c.Int("count")
	filterRepo := c.String("repo")
	titleLenth := 40

	options := &vsts.PullRequestListOptions{State: state}
	pulls, _, error := client.PullRequests.List(options)
	if error != nil {
		fmt.Println(error)
	}

	if len(pulls) == 0 {
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "Repo", "Title", "Status", "Created")

	for index := 0; index <= count; index++ {
		pull := pulls[index]

		repoName := pull.Repo.Name

		// Filter on branches
		matched, _ := regexp.MatchString(".*"+filterRepo+".*", repoName)
		if matched == false {
			continue
		}

		title := pull.Title
		if len(title) > titleLenth {
			title = title[0:titleLenth] + "..."
		}
		status := pull.Status

		// Deal with date formatting
		when, error := time.Parse(time.RFC3339, pull.Created)
		whens := when.Format("2006-01-02 15:04:05")
		if error != nil {
			whens = pull.Created
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", repoName, title, status, whens)
	}
	w.Flush()
}
