package pushover

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/mchmarny/knative-build-status-notifs/pkg/build"
)

const (
	pushoverAPIEndpoint = "https://api.pushover.net/1/messages.json"
)

var (
	appToke    = os.Getenv("APP_TOKEN")
	userToken  = os.Getenv("USR_TOKEN")
	sendStatus = []string{"SUCCESS", "FAILURE", "INTERNAL_ERROR", "TIMEOUT"}
)

// Send sends the message
func Send(msg *build.CloudBuildNotification) error {

	if msg == nil {
		return fmt.Errorf("Null message on send: %v", msg)
	}

	// check if status is to be sent
	if !isStatusForSend(msg.Status) {
		return nil
	}

	title := fmt.Sprintf("Build Status - %s", msg.Substitutions.AppName)
	body := fmt.Sprintf("Release %s in %s repo was built and pushed to %s cluster.\nFinal Status: %s",
		msg.Source.RepoSource.TagName, msg.Source.RepoSource.RepoName, msg.Substitutions.ClusterName, msg.Status)

	// args
	urlValues := url.Values{}
	urlValues.Add("token", appToke)
	urlValues.Add("user", userToken)
	urlValues.Add("title", title)
	urlValues.Add("message", body)

	// request
	req, err := http.NewRequest("POST", pushoverAPIEndpoint, strings.NewReader(urlValues.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// send
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK {
		return fmt.Errorf("Invalid response status code: %d", resp.StatusCode)
	}

	return nil

}

func isStatusForSend(a string) bool {
	for _, b := range sendStatus {
		if b == a {
			return true
		}
	}
	return false
}
