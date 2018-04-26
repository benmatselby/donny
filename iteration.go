package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

// ListIterations will call the VSTS API and get a list of iterations
func ListIterations(c *cli.Context) {
	iterations, error := client.Iterations.List(team)
	if error != nil {
		fmt.Println(error)
	}

	for index := 0; index < len(iterations); index++ {
		fmt.Println(iterations[index].Name)
	}
}

// ListItemsInIteration will call the VSTS API and get a list of items for an iteration
func ListItemsInIteration(c *cli.Context) {
	args := c.Args()
	if len(args) < 1 {
		fmt.Printf("Please specify an iteration\n")
		cli.ShowSubcommandHelp(c)
		return
	}
	iterationName := args[0]

	boardName := c.String("board")

	// Get the iteration by name
	iteration, error := client.Iterations.GetByName(team, iterationName)
	if error != nil {
		fmt.Println(error)
	}

	// Get the items for the iteration we have found
	workItems, error := client.WorkItems.GetForIteration(team, *iteration)
	if error != nil {
		fmt.Println(error)
	}
	items := make(map[string][]string)

	// Now build a map|slice|array (!) of
	// BoardColumn => Items[]
	for index := 0; index < len(workItems); index++ {
		item := workItems[index]

		if item.Fields.Type == "Task" {
			continue
		}
		key := item.Fields.BoardColumn
		value := fmt.Sprintf("* (%g) %s", item.Fields.Points, item.Fields.Title)
		items[key] = append(items[key], value)
	}

	// Get the board layout so we now how to render the columns in the right order
	boards, error := client.Boards.List(team)
	if error != nil {
		fmt.Println(error)
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
				for index := range items[column.Name] {
					asList += items[column.Name][index] + "\n"
				}
			}
			fmt.Println(asList)
		}
	}
}
