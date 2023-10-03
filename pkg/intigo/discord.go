package intitools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type discordMessage struct {
	Embeds []discordMsgEmbeds `json:"embeds"`
}

type discordMsgEmbeds struct {
	Color       int          `json:"color"`
	Title       string       `json:"title"`
	URL         string       `json:"url"`
	Description string       `json:"description"`
	Thumbnail   discordThumb `json:"thumbnail"`
}

type discordThumb struct {
	URL string `json:"url"`
}

func (c *Client) DiscordSend(ctx context.Context, message string) error {
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

	err = c.Ratelimiter.Wait(ctx) // This is a blocking call. Honors the rate limit
	if err != nil {
		return err
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode < 200 || res.StatusCode > 299 {
		log.Print(message)
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Printf(bodyString)

		return fmt.Errorf("cannot send message. Error code: %d", res.StatusCode)
	}

	return nil
}

func (c *Client) DiscordFormatActivity(a Activity) (string, error) {

	var message string

	submissionLink := fmt.Sprintf("https://app.intigriti.com/researcher/submissions/%s/%s",
		url.PathEscape(a.Programid), url.PathEscape(a.Submissioncode))
	submissionTitle := fmt.Sprintf("[%s] %s", a.Programname, a.Submissiontitle)

	programLink := fmt.Sprintf("https://app.intigriti.com/researcher/programs/%s/%s",
		url.PathEscape(a.Companyhandle), url.PathEscape(a.Programhandle))
	programTitle := a.Programname

	iconUrl := fmt.Sprintf("https://api.intigriti.com/file/api/file/%s", a.Programlogoid)

	var link string
	var title string

	switch d := a.Discriminator; d {

	case 1:
		userRole := a.User.Role
		// Do not send notifications about our own messages
		if userRole == "RESEARCHER" {
			return "", fmt.Errorf("empty message")
		}

		message = fmt.Sprintf("New **message** from *%s* (%s)",
			a.User.Username, userRole)
		link = submissionLink
		title = submissionTitle

	//	2	Submission 	- Status change
	case 2:
		newState := c.GetSubmissionState(a.Newstate.Status)
		// If status is Closed add reason
		if a.Newstate.Status == 4 {
			newState += " as " + c.GetClosedState(a.Newstate.Closereason)
		}

		message = fmt.Sprintf("The **status** changed to `%s`", newState)
		link = submissionLink
		title = submissionTitle

	//	3	Submission 	- Change Severity
	case 3:
		message = fmt.Sprintf("The **severity** changed to `%s`", c.GetSeverity(a.Newseverityid))
		link = submissionLink
		title = submissionTitle

	//	5 	Submission 	- Payout
	case 5:
		message = fmt.Sprintf("New payout **%s %.f** :partying_face:", a.NewPayoutAmount.Currency, a.NewPayoutAmount.Value)
		link = submissionLink
		title = submissionTitle

	//	7 	Submission 	- Change vulnerable endpoint
	case 7:
		message = fmt.Sprintf("The **endpoint / vulnerable component** changed")
		link = submissionLink
		title = submissionTitle
	//	8 	Submission 	- User changed vulnerability type
	case 8:
		message = fmt.Sprintf("**@%s** changed vulnerability **type**", a.UserName)
		link = submissionLink
		title = submissionTitle

	//	9 	Submission 	- User requires additional feedback
	case 9:
		message = fmt.Sprintf("**@%s** requires additional feedback", a.UserName)
		link = submissionLink
		title = submissionTitle

	//	10	Submission 	- User provided feedback
	case 10:
		message = fmt.Sprintf("**@%s** provided additional feedback", a.UserName)
		link = submissionLink
		title = submissionTitle

	case 11:
		message = fmt.Sprintf("**@%s** stopped requesting feedback", a.UserName)
		link = submissionLink
		title = submissionTitle

	//	20 	Program		- Status Change
	case 20:
		message = fmt.Sprintf("Program changed **status** to `%s`", c.GetProgramState(a.Newstatusid))
		link = programLink
		title = programTitle
	//	22 	Program		- Change description
	case 22:
		descr := a.Description
		if len(descr) > 250 {
			descr = descr[:250]
		}
		message = fmt.Sprintf("Program changed description: \n```%s```", descr)
		link = programLink
		title = programTitle
	//	23 	Program		- Update bounties
	case 23:
		message = fmt.Sprintf("Program updated **bounties**")
		link = programLink
		title = programTitle
	//	24 	Program		- Update scope
	case 24:
		message = fmt.Sprintf("Program updated **scope**")
		link = programLink
		title = programTitle
	//	25 	Program		- Update out of scope
	case 25:
		message = fmt.Sprintf("Program updated **out of scope**")
		link = programLink
		title = programTitle
	//	26 	Program		- Update FAQ
	case 26:
		message = fmt.Sprintf("Program updated **FAQ**")
		link = programLink
		title = programTitle
	//	27 	Program		- Update domains
	case 27:
		message = fmt.Sprintf("Program updated **domains**")
		link = programLink
		title = programTitle
	//	28 	Program		- Update rules of engagement
	case 28:
		message = fmt.Sprintf("Program updated **rules of engagement**")
		link = programLink
		title = programTitle
	//	29 	Program		- Update severity assessment
	case 29:
		message = fmt.Sprintf("Program updated **severity assessment**")
		link = programLink
		title = programTitle
		//	47 	Program		- Program update published
	case 47:
		descr := a.Description
		if len(descr) > 250 {
			descr = descr[:250]
		}
		message = fmt.Sprintf("Program published an update: **%s**\n```%s```", a.Title, descr)
		link = programLink
		title = programTitle

	}

	if message == "" {
		message = fmt.Sprintf("Unknown message type: %d", a.Discriminator)
	}

	embedMsg := discordMsgEmbeds{
		Title:       title,
		URL:         link,
		Description: message,
		Thumbnail: discordThumb{
			URL: iconUrl,
		},
	}

	embed := make([]discordMsgEmbeds, 0)
	embed = append(embed, embedMsg)
	discordMsg := discordMessage{
		Embeds: embed,
	}

	jsonMsg, err := json.Marshal(discordMsg)

	if err != nil {
		return "", err
	}

	return string(jsonMsg), nil

}
