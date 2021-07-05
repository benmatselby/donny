package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/benmatselby/go-azuredevops/azuredevops"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ListIterationsOptions defines what arguments/options the user can provide for the
// command.
type ListIterationsOptions struct {
	Args []string
	Team string
}

func NewListIterationsCommand(client *azuredevops.Client) *cobra.Command {
	var opts ListIterationsOptions
	cmd := &cobra.Command{
		Use:   "iterations",
		Short: "Provide a list of iterations",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args

			iterations, err := client.Iterations.List(opts.Team)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not list iterations: %v", err)
				os.Exit(2)
			}

			for index := 0; index < len(iterations); index++ {
				fmt.Println(iterations[index].Name)
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.Team, "team", viper.GetString("AZURE_DEVOPS_TEAM"), "The team to get iterations for")

	return cmd
}

// ListItemsInIterationOptions defines what arguments/options the user can provide for the
// command.
type ListItemsInIterationOptions struct {
	Args      []string
	Team      string
	Iteration string
	Board     string
	HideTag   string
	ShowTag   string
}

// NewListItemsInIteration will call the API and get a list of items for an iteration
func NewListItemsInIteration(client *azuredevops.Client) *cobra.Command {
	var opts ListItemsInIterationOptions
	cmd := &cobra.Command{
		Use:   "iteration-items",
		Short: "Provide a list of items in an iterations",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args

			items, err := getWorkItems(client, opts.Team, opts.Iteration)
			if err != nil {
				return err
			}

			workItems := getWorkItemsByBoardColumn(items)

			// Get the board layout so we now how to render the columns in the right order
			boards, err := client.Boards.List(opts.Team)
			if err != nil {
				return fmt.Errorf("could not list the boards: %v", err)
			}

			// We need to get the specific board we are interested in
			for _, board := range boards {
				if board.Name == opts.Board {
					b, err := client.Boards.Get(opts.Team, board.ID)
					if err != nil {
						return fmt.Errorf("could not get board: %v", err)
					}

					// Now we want a string we can display
					asList := ""
					for _, column := range b.Columns {
						asList += "\n" + column.Name + "\n"
						asList += strings.Repeat("=", len(column.Name)) + "\n"
						for _, item := range workItems[column.Name] {
							if opts.HideTag != "" && stringInSlice(opts.HideTag, item.Fields.TagList) {
								continue
							}

							if opts.ShowTag != "" && !stringInSlice(opts.ShowTag, item.Fields.TagList) {
								continue
							}
							asList += fmt.Sprintf("* (%g) %s\n", item.Fields.Points, item.Fields.Title)
						}
					}
					fmt.Println(asList)

					break
				}
			}

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.Team, "team", viper.GetString("AZURE_DEVOPS_TEAM"), "The team the iteration belongs to")
	flags.StringVar(&opts.Iteration, "iteration", "", "The iteration name")
	flags.StringVar(&opts.Board, "board", "User Stories", "The board name")
	flags.StringVar(&opts.HideTag, "hide-tags", "", "Items to hide based on tags")
	flags.StringVar(&opts.ShowTag, "show-tags", "", "Items to show based on tags")

	return cmd
}

