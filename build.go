package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/benmatselby/go-azuredevops/azuredevops"
	"github.com/urfave/cli"
)

// ListBuilds will call the API and get a list of builds
func ListBuilds(c *cli.Context) {
	count := c.Int("count")
	filterBranch := c.String("branch")

	options := &azuredevops.BuildsListOptions{
		Count: count,
	}
	builds, err := client.Builds.List(options)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to get a list of builds: %v", err)
		os.Exit(2)
	}

	if len(builds) == 0 {
		return
	}

	renderBuilds(builds, len(builds), filterBranch)
}

// ListBuildOverview will call the API and get a list of builds for a given path
func ListBuildOverview(c *cli.Context) {
	filterBranch := c.String("branch")
	path := c.String("path")

	buildDefOpts := azuredevops.BuildDefinitionsListOptions{Path: "\\" + path}
	definitions, err := client.BuildDefinitions.List(&buildDefOpts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to get a list of build definitions: %v", err)
		os.Exit(2)
	}

	results := make(chan azuredevops.Build)
	var wg sync.WaitGroup
	wg.Add(len(definitions))

	go func() {
		wg.Wait()
		close(results)
	}()

	for _, definition := range definitions {
		go func(definition azuredevops.BuildDefinition) {
			defer wg.Done()

			for _, branchName := range strings.Split(filterBranch, ",") {
				builds, err := getBuildsForBranch(definition.ID, branchName)
				if err != nil {
					fmt.Printf("unable to get builds for definition %s: %v", definition.Name, err)
				}
				if len(builds) > 0 {
					results <- builds[0]
				}
			}
		}(definition)
	}

	var builds []azuredevops.Build
	for result := range results {
		builds = append(builds, result)
	}

	sort.Slice(builds, func(i, j int) bool { return builds[i].Definition.Name < builds[j].Definition.Name })

	renderBuilds(builds, len(builds), ".*")
}

func getBuildsForBranch(defID int, branchName string) ([]azuredevops.Build, error) {
	buildOpts := azuredevops.BuildsListOptions{Definitions: strconv.Itoa(defID), Branch: "refs/heads/" + branchName, Count: 1}
	build, err := client.Builds.List(&buildOpts)
	return build, err
}

func renderBuilds(builds []azuredevops.Build, count int, filterBranch string) {
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
