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
26 	Program		- Update FAQ
27 	Program		- Update domains
47 	Program		- update

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
	Discriminator   int           `json:"discriminator"`
	Newstatusid     int           `json:"newStatusId"`
	Oldstatusid     int           `json:"oldStatusId"`
	Trigger         int           `json:"trigger"`
	Title           string        `json:"title"`
	Description     string        `json:"description"`
	Newstate        ResponseState `json:"newState"`
	User            ResponseUser  `json:"user"`
	UserName        string        `json:"username"`
	Newseverityid   int           `json:"newSeverityId"`
	NewPayoutAmount float32       `json:"newPayoutAmount"`
	NewPayoutType   int           `json:"newPayoutType"`
	Submissioncode  string        `json:"submissionCode"`
	Submissiontitle string        `json:"submissionTitle"`
	Createdat       int           `json:"createdAt"`
	Programid       string        `json:"programId"`
	Programlogoid   string        `json:"programLogoId"`
	Programname     string        `json:"programName"`
	Programhandle   string        `json:"programHandle"`
	Companyhandle   string        `json:"companyHandle"`
	Newendpoint     string        `json:"newEndpointVulnerableComponent"`
}

type ActivityOptions struct {
	ProgramId          string
	ShowHiddenPrograms bool
	StartDate          int64
}

const messageTemplate = `{
	"text": "%s",
	"mrkdwn": true,
	"blocks": [
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "%s"
			}
		},
		{
			"type": "divider"
		}
	]
}`

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

func (c *Client) CheckActivity(ctx context.Context) (int, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/core/researcher/dashboard/activity/amount?lastviewed=%d", c.ApiURL, c.LastViewed), nil)
	if err != nil {
		return 0, err
	}

	req = req.WithContext(ctx)

	res := 0

	if err := c.sendRequest(req, &res); err != nil {
		return 0, err
	}

	log.Printf("Checking for new activities: %d\n", res)

	return res, nil
}

func (c *Client) FormatActivityMessage(a Activity) (string, error) {

	submmissionStates := []string{
		"Dummy",
		"Unknown",
		"Pending",
		"Accepted",
		"Closed",
		"Archived",
	}

	closedStates := []string{
		"Dummy",
		"Unknown",
		"Duplicate",
		"Unknown",
		"Informative",
	}

	severityIds := []string{
		"Dummy",
		"Undecided",
		"Low",
		"Medium",
		"High",
		"Critical",
		"Exceptional",
	}

	programStates := []string{

		"Dummy",
		"Unknown",
		"Unknown",
		"Open",
		"Suspended",
		"Closing",
		"Closed",
	}

	var tmp = messageTemplate
	var message string

	submissionLink := fmt.Sprintf("*%s* <https://app.intigriti.com/researcher/submissions/%s/%s|%s>",
		a.Programname, a.Programid, a.Submissioncode, a.Submissiontitle)

	programLink := fmt.Sprintf("<https://app.intigriti.com/researcher/programs/%s/website/detail|%s>",
		a.Programhandle, a.Programname)

	switch d := a.Discriminator; d {

	case 1:
		userRole := a.User.Role
		// Do not send notifications about our own messages
		if userRole != "RESEARCHER" {
			message = fmt.Sprintf("%s\\nNew *message* from *%s* (%s)",
				submissionLink, a.User.Username, userRole)
		}

	//	2	Submission 	- Status change
	case 2:
		newState := submmissionStates[a.Newstate.Status]
		// If status is Closed add reason
		if a.Newstate.Status == 4 {
			newState += " as " + closedStates[a.Newstate.Closereason]
		}

		message = fmt.Sprintf("%s\\nThe *status* changed to `%s`", submissionLink, newState)

	//	3	Submission 	- Change Severity
	case 3:
		message = fmt.Sprintf("%s\\nThe *severity* changed to `%s`", submissionLink, severityIds[a.Newseverityid])

	//	5 	Submission 	- Payout
	case 5:
		message = fmt.Sprintf("%s\\nNew payout *€%.f* :partying_face:", submissionLink, a.NewPayoutAmount)

	//	7 	Submission 	- Change vulnerable endpoint
	case 7:
		message = fmt.Sprintf("%s\\nThe *endpoint / vulnerable component* changed", submissionLink)
	//	9 	Submission 	- User requires additional feedback
	case 9:
		message = fmt.Sprintf("%s\\n*%s* requires additional feedback", submissionLink, a.UserName)
	//	10	Submission 	- User provided feedback
	case 10:
		message = fmt.Sprintf("%s\\n*%s* provided additional feedback", submissionLink, a.UserName)
	//	20 	Program		- Status Change
	case 20:
		message = fmt.Sprintf("%s changed *program status* to `%s`", programLink, programStates[a.Newstatusid])
	//	23 	Program		- Update bounties
	case 23:
		message = fmt.Sprintf("%s updated *bounties*", programLink)
	//	24 	Program		- Update scope
	case 24:
		message = fmt.Sprintf("%s updated *scope*", programLink)
	//	26 	Program		- Update FAQ
	case 26:
		message = fmt.Sprintf("%s updated *FAQ*", programLink)
	//	27 	Program		- Update domains
	case 27:
		message = fmt.Sprintf("%s updated *domains*", programLink)
	//	47 	Program		- Program update published
	case 47:
		message = fmt.Sprintf("%s published a program update: *%s*\\n```%s```", programLink, a.Title, a.Description[:500])

	}
	if message == "" {
		return "", fmt.Errorf("Unknown message type")
	}
	return fmt.Sprintf(tmp, message, message), nil
}
