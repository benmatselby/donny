package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/benmatselby/go-vsts/vsts"
	"github.com/urfave/cli"
)

// ShowGitBranchInfo will get branch information for a repo
func ShowGitBranchInfo(c *cli.Context) {
	if len(c.Args()) < 1 {
		cli.ShowSubcommandHelp(c)
		os.Exit(2)
	}

	repo := c.Args()[0]

	gitRefOps := vsts.GitRefListOptions{
		IncludeStatuses:    true,
		LatestStatusesOnly: true,
	}
	refs, _, err := client.Git.ListRefs(repo, "heads", &gitRefOps)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to get git refs for %s: %+v", repo, err)
		os.Exit(2)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.FilterHTML)

	for _, ref := range refs {
		state := appUnknown
		if len(ref.Statuses) > 0 {
			state = ref.Statuses[0].State

			if state == "failed" {
				state = appFailure
			} else if state == "pending" {
				state = appPending
			} else {
				state = appSuccess
			}
		}

		fmt.Fprintf(w, "%s \t%s\n", state, ref.Name)
	}

	w.Flush()
}