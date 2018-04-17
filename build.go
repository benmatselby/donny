package main

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

// ListBuilds will call the VSTS API and get a list of builds
func ListBuilds(c *cli.Context) {
	builds, error := client.Builds.List()
	if error != nil {
		fmt.Println(error)
	}

	for index := 0; index < len(builds); index++ {
		name := builds[index].Definition.Name
		// status := builds[index].Status
		result := builds[index].Result

		if result == "failed" {
			color.Red(name)
		} else if result == "succeeded" {
			color.Green(name)
		}
	}
}
