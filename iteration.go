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

// ListCardsInIteration will call the VSTS API and get a list of cards for an iteration
func ListCardsInIteration(c *cli.Context) {
	args := c.Args()
	if len(args) < 1 {
		fmt.Printf("Please specify an iteration\n")
		cli.ShowSubcommandHelp(c)
		return
	}
	iterationName := args[0]

	iteration, error := client.Iterations.GetByName(team, iterationName)
	if error != nil {
		fmt.Println(error)
	}

	workItems, error := client.WorkItems.GetForIteration(team, *iteration)
	if error != nil {
		fmt.Println(error)
	}
	x := make(map[string][]string)

	for index := 0; index < len(workItems); index++ {
		key := workItems[index].Fields.State
		value := fmt.Sprintf("* %s", workItems[index].Fields.Title)
		x[key] = append(x[key], value)
	}

	asList := ""
	for state := range x {
		asList += "\n" + state + "\n"
		asList += strings.Repeat("=", len(state)) + "\n"
		for item := range x[state] {
			asList += x[state][item] + "\n"
		}
	}
	fmt.Println(asList)
}
