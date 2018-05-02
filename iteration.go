package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/benmatselby/go-vsts/vsts"
	"github.com/urfave/cli"
)

// ListIterations will call the VSTS API and get a list of iterations
func ListIterations(c *cli.Context) {
	iterations, error := client.Iterations.List(team)
	if error != nil {
		fmt.Println(error.Error())
		return
	}

	for index := 0; index < len(iterations); index++ {
		fmt.Println(iterations[index].Name)
	}
}

// ListItemsInIteration will call the VSTS API and get a list of items for an iteration
func ListItemsInIteration(c *cli.Context) {
	if !checkIterationDeclared(c) {
		return
	}

	iterationName := c.Args()[0]
	boardName := c.String("board")
	hideTag := c.String("hide-tag")
	showTag := c.String("filter-tag")

	items, error := getWorkItems(team, iterationName)
	if error != nil {
		fmt.Println(error.Error())
		return
	}

	workItems, error := getWorkItemsByBoardColumn(items)
	if error != nil {
		fmt.Println(error.Error())
		return
	}

	// Get the board layout so we now how to render the columns in the right order
	boards, error := client.Boards.List(team)
	if error != nil {
		fmt.Println(error.Error())
		return
	}

	// We need to get the specific board we are interested in
	for _, board := range boards {
		if board.Name == boardName {
			b, error := client.Boards.Get(team, board.ID)
			if error != nil {
				fmt.Println(error)
			}

			// Now we want a string we can display
			asList := ""
			for _, column := range b.Columns {
				asList += "\n" + column.Name + "\n"
				asList += strings.Repeat("=", len(column.Name)) + "\n"
				for _, item := range workItems[column.Name] {
					if hideTag != "" && stringInSlice(hideTag, item.Fields.TagList) {
						continue
					}

					if showTag != "" && !stringInSlice(showTag, item.Fields.TagList) {
						continue
					}
					asList += fmt.Sprintf("* (%g) %s\n", item.Fields.Points, item.Fields.Title)
				}
			}
			fmt.Println(asList)
		}
	}
}

// ShowIterationBurndown will display column based data that helps with a daily burndown
func ShowIterationBurndown(c *cli.Context) {
	if !checkIterationDeclared(c) {
		return
	}

	iterationName := c.Args()[0]
	boardName := c.String("board")

	items, error := getWorkItems(team, iterationName)
	if error != nil {
		fmt.Println(error.Error())
		return
	}

	workItems, error := getWorkItemsByBoardColumn(items)
	if error != nil {
		fmt.Println(error.Error())
		return
	}

	// Get the board layout so we now how to render the columns in the right order
	boards, error := client.Boards.List(team)
	if error != nil {
		fmt.Println(error.Error())
		return
	}

	// We need to get the specific board we are interested in
	for _, board := range boards {
		if board.Name == boardName {
			b, error := client.Boards.Get(team, board.ID)
			if error != nil {
				fmt.Println(error)
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
		}
	}
}

// ShowIterationPeopleBreakdown will display people based data for the iteration
func ShowIterationPeopleBreakdown(c *cli.Context) {
	if !checkIterationDeclared(c) {
		return
	}

	iterationName := c.Args()[0]

	items, error := getWorkItems(team, iterationName)
	if error != nil {
		fmt.Println(error.Error())
		return
	}

	workItems, error := getWorkItemsByPerson(items)
	if error != nil {
		fmt.Println(error.Error())
		return
	}

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
}

func checkIterationDeclared(c *cli.Context) bool {
	if len(c.Args()) < 1 {
		fmt.Printf("Please specify an iteration\n")
		cli.ShowSubcommandHelp(c)
		return false
	}

	return true
}

func getWorkItemsByPerson(workItems []vsts.WorkItem) (map[string][]vsts.WorkItem, error) {
	items := make(map[string][]vsts.WorkItem)

	// Now build a map|slice|array (!) of
	// Person => Items[]
	for index := 0; index < len(workItems); index++ {
		item := workItems[index]
		key := item.Fields.AssignedTo
		items[key] = append(items[key], item)
	}

	return items, nil
}

func getWorkItemsByBoardColumn(workItems []vsts.WorkItem) (map[string][]vsts.WorkItem, error) {
	items := make(map[string][]vsts.WorkItem)

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

	return items, nil
}

func getWorkItems(team string, iterationName string) ([]vsts.WorkItem, error) {
	// Get the iteration by name
	iteration, error := client.Iterations.GetByName(team, iterationName)
	if error != nil {
		return nil, error
	}

	if iteration == nil {
		return nil, fmt.Errorf("Unable to find iteration")
	}

	// Get the items for the iteration we have found
	workItems, error := client.WorkItems.GetForIteration(team, *iteration)
	if error != nil {
		return nil, error
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
