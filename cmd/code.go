package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/benmatselby/donny/ui"
	"github.com/benmatselby/go-azuredevops/azuredevops"
	"github.com/spf13/cobra"
)

// ListBranchInfoOptions defines what arguments/options the user can provide for the
// command.
type ListBranchInfoOptions struct {
	Args []string
	Repo string
}

func NewListBranchInfoCommand(client *azuredevops.Client) *cobra.Command {
	var opts ListBranchInfoOptions
	cmd := &cobra.Command{
		Use:   "branch",
		Short: "Provide a list of pull requests",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args

			if opts.Repo == "" {
				return fmt.Errorf("expected argument of repo")
			}
			gitRefOps := azuredevops.GitRefListOptions{
				IncludeStatuses:    true,
				LatestStatusesOnly: true,
			}
			refs, _, err := client.Git.ListRefs(opts.Repo, "heads", &gitRefOps)
			if err != nil {
				return fmt.Errorf("unable to get git refs for %s: %+v", opts.Repo, err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.FilterHTML)

			for _, ref := range refs {
				state := ui.AppUnknown
				if len(ref.Statuses) > 0 {
					state = ref.Statuses[0].State

					if state == "failed" {
						state = ui.AppFailure
					} else if state == "pending" {
						state = ui.AppPending
					} else {
						state = ui.AppSuccess
					}
				}

				fmt.Fprintf(w, "%s \t%s\n", state, ref.Name)
			}

			w.Flush()

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.Repo, "repo", "", "The repo to filter on")

	return cmd
}
