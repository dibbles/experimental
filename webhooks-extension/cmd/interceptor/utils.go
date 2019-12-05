/*
 Copyright 2020 The Tekton Authors
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"github.com/google/go-github/github"
	gitlab "github.com/xanzy/go-gitlab"
	"log"
	"net/http"
	"strings"
)

const (
	RequiredRepositoryHeader = "Wext-Repository-Url"
	RequiredEventHeader      = "Wext-Incoming-Event"
	RequiredActionsHeader    = "Wext-Incoming-Actions"
)

type ghPushPayload struct {
	github.PushEvent
	WebhookBranch string `json:"webhooks-tekton-git-branch"`
}

type ghPullRequestPayload struct {
	github.PullRequestEvent
	WebhookBranch string `json:"webhooks-tekton-git-branch"`
}

type glPushPayload struct {
	gitlab.PushEvent
	WebhookBranch string `json:"webhooks-tekton-git-branch"`
}

type glPullRequestPayload struct {
	gitlab.MergeEvent
	WebhookBranch string `json:"webhooks-tekton-git-branch"`
}

func Validate(request *http.Request, httpsCloneURL, eventHeader, pullRequestAction, foundTriggerName string) (bool, error) {

	wantedRepoURL := request.Header.Get(RequiredRepositoryHeader)
	wantedActions := request.Header[RequiredActionsHeader]
	wantedEvents := request.Header[RequiredEventHeader]

	if sanitizeGitInput(httpsCloneURL) == sanitizeGitInput(wantedRepoURL) {
		if request.Header.Get(RequiredEventHeader) != "" {
			foundEvent := request.Header.Get(eventHeader)
			events := strings.Split(wantedEvents[0], ",")
			eventMatch := false
			for _, event := range events {
				if strings.TrimSpace(event) == foundEvent {
					eventMatch = true
					if len(wantedActions) == 0 {
						log.Printf("[%s] Validation PASS (repository URL, secret payload, event type checked)", foundTriggerName)
						return true, nil
					} else {
						actions := strings.Split(wantedActions[0], ",")
						for _, action := range actions {
							if action == pullRequestAction {
								log.Printf("[%s] Validation PASS (repository URL, secret payload, event type, action:%s checked)", foundTriggerName, action)
								return true, nil
							}
						}
					}
				}
			}
			if !eventMatch {
				log.Printf("[%s] Validation FAIL (event type does not match, got %s but wanted one of %s)", foundTriggerName, foundEvent, wantedEvents)
				return false, errors.New("Validator failed as event type does not not match")
			}
			if len(wantedActions) > 0 {
				log.Printf("[%s] Validation FAIL (action type does not match, got %s but wanted one of %s)", foundTriggerName, pullRequestAction, wantedActions)
				return false, errors.New("Validator failed as action does not not match")
			}
			//In theory you wouldn't get here?
			log.Printf("[%s] Validation FAIL (unable to match attributes)", foundTriggerName)
			return false, errors.New("Validator failed")
		}
		// Repository URL matches and no event type restrictions active
		log.Printf("[%s] Validation PASS (repository URL and secret payload checked)", foundTriggerName)
		return true, nil
	}

	log.Printf("[%s] Validation FAIL (repository URLs do not match, got %s but wanted %s)", foundTriggerName, sanitizeGitInput(httpsCloneURL), sanitizeGitInput(wantedRepoURL))
	return false, errors.New("Validator failed as repository URLs do not match")

}

func addBranch(webhookEvent interface{}) ([]byte, error) {
	switch event := webhookEvent.(type) {
	case *github.PushEvent:
		toReturn := ghPushPayload{
			PushEvent:     *event,
			WebhookBranch: event.GetRef()[strings.LastIndex(event.GetRef(), "/")+1:],
		}
		return json.Marshal(toReturn)
	case *github.PullRequestEvent:
		ref := event.GetPullRequest().GetHead().GetRef()
		toReturn := ghPullRequestPayload{
			PullRequestEvent: *event,
			WebhookBranch:    ref[strings.LastIndex(ref, "/")+1:],
		}
		return json.Marshal(toReturn)
	case *gitlab.PushEvent:
		ref := event.Ref
		toReturn := glPushPayload{
			PushEvent:     *event,
			WebhookBranch: ref[strings.LastIndex(ref, "/")+1:],
		}
		return json.Marshal(toReturn)
	case *gitlab.MergeEvent:
		toReturn := glPullRequestPayload{
			MergeEvent:    *event,
			WebhookBranch: event.ObjectAttributes.TargetBranch,
		}
		return json.Marshal(toReturn)
	default:
		//error return
		return []byte{}, errors.New("Unsupported event type received in addBranch()")
	}
}

func sanitizeGitInput(input string) string {
	asLower := strings.ToLower(input)
	noGitSuffix := strings.TrimSuffix(asLower, ".git")
	noHTTPSPrefix := strings.TrimPrefix(noGitSuffix, "https://")
	noHTTPrefix := strings.TrimPrefix(noHTTPSPrefix, "http://")
	return noHTTPrefix
}
