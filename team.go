package main

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/benmatselby/go-azuredevops/azuredevops"
	"github.com/urfave/cli"
)

// ListTeams will return all the teams for the organisation
func ListTeams(c *cli.Context) {
	filterMine := c.Bool("mine")

	opts := azuredevops.TeamsListOptions{
		Mine: filterMine,
	}
	teams, _, err := client.Teams.List(&opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to get teams: %+v", err)
		os.Exit(2)
	}

	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Name < teams[j].Name
	})

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.FilterHTML)

	for _, team := range teams {
		fmt.Fprintf(w, "%s\n", team.Name)
	}

	w.Flush()
}
