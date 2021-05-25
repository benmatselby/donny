package cmd

import (
	"fmt"
	"os"
	"regexp"
	"text/tabwriter"
	"time"

	"github.com/benmatselby/donny/ui"
	"github.com/benmatselby/go-azuredevops/azuredevops"
	"github.com/spf13/cobra"
)

// ListPullRequestsOptions defines what arguments/options the user can provide for the
// command.
type ListPullRequestsOptions struct {
	Args  []string
	Count int
	Repo  string
	State string
}

func NewListPullRequestsCommand(client *azuredevops.Client) *cobra.Command {
	var opts ListPullRequestsOptions
	cmd := &cobra.Command{
		Use:   "prs",
		Short: "Provide a list of pull requests",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args

			options := &azuredevops.PullRequestListOptions{State: opts.State}
			pulls, _, err := client.PullRequests.List(options)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not get list of pull requests: %v", err)
				os.Exit(2)
			}

			if len(pulls) < opts.Count {
				opts.Count = len(pulls)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.FilterHTML)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", "", "ID", "Title", "Repo", "Created")

			for index := 0; index < opts.Count; index++ {
				pull := pulls[index]

				repoName := pull.Repo.Name

				// Filter on branches
				matched, _ := regexp.MatchString(".*"+opts.Repo+".*", repoName)
				if !matched {
					continue
				}

				title := pull.Title
				status := pull.Status

				// Deal with date formatting
				when, err := time.Parse(time.RFC3339, pull.Created)
				created := when.Format(ui.AppDateTimeFormat)
				if err != nil {
					created = pull.Created
				}

				var result string
				if status == "completed" {
					result = ui.AppSuccess
				} else if status == "abandoned" {
					result = ui.AppStale
				} else {
					result = ui.AppPending
				}

				fmt.Fprintf(w, "%s \t%d\t%s\t%s\t%s\n", result, pull.ID, title, repoName, created)
			}

			w.Flush()

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.Repo, "repo", "", "The repo to filter on")
	flags.IntVar(&opts.Count, "count", 10, "Number of pull requests to show")
	flags.StringVar(&opts.State, "state", "active", "The state of the pull requests we want to filter on")

	return cmd
}
