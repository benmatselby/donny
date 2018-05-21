package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/benmatselby/go-vsts/vsts"
	"github.com/urfave/cli"
)

// ListDeliveryPlans will call the VSTS API and get a list of delivery plans
func ListDeliveryPlans(c *cli.Context) {
	options := &vsts.DeliveryPlansListOptions{}
	plans, _, err := client.DeliveryPlans.List(options)
	if err != nil {
		fmt.Println(err)
	}

	if len(plans) == 0 {
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.FilterHTML)
	fmt.Fprintf(w, "%s\t%s\n", "Name", "Created")

	for _, plan := range plans {
		// Deal with date formatting for the finish time
		created, err := time.Parse(time.RFC3339, plan.Created)
		createdOn := created.Format(appDateFormat)
		if err != nil {
			createdOn = plan.Created
		}

		fmt.Fprintf(w, "%s\t%s\n", plan.Name, createdOn)
	}
	w.Flush()
}

// GetDeliveryPlanTimeLine will call the VSTS API and get a list of delivery plans
func GetDeliveryPlanTimeLine(c *cli.Context) {
	planName := c.Args()[0]

	options := &vsts.DeliveryPlansListOptions{}
	plans, _, err := client.DeliveryPlans.List(options)
	if err != nil {
		fmt.Println(err)
	}

	if len(plans) == 0 {
		return
	}

	for _, plan := range plans {
		if plan.Name == planName {
			timeline, err := client.DeliveryPlans.GetTimeLine(plan.ID)
			if err != nil {
				fmt.Println(err)
			}

			start, _ := time.Parse(time.RFC3339, timeline.StartDate)
			end, _ := time.Parse(time.RFC3339, timeline.EndDate)
			fmt.Printf("Name:     %s\n", plan.Name)
			fmt.Printf("Start:    %s\n", start.Format(appDateFormat))
			fmt.Printf("End:      %s\n", end.Format(appDateFormat))
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
						iStart.Format(appDateFormat),
						iEnd.Format(appDateFormat),
					)

					fmt.Println(title)
					fmt.Println(strings.Repeat("-", len(title)))

					for _, item := range iteration.WorkItems {
						fmt.Printf(
							" * %v - %s\n",
							item[vsts.DeliveryPlanWorkItemIDKey],
							item[vsts.DeliveryPlanWorkItemNameKey],
						)
					}

					fmt.Println()
				}

				fmt.Println()
			}
			return
		}
	}
}
