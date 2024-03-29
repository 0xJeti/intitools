package intitools

/*	Supported activities

2	Submission 	- Status change
3	Submission 	- Change Severity
5 	Submission 	- Payout
7 	Submission 	- Change vulnerable endpoint
9 	Submission 	- User requires additional feedback
10	Submission 	- User provided feedback
20 	Program		- Status Change
23 	Program		- Update bounties
24 	Program		- Update scope
25 	Program		- Update out of scope
26 	Program		- Update FAQ
27 	Program		- Update domains
28  Program 	- Update rules of engagement
29  Program 	- Update out of scope
47 	Program		- update program

*/
import (
	"context"
	"fmt"
	"log"
	"net/http"
)

type ActivityList struct {
	Completed  bool       `json:"completed"`
	Activities []Activity `json:"activities"`
}

type Activity struct {
	Discriminator   int            `json:"discriminator"`
	Newstatusid     int            `json:"newStatusId"`
	Oldstatusid     int            `json:"oldStatusId"`
	Trigger         int            `json:"trigger"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	Newstate        ResponseState  `json:"newState"`
	User            ResponseUser   `json:"user"`
	UserName        string         `json:"username"`
	Newseverityid   int            `json:"newSeverityId"`
	NewPayoutAmount ResponsePayout `json:"newPayoutAmount"`
	NewPayoutType   int            `json:"newPayoutType"`
	Submissioncode  string         `json:"submissionCode"`
	Submissiontitle string         `json:"submissionTitle"`
	CreatedAt       int64          `json:"createdAt"`
	Programid       string         `json:"programId"`
	Programlogoid   string         `json:"programLogoId"`
	Programname     string         `json:"programName"`
	Programhandle   string         `json:"programHandle"`
	Companyhandle   string         `json:"companyHandle"`
	Newendpoint     string         `json:"newEndpointVulnerableComponent"`
}

type ActivityOptions struct {
	ProgramId          string
	ShowHiddenPrograms bool
	StartDate          int64
}

func (c *Client) GetActivities(ctx context.Context) (*ActivityList, error) {

	apiURL := fmt.Sprintf("%s/core/researcher/dashboard/activity", c.ApiURL)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := ActivityList{}

	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) CheckActivity() (int, error) {

	ctx := c.HttpCtx
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/core/researcher/dashboard/activity/amount?lastviewed=%d", c.ApiURL, c.LastViewed), nil)
	if err != nil {
		return 0, err
	}

	req = req.WithContext(ctx)

	res := 0

	if err := c.sendRequest(req, &res); err != nil {
		return 0, err
	}
	// Log only freshly found activites. This was filling up the log.
	if res > 0 {
		log.Printf("Found %d new activities.\n", res)
	}
	return res, nil
}
func (c *Client) GetSubmissionState(state int) string {
	submmissionStates := []string{
		"Dummy",
		"Triage",
		"Pending",
		"Accepted",
		"Closed",
		"Archived",
		"Unknown: 6",
		"Unknown: 7",
	}

	return submmissionStates[state]
}

func (c *Client) GetClosedState(state int) string {
	closedStates := []string{
		"Dummy",
		"Resolved",
		"Duplicate",
		"Accepted Risk",
		"Informative",
		"Out Of Scope",
		"Spam",
		"Not Applicable",
	}

	return closedStates[state]
}

func (c *Client) GetSeverity(severity int) string {
	severityIds := []string{
		"Dummy",
		"Undecided",
		"Low",
		"Medium",
		"High",
		"Critical",
		"Exceptional",
		"Undecided",
		"Unknown: 8",
	}

	return severityIds[severity]
}

func (c *Client) GetProgramState(program int) string {
	programStates := [1006]string{
		"Dummy",
		"Unknown: 1",
		"Unknown: 2",
		"Open",
		"Suspended",
		"Closing",
		"Closed",
		"Deleted",
	}
	programStates[1001] = "Draft"
	programStates[1002] = "Enrolling"
	programStates[1003] = "Open"
	programStates[1004] = "Closing"
	programStates[1005] = "Closed"

	return programStates[program]
}