// NewIterationBurndownCommand will display column based data that helps with a daily burndown
func NewIterationBurndownCommand(client *azuredevops.Client) *cobra.Command {
	var opts ListItemsInIterationOptions
	cmd := &cobra.Command{
		Use:   "iteration-burndown",
		Short: "Provide a burndown of the iterations",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args

			items, err := getWorkItems(client, opts.Team, opts.Iteration)
			if err != nil {
				return err
			}

			workItems := getWorkItemsByBoardColumn(items)

			// Get the board layout so we now how to render the columns in the right order
			boards, err := client.Boards.List(opts.Team)
			if err != nil {
				return fmt.Errorf("could not list boards: %v", err)
			}

			// We need to get the specific board we are interested in
			for _, board := range boards {
				if board.Name == opts.Board {
					b, err := client.Boards.Get(opts.Team, board.ID)
					if err != nil {
						return fmt.Errorf("could not get board: %v", err)
					}

					// Now we want a string we can display
					w := tabwriter.NewWriter(os.Stdout, 0, 1, 1, ' ', 0)
					fmt.Fprintf(w, "%s\t%s\t%s\n", "Column", "Items", "Points")
					fmt.Fprintf(w, "%s\t%s\t%s\n", "------", "-----", "------")
					totalItems := 0
					totalPoints := 0.0
					for _, column := range b.Columns {
						points := 0.0
						itemCount := len(workItems[column.Name])

						for _, item := range workItems[column.Name] {
							points += item.Fields.Points
						}
						totalPoints += points
						totalItems += itemCount

						fmt.Fprintf(w, "%s\t%d\t%g\n", column.Name, itemCount, points)
					}

					fmt.Fprintf(w, "%s\t%s\t%s\n", "------", "", "")
					fmt.Fprintf(w, "%s\t%d\t%g\n", "Total", totalItems, totalPoints)
					fmt.Fprintf(w, "%s\t%s\t%s\n", "------", "", "")

					w.Flush()

					break
				}
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.Team, "team", viper.GetString("AZURE_DEVOPS_TEAM"), "The team the iteration belongs to")
	flags.StringVar(&opts.Iteration, "iteration", "", "The iteration name")
	flags.StringVar(&opts.Board, "board", "User Stories", "The board name")
	flags.StringVar(&opts.HideTag, "hide-tags", "", "Items to hide based on tags")
	flags.StringVar(&opts.ShowTag, "show-tags", "", "Items to show based on tags")

	return cmd
}

// NewIterationPeopleBreakdownCommand will display column based data that show person breakdown
func NewIterationPeopleBreakdownCommand(client *azuredevops.Client) *cobra.Command {
	var opts ListItemsInIterationOptions
	cmd := &cobra.Command{
		Use:   "iteration-people",
		Short: "Provide a person breakdown for the iteration",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.Args = args

			items, err := getWorkItems(client, opts.Team, opts.Iteration)
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				os.Exit(2)
			}

			workItems := getWorkItemsByPerson(items)

			w := tabwriter.NewWriter(os.Stdout, 0, 1, 1, ' ', 0)
			fmt.Fprintf(w, "%s\t%s\t%s\n", "Person", "Items", "Points")
			fmt.Fprintf(w, "%s\t%s\t%s\n", "------", "-----", "------")
			totalItems := 0
			totalPoints := 0.0
			for person, items := range workItems {
				points := 0.0
				itemCount := len(items)

				if person == "" {
					person = "Unassigned"
				}

				// Cut the email address out
				person = strings.Split(person, "<")[0]

				for _, item := range items {
					points += item.Fields.Points
				}
				totalPoints += points
				totalItems += itemCount

				fmt.Fprintf(w, "%s\t%d\t%g\n", person, itemCount, points)
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", "------", "", "")
			fmt.Fprintf(w, "%s\t%d\t%g\n", "Total", totalItems, totalPoints)
			fmt.Fprintf(w, "%s\t%s\t%s\n", "------", "", "")
			w.Flush()

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&opts.Team, "team", viper.GetString("AZURE_DEVOPS_TEAM"), "The team the iteration belongs to")
	flags.StringVar(&opts.Iteration, "iteration", "", "The iteration name")
	flags.StringVar(&opts.Board, "board", "Stories", "The board name")
	flags.StringVar(&opts.HideTag, "hide-tags", "", "Items to hide based on tags")
	flags.StringVar(&opts.ShowTag, "show-tags", "", "Items to show based on tags")

	return cmd
}

// getWorkItemsByPerson will return the work items grouped by person
func getWorkItemsByPerson(workItems []azuredevops.WorkItem) map[string][]azuredevops.WorkItem {
	items := make(map[string][]azuredevops.WorkItem)

	// Now build a map|slice|array (!) of
	// Person => Items[]
	for index := 0; index < len(workItems); index++ {
		item := workItems[index]
		key := item.Fields.AssignedTo
		items[key.DisplayName] = append(items[key.DisplayName], item)
	}

	return items
}

func getWorkItemsByBoardColumn(workItems []azuredevops.WorkItem) map[string][]azuredevops.WorkItem {
	items := make(map[string][]azuredevops.WorkItem)

	// Now build a map|slice|array (!) of
	// BoardColumn => Items[]
	for index := 0; index < len(workItems); index++ {
		item := workItems[index]

		if item.Fields.Type == "Task" {
			continue
		}
		key := item.Fields.BoardColumn
		items[key] = append(items[key], item)
	}

	return items
}

func getWorkItems(client *azuredevops.Client, team string, iterationName string) ([]azuredevops.WorkItem, error) {
	iteration, err := client.Iterations.GetByName(team, iterationName)
	if err != nil {
		return nil, fmt.Errorf("could not get work items for %s: %v", iterationName, err)
	}

	if iteration == nil {
		return nil, fmt.Errorf("unable to find iteration: %s", iterationName)
	}

	// Get the items for the iteration we have found
	workItems, err := client.WorkItems.GetForIteration(team, *iteration)
	if err != nil {
		return nil, err
	}

	return workItems, nil
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
