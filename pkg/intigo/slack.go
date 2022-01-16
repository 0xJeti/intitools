package intitools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type slackMessage struct {
	Text   string       `json:"text"`
	Mrkdwn bool         `json:"mrkdwn"`
	Blocks []slackBlock `json:"blocks"`
}

type slackBlock struct {
	Type string         `json:"type"`
	Text slackBlockText `json:"text"`
}

type slackBlockText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (c *Client) SlackSend(message string) error {
	webhookURL := c.WebhookURL

	if webhookURL == "" {
		return fmt.Errorf("Webhook not defined.")
	}

	jsonStr := []byte(message)
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	resp, err := ioutil.ReadAll(res.Body)

	if string(resp) != "ok" {
		return fmt.Errorf("cannot send message - %s", string(resp))
	}

	return nil
}

func (c *Client) SlackFormatActivity(a Activity) string {

	var message string

	submissionLink := fmt.Sprintf("*%s* <https://app.intigriti.com/researcher/submissions/%s/%s|%s>",
		url.PathEscape(a.Programname), url.PathEscape(a.Programid), a.Submissioncode, a.Submissiontitle)
	programLink := fmt.Sprintf("<https://app.intigriti.com/researcher/programs/%s/%s|%s>",
		url.PathEscape(a.Companyhandle), url.PathEscape(a.Programhandle), a.Programname)

	switch d := a.Discriminator; d {

	case 1:
		userRole := a.User.Role
		// Do not send notifications about our own messages
		if userRole == "RESEARCHER" {
			return ""
		}

		message = fmt.Sprintf("%s\nNew *message* from *%s* (%s)",
			submissionLink, a.User.Username, userRole)

	//	2	Submission 	- Status change
	case 2:
		newState := c.GetSubmissionState(a.Newstate.Status)
		// If status is Closed add reason
		if a.Newstate.Status == 4 {
			newState += " as " + c.GetClosedState(a.Newstate.Closereason)
		}

		message = fmt.Sprintf("%s\nThe *status* changed to `%s`", submissionLink, newState)

	//	3	Submission 	- Change Severity
	case 3:
		message = fmt.Sprintf("%s\nThe *severity* changed to `%s`", submissionLink, c.GetSeverity(a.Newseverityid))

	//	5 	Submission 	- Payout
	case 5:
		message = fmt.Sprintf("%s\nNew payout *%s %.f* :partying_face:", submissionLink, a.NewPayoutAmount.Currency, a.NewPayoutAmount.Value)

	//	7 	Submission 	- Change vulnerable endpoint
	case 7:
		message = fmt.Sprintf("%s\nThe *endpoint / vulnerable component* changed", submissionLink)
	//	8 	Submission 	- User changed vulnerability type
	case 8:
		message = fmt.Sprintf("%s\n*%s* changed *vulnerability type*", submissionLink, a.UserName)
	//	9 	Submission 	- User requires additional feedback
	case 9:
		message = fmt.Sprintf("%s\n*%s* requires additional feedback", submissionLink, a.UserName)
	//	10	Submission 	- User provided feedback
	case 10:
		message = fmt.Sprintf("%s\n*%s* provided additional feedback", submissionLink, a.UserName)
	//	20 	Program		- Status Change
	case 11:
		message = fmt.Sprintf("%s\n*%s* stopped requesting feedback", submissionLink, a.UserName)
	//	20 	Program		- Status Change
	case 20:
		message = fmt.Sprintf("%s changed *program status* to `%s`", programLink, c.GetProgramState(a.Newstatusid))
	//	22 	Program		- Update description
	case 22:
		descr := a.Description
		if len(descr) > 500 {
			descr = fmt.Sprintf("%s [...]", descr[:500])
		}
		message = fmt.Sprintf("%s changed description: \n```%s```", programLink, descr)
	//	23 	Program		- Update bounties
	case 23:
		message = fmt.Sprintf("%s updated *bounties*", programLink)
	//	24 	Program		- Update scope
	case 24:
		message = fmt.Sprintf("%s updated *scope*", programLink)
	//	25 	Program		- Update out of scope
	case 25:
		message = fmt.Sprintf("%s updated *out of scope*", programLink)
	//	26 	Program		- Update FAQ
	case 26:
		message = fmt.Sprintf("%s updated *FAQ*", programLink)
	//	27 	Program		- Update domains
	case 27:
		message = fmt.Sprintf("%s updated *domains*", programLink)
	//	28 	Program		- Update rules of engagement
	case 28:
		message = fmt.Sprintf("%s updated *rules of engagement*", programLink)
	//	29 	Program		- Update severity assessment
	case 29:
		message = fmt.Sprintf("%s updated *severity assessment*", programLink)
		//	47 	Program		- Program update published
	case 47:
		descr := a.Description
		if len(descr) > 500 {
			descr = fmt.Sprintf("%s [...]", descr[:500])
		}
		message = fmt.Sprintf("%s published a program update: *%s*\n```%s```", programLink, a.Title, descr)

	}
	if message == "" {
		message = fmt.Sprintf("Unknown message type: %d", a.Discriminator)
	}

	blockMsg := slackBlock{
		Type: "section",
		Text: slackBlockText{
			Text: message,
			Type: "mrkdwn",
		},
	}

	block := make([]slackBlock, 0)
	block = append(block, blockMsg)
	slackMsg := slackMessage{
		Text:   message,
		Mrkdwn: true,
		Blocks: block,
	}

	jsonMsg, err := json.Marshal(slackMsg)

	if err != nil {
		return ""
	}

	return string(jsonMsg)

}
