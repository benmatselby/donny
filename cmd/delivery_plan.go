package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/benmatselby/donny/ui"
	"github.com/benmatselby/go-azuredevops/azuredevops"
	"github.com/spf13/cobra"
)

func NewListDeliveryPlansCommand(client *azuredevops.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plans",
		Short: "Provide a list of delivery plans",
		RunE: func(cmd *cobra.Command, args []string) error {
			options := &azuredevops.DeliveryPlansListOptions{}
			plans, _, err := client.DeliveryPlans.List(options)
			if err != nil {
				return fmt.Errorf("unable to get a list of delivery plans: %v", err)
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.FilterHTML)
			fmt.Fprintf(w, "%s\t%s\n", "Name", "Created")

			for _, plan := range plans {
				// Deal with date formatting for the finish time
				created, err := time.Parse(time.RFC3339, plan.Created)
				createdOn := created.Format(ui.AppDateFormat)
				if err != nil {
					createdOn = plan.Created
				}

				fmt.Fprintf(w, "%s\t%s\n", plan.Name, createdOn)
			}
			w.Flush()
			return nil
		},
	}

	return cmd
}

// GetDeliveryPlanOptions defines what arguments/options the user can provide for
// the `repos` command.
type GetDeliveryPlanOptions struct {
	Args     []string
	Plan     string
	ShowTags bool
}

func NewGetDeliveryPlanCommand(client *azuredevops.Client) *cobra.Command {
	var opts GetDeliveryPlanOptions

	cmd := &cobra.Command{
		Use:   "plan",
		Short: "Get information about a delivery plan",
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Plan == "" {
				return fmt.Errorf("please provide a delivery plan")
			}

			options := &azuredevops.DeliveryPlansListOptions{}
			plans, _, err := client.DeliveryPlans.List(options)
			if err != nil {
				return fmt.Errorf("unable to get a list of delivery plans: %v", err)
			}

			for _, plan := range plans {
				if plan.Name == opts.Plan {
					timeline, err := client.DeliveryPlans.GetTimeLine(plan.ID, "", "")
					if err != nil {
						return fmt.Errorf("unable to get the delivery plan timeline time for %s: %v", plan.ID, err)
					}

					start, _ := time.Parse(time.RFC3339, timeline.StartDate)
					end, _ := time.Parse(time.RFC3339, timeline.EndDate)
					fmt.Printf("Name:     %s\n", plan.Name)
					fmt.Printf("Start:    %s\n", start.Format(ui.AppDateFormat))
					fmt.Printf("End:      %s\n", end.Format(ui.AppDateFormat))
					fmt.Printf("Revision: %d\n\n", timeline.Revision)

					for _, team := range timeline.Teams {
						fmt.Println(team.Name)
						fmt.Println(strings.Repeat("=", len(team.Name)))
						fmt.Println()

						for _, iteration := range team.Iterations {
							iStart, _ := time.Parse(time.RFC3339, iteration.StartDate)
							iEnd, _ := time.Parse(time.RFC3339, iteration.EndDate)
							title := fmt.Sprintf(
								"%s (%s - %s)",
								iteration.Name,
								iStart.Format(ui.AppDateFormat),
								iEnd.Format(ui.AppDateFormat),
							)

							fmt.Println(title)
							fmt.Println(strings.Repeat("-", len(title)))
							fmt.Println()

							for _, item := range iteration.WorkItems {
								line := fmt.Sprintf(
									"* %v - %s",
									item[azuredevops.DeliveryPlanWorkItemIDKey],
									item[azuredevops.DeliveryPlanWorkItemNameKey],
								)

								if opts.ShowTags {
									line += fmt.Sprintf(" (%s)", item[azuredevops.DeliveryPlanWorkItemTagKey])
								}

								fmt.Println(line)
							}

							fmt.Println()
							fmt.Println()
						}

						fmt.Println()
					}
					break
				}
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.Plan, "plan", "", "The delivery plan to show")
	flags.BoolVar(&opts.ShowTags, "show-tags", false, "Should tags be displayed")

	return cmd
}
