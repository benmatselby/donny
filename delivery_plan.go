package main

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/benmatselby/go-vsts/vsts"
	"github.com/urfave/cli"
)

// ListDeliveryPlans will call the VSTS API and get a list of delivery plans
func ListDeliveryPlans(c *cli.Context) {

	options := &vsts.DeliveryPlansListOptions{}
	plans, _, error := client.DeliveryPlans.List(options)
	if error != nil {
		fmt.Println(error)
	}

	if len(plans) == 0 {
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.FilterHTML)
	fmt.Fprintf(w, "%s\t%s\n", "Name", "Created")

	for _, plan := range plans {
		// Deal with date formatting for the finish time
		created, error := time.Parse(time.RFC3339, plan.Created)
		createdOn := created.Format(appDateFormat)
		if error != nil {
			createdOn = plan.Created
		}

		fmt.Fprintf(w, "%s\t%s\n", plan.Name, createdOn)
	}
	w.Flush()
}
