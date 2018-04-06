package vsts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Client for interacting with VSTS
type Client struct {
	Account   string
	Project   string
	Team      string
	AuthToken string
}

// IterationsResponse describes the iterations response
type IterationsResponse struct {
	Iterations []Iteration `json:"value"`
}

// Iteration describes an iteration
type Iteration struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	URL  string `json:"url"`
}

// WorkItemsResponse describes the relationships between work items in VSTS
type WorkItemsResponse struct {
	WorkItemRelationships []WorkItemRelationship `json:"workItemRelations"`
}

// WorkItemRelationship describes the workitem section of the response
type WorkItemRelationship struct {
	Target WorkItemRelation `json:"target"`
}

// WorkItemRelation describes an intermediary between iterations and work items
type WorkItemRelation struct {
	ID int `json:"id"`
}

// WorkItemListResponse describes the list reponse for work items
type WorkItemListResponse struct {
	WorkItems []WorkItem `json:"value"`
}

// WorkItem describes an individual work item in TFS
type WorkItem struct {
	ID     int            `json:"id"`
	Rev    int            `json:"rev"`
	Fields WorkItemFields `json:"fields"`
}

// AsListItem renders a work item in a list format
func (wi WorkItem) AsListItem() string {
	return fmt.Sprintf("* %s", wi.Fields.Title)
}

// WorkItemFields describes all the fields for a given work item
type WorkItemFields struct {
	ID    int    `json:"System.Id"`
	Title string `json:"System.Title"`
	State string `json:"System.State"`
	Type  string `json:"System.WorkItemType"`
}

// New gets the VSTS Client
func New(account string, project string, team string, token string) *Client {
	return &Client{
		Account:   account,
		Project:   project,
		Team:      team,
		AuthToken: token,
	}
}

// Execute runs all the http requests to VSTS
func (c *Client) Execute(request *http.Request, r interface{}) (*http.Response, error) {
	request.SetBasicAuth("", c.AuthToken)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Request to %s responded with status %d", request.URL, response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("Decoding json response from %s failed: %v", request.URL, err)
	}

	return response, nil
}

// GetIterations returns list of the iterations available to the user in VSTS
func (c *Client) GetIterations() ([]Iteration, error) {
	url := fmt.Sprintf(
		"https://%s.visualstudio.com/%s/%s/_apis/work/teamsettings/iterations?api-version=4.1-preview",
		c.Account,
		url.PathEscape(c.Project),
		url.PathEscape(c.Team),
	)

	request, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}
	var ir IterationsResponse
	_, err = c.Execute(request, &ir)

	return ir.Iterations, nil
}

// GetIterationByName will search the iterations for the account and project
// and return a single iteration if the names match
func (c *Client) GetIterationByName(name string) (*Iteration, error) {
	iterations, error := c.GetIterations()
	if error != nil {
		return nil, error
	}

	for index := 0; index < len(iterations); index++ {
		if name == iterations[index].Name {
			iteration := iterations[index]
			return &iteration, nil
		}
	}

	return nil, nil
}

// GetWorkItemIdsForIteration will return an array of ids for a given iteration
func (c *Client) GetWorkItemIdsForIteration(iterationName string) ([]int, error) {
	iteration, error := c.GetIterationByName(iterationName)
	if error != nil {
		return nil, error
	}

	if iteration == nil {
		return nil, nil
	}

	URL := fmt.Sprintf(
		"https://%s.visualstudio.com/%s/%s/_apis/work/teamsettings/iterations/%s/workitems?api-version=%s",
		c.Account,
		url.PathEscape(c.Project),
		url.PathEscape(c.Team),
		iteration.ID,
		"4.1-preview",
	)

	request, err := http.NewRequest("GET", URL, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}

	var response WorkItemsResponse
	_, err = c.Execute(request, &response)

	var queryIds []int
	for index := 0; index < len(response.WorkItemRelationships); index++ {
		relationship := response.WorkItemRelationships[index]
		queryIds = append(queryIds, relationship.Target.ID)
	}

	return queryIds, nil
}

// GetWorkItemsForIteration will get a list of work items based on an iteration name
func (c *Client) GetWorkItemsForIteration(iterationName string) ([]WorkItem, error) {
	queryIds, error := c.GetWorkItemIdsForIteration(iterationName)

	if error != nil {
		return nil, error
	}

	var workIds []string
	for index := 0; index < len(queryIds); index++ {
		workIds = append(workIds, strconv.Itoa(queryIds[index]))
	}

	// Now we want to pad out the fields for the work items
	URL := fmt.Sprintf(
		"https://%s.visualstudio.com/%s/_apis/wit/workitems?ids=%s&fields=%s&api-version=%s",
		c.Account,
		url.PathEscape(c.Project),
		strings.Join(workIds, ","),
		// https://docs.microsoft.com/en-us/rest/api/vsts/wit/work%20item%20types%20field/list
		"System.Id,System.Title,System.State,System.WorkItemType",
		"4.1-preview",
	)

	request, err := http.NewRequest("GET", URL, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}

	var response WorkItemListResponse
	_, err = c.Execute(request, &response)

	return response.WorkItems, nil
}

// GetWorkItemsForIterationByState can build a string response of state => item mappings
func (c *Client) GetWorkItemsForIterationByState(iterationName string) string {
	workItems, error := c.GetWorkItemsForIteration(iterationName)
	if error != nil {
		return ""
	}
	x := make(map[string][]string)

	for index := 0; index < len(workItems); index++ {
		key := workItems[index].Fields.State
		value := workItems[index].AsListItem()
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
	return asList
}
