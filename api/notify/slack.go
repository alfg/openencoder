package notify

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// SendSlackMessage sends a slack webhook post with a message.
// url is the webhook.
func SendSlackMessage(url string, message string) error {
	payload, _ := json.Marshal(map[string]interface{}{
		"attachments": []map[string]string{
			map[string]string{
				"text":  message,
				"color": "good",
			},
		},
	})
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	return err
}
