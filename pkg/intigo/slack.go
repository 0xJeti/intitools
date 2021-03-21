package intitools

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *Client) SlackSend(message string) error {
	webhookURL := c.SlackWebhookURL

	if webhookURL == "" {
		return fmt.Errorf("Slack webhook not defined.")
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
