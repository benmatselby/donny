package main

import (
	"fmt"

	"github.com/urfave/cli"
)

func ListBuilds(c *cli.Context) {
	builds, error := client.Builds.List()
	if error != nil {
		fmt.Println(error)
	}

	for index := 0; index < len(builds); index++ {
		fmt.Println(builds[index].Definition.Name)
	}
}
