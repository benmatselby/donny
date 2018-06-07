package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/benmatselby/go-vsts/vsts"
	"github.com/urfave/cli"
)

// ListBuilds will call the VSTS API and get a list of builds
func ListBuilds(c *cli.Context) {
	count := c.Int("count")
	filterBranch := c.String("branch")

	options := &vsts.BuildsListOptions{}
	builds, err := client.Builds.List(options)
	if err != nil {
		fmt.Println(err)
	}

	if len(builds) == 0 {
		return
	}

	if len(builds) < count {
		count = len(builds)
	}

	renderBuilds(builds, count, filterBranch)
}

// ListBuildOverview will call the VSTS API and get a list of builds for a given path
func ListBuildOverview(c *cli.Context) {
	filterBranch := c.String("branch")
	path := c.String("path")

	buildDefOpts := vsts.BuildDefinitionsListOptions{Path: "\\" + path}
	definitions, err := client.BuildDefinitions.List(&buildDefOpts)
	if err != nil {
		fmt.Printf("unable to get a list of build definitions: %v", err)
		return
	}

	var builds []vsts.Build
	for _, definition := range definitions {
		for _, branchName := range strings.Split(filterBranch, ",") {
			build, err := getBuildsForBranch(definition.ID, branchName)
			if err != nil {
				fmt.Printf("unable to get builds for definition %s: %v", definition.Name, err)
			}
			if len(build) > 0 {
				builds = append(builds, build[0])
			}
		}
	}

	renderBuilds(builds, len(builds), ".*")
}

func getBuildsForBranch(defID int, branchName string) ([]vsts.Build, error) {
	buildOpts := vsts.BuildsListOptions{Definitions: strconv.Itoa(defID), Branch: "refs/heads/" + branchName, Count: 1}
	build, err := client.Builds.List(&buildOpts)
	return build, err
}

func renderBuilds(builds []vsts.Build, count int, filterBranch string) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.FilterHTML)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", "", "Name", "Branch", "Build", "Finished")
	for index := 0; index < count; index++ {
		build := builds[index]
		name := build.Definition.Name
		result := build.Result
		status := build.Status
		buildNo := build.BuildNumber
		branch := build.Branch

		// Deal with date formatting for the finish time
		finish, err := time.Parse(time.RFC3339, builds[index].FinishTime)
		finishAt := finish.Format(appDateTimeFormat)
		if err != nil {
			finishAt = builds[index].FinishTime
		}

		// Filter on branches
		matched, _ := regexp.MatchString(".*"+filterBranch+".*", branch)
		if matched == false {
			continue
		}

		if status == "inProgress" {
			result = appProgress
		} else if status == "notStarted" {
			result = appPending
		} else {
			if result == "failed" {
				result = appFailure
			} else {
				result = appSuccess
			}
		}

		fmt.Fprintf(w, "%s \t%s\t%s\t%s\t%s\n", result, name, branch, buildNo, finishAt)
	}

	w.Flush()
}
