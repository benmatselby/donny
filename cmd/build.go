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

// ListBuildsOptions defines what arguments/options the user can provide for the
// command.
type ListBuildsOptions struct {
	Args   []string
	Repo   string
	Branch string
	Count  int
}

func NewListBuildsCommand(client *azuredevops.Client) *cobra.Command {
	var opts ListBuildsOptions
	cmd := &cobra.Command{
		Use:   "builds",
		Short: "Provide a list of builds",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args

			options := &azuredevops.BuildsListOptions{
				Count: opts.Count,
			}

			builds, err := client.Builds.List(options)
			if err != nil {
				return fmt.Errorf("unable to get a list of builds: %v", err)
			}

			if len(builds) > 0 {
				renderBuilds(builds, len(builds), opts.Branch)
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.Repo, "repo", "", "The repo to filter on")
	flags.IntVar(&opts.Count, "count", 10, "How many builds to show")
	flags.StringVar(&opts.Branch, "branch", "", "The branch to filter on")

	return cmd
}

func renderBuilds(builds []azuredevops.Build, count int, filterBranch string) {
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
		finishAt := finish.Format(ui.AppDateTimeFormat)
		if err != nil {
			finishAt = builds[index].FinishTime
		}

		// Filter on branches
		matched, _ := regexp.MatchString(".*"+filterBranch+".*", branch)
		if !matched {
			continue
		}

		if status == "inProgress" {
			result = ui.AppProgress
		} else if status == "notStarted" {
			result = ui.AppPending
		} else {
			if result == "failed" {
				result = ui.AppFailure
			} else {
				result = ui.AppSuccess
			}
		}

		fmt.Fprintf(w, "%s \t%s\t%s\t%s\t%s\n", result, name, branch, buildNo, finishAt)
	}

	w.Flush()
}
