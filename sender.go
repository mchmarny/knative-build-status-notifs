package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

var (
	sendNotifsForApp = os.Getenv("NOTIFS_FOR_APP")
	sendStatus       = []string{"WORKING", "SUCCESS", "FAILURE", "INTERNAL_ERROR", "TIMEOUT"}
	token            = os.Getenv("SLACK_API_TOKEN")
	channel          = os.Getenv("SLACK_BUILD_STATUS_CHANNEL")
)

// Send sends the message
func send(msg *CloudBuildNotification) error {

	if msg == nil {
		return fmt.Errorf("Null message on send: %v", msg)
	}

	// to limit number of notifications in Pushover check the name of the app
	if sendNotifsForApp != msg.Substitutions.AppName {
		log.Printf("App not for send. Wanted %s, Got: %s", sendNotifsForApp, msg.Substitutions.AppName)
		return nil
	}

	// check if status is to be sent
	if !isStatusForSend(msg.Status) {
		log.Printf("Status not for send: %s != [%s]", msg.Status, strings.Join(sendStatus, ","))
		return nil
	}

	api := slack.New(token)

	a1 := slack.Attachment{
		Title:     fmt.Sprintf("Trigger: %s", msg.Substitutions.AppName),
		TitleLink: msg.LogURL,
	}
	a1.Fields = []slack.AttachmentField{
		slack.AttachmentField{
			Title: "Tag",
			Value: fmt.Sprintf("Git repo *%s* was tagged: *%s*",
				msg.Source.RepoSource.RepoName,
				msg.Source.RepoSource.TagName),
		},
		slack.AttachmentField{
			Title: "Build",
			Value: fmt.Sprintf("Tag *%s* was built: *%s*",
				msg.Source.RepoSource.TagName,
				msg.Status),
		},
		slack.AttachmentField{
			Title: "Deploy",
			Value: fmt.Sprintf("Image *%s* was deployed to cluster: *%s*",
				msg.Substitutions.AppName,
				msg.Substitutions.ClusterName),
		},
	}

	_, _, err := api.PostMessage(channel,
		slack.MsgOptionText("Cloud Build Status", false),
		slack.MsgOptionAttachments(a1))

	return err

}

func isStatusForSend(a string) bool {
	for _, b := range sendStatus {
		if b == a {
			return true
		}
	}
	return false
}
