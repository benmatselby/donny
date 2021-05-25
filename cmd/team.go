package cmd

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/benmatselby/go-azuredevops/azuredevops"
	"github.com/spf13/cobra"
)

// ListTeamsOptions defines what arguments/options the user can provide for the
// command.
type ListTeamsOptions struct {
	Args []string
	Mine bool
}

func NewListTeamsCommand(client *azuredevops.Client) *cobra.Command {
	var opts ListTeamsOptions
	cmd := &cobra.Command{
		Use:   "teams",
		Short: "Provide a list of teams (Defaults to teams you are in)",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args

			apiOpts := azuredevops.TeamsListOptions{}

			if opts.Mine {
				apiOpts.Mine = opts.Mine
			}
			teams, _, err := client.Teams.List(&apiOpts)
			if err != nil {
				return fmt.Errorf("unable to get teams: %+v", err)
			}

			sort.Slice(teams, func(i, j int) bool {
				return teams[i].Name < teams[j].Name
			})

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.FilterHTML)

			for _, team := range teams {
				fmt.Fprintf(w, "%s\n", team.Name)
			}

			w.Flush()
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&opts.Mine, "mine", false, "List only my teams")

	return cmd
}
