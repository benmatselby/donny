package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/benmatselby/go-vsts/vsts"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// ListPullRequests will call the VSTS API and get a list of iterations
func ListPullRequests(c *cli.Context) {
	state := c.String("state")
	count := c.Int("count")
	verbose := c.Bool("verbose")
	filterRepo := c.String("repo")

	options := &vsts.PullRequestListOptions{State: state}
	pulls, _, error := client.PullRequests.List(options)
	if error != nil {
		fmt.Println(error)
	}

	if len(pulls) == 0 {
		return
	}

	if len(pulls) < count {
		count = len(pulls)
	}

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
		when, error := time.Parse(time.RFC3339, pull.Created)
		whens := when.Format("2006-01-02 15:04:05")
		if error != nil {
			whens = pull.Created
		}

		// Colourise the title based on state
		if status == "completed" {
			color.Green("#%d %s\n", pull.ID, title)
		} else if status == "abandoned" {
			color.Red("#%d %s\n", pull.ID, title)
		} else {
			color.Yellow("#%d %s\n", pull.ID, title)
		}
		if verbose && pull.Description != "" {
			fmt.Printf("%s\n", pull.Description)
		}
		fmt.Printf("%s\n", repoName)
		fmt.Printf("%v\n", whens)

		fmt.Println("")
	}
}
