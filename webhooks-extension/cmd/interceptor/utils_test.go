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
	"github.com/google/go-github/github"
	gitlab "github.com/xanzy/go-gitlab"
	"testing"
	"time"
)

func TestSanitizeGitInput(t *testing.T) {
	// Key = URL to process
	// Value = Expected return
	urls := make(map[string]string)
	urls["https://github.com/foo/bar.git"] = "github.com/foo/bar"
	urls["https://github.com/foo/bar"] = "github.com/foo/bar"
	urls["http://github.com/foo/bar.git"] = "github.com/foo/bar"
	urls["http://github.com/foo/bar"] = "github.com/foo/bar"
	urls["github.com/foo/bar"] = "github.com/foo/bar"
	urls["HTTPS://github.com/foo/bar.GIT"] = "github.com/foo/bar"
	urls["hTtP://GiThUb.CoM/FoO/BaR"] = "github.com/foo/bar"
	urls["http://something.else/foo/bar/wibble.git"] = "something.else/foo/bar/wibble"

	for url := range urls {
		sanitized := sanitizeGitInput(url)
		if urls[url] != sanitized {
			t.Errorf("Error santizeGitInput returned unexpected value processing %s, %s was returned but expected %s", url, sanitized, urls[url])
		}
	}
}

func TestAddBranch(t *testing.T) {
	ref := "/blah/head/foo"

	// GitHub Push
	ghPushEventExpectedResult := "{\"ref\":\"/blah/head/foo\",\"webhooks-tekton-git-branch\":\"foo\"}"
	ghPushEvent := github.PushEvent{
		Ref: &ref,
	}
	payload, _ := addBranch(&ghPushEvent)
	if ghPushEventExpectedResult != string(payload) {
		t.Errorf("GitHub push event result unexpected, received %s, expected %s", string(payload), ghPushEventExpectedResult)
	}

	// GitHub Pull Request
	ghPullEventExpectedResult := "{\"pull_request\":{\"head\":{\"ref\":\"/blah/head/foo\"}},\"webhooks-tekton-git-branch\":\"foo\"}"
	ghPullEvent := github.PullRequestEvent{
		PullRequest: &github.PullRequest{
			Head: &github.PullRequestBranch{
				Ref: &ref,
			},
		},
	}
	payload, _ = addBranch(&ghPullEvent)
	if ghPullEventExpectedResult != string(payload) {
		t.Errorf("GitHub pull request event result unexpected, received %s, expected %s", string(payload), ghPullEventExpectedResult)
	}

	// GitLab Push
	glPushEventExpectedResult := "{\"object_kind\":\"\",\"before\":\"\",\"after\":\"\",\"ref\":\"/blah/head/foo\",\"checkout_sha\":\"\",\"user_id\":0,\"user_name\":\"\",\"user_username\":\"\",\"user_email\":\"\",\"user_avatar\":\"\",\"project_id\":0,\"project\":{\"name\":\"\",\"description\":\"\",\"avatar_url\":\"\",\"git_ssh_url\":\"\",\"git_http_url\":\"\",\"namespace\":\"\",\"path_with_namespace\":\"\",\"default_branch\":\"\",\"homepage\":\"\",\"url\":\"\",\"ssh_url\":\"\",\"http_url\":\"\",\"web_url\":\"\",\"visibility\":\"\"},\"repository\":null,\"commits\":null,\"total_commits_count\":0,\"webhooks-tekton-git-branch\":\"foo\"}"
	glPushEvent := gitlab.PushEvent{
		Ref: ref,
	}
	payload, _ = addBranch(&glPushEvent)
	if glPushEventExpectedResult != string(payload) {
		t.Errorf("GitLab push event result unexpected, received %s, expected %s", string(payload), glPushEventExpectedResult)
	}

	// Unsupported Event
	unsupportedEvent := github.StarEvent{
		Action: &ref,
	}
	payload, err := addBranch(&unsupportedEvent)
	if "" != string(payload) {
		t.Errorf("Unsupported event result unexpected, received %s, expected \"\"", string(payload))
	}
	if err.Error() != "Unsupported event type received in addBranch()" {
		t.Errorf("Unexpected error received")
	}

}

func TestAddBranchGitLabMergeRequest(t *testing.T) {
	//We need to mock up more of a struct for gitlab merge requests,
	//so we have this in a seperate test just for some code clartiy
	type ObjectAttributes struct {
		ID                       int                 `json:"id"`
		TargetBranch             string              `json:"target_branch"`
		SourceBranch             string              `json:"source_branch"`
		SourceProjectID          int                 `json:"source_project_id"`
		AuthorID                 int                 `json:"author_id"`
		AssigneeID               int                 `json:"assignee_id"`
		Title                    string              `json:"title"`
		CreatedAt                string              `json:"created_at"` // Should be *time.Time (see Gitlab issue #21468)
		UpdatedAt                string              `json:"updated_at"` // Should be *time.Time (see Gitlab issue #21468)
		StCommits                []*gitlab.Commit    `json:"st_commits"`
		StDiffs                  []*gitlab.Diff      `json:"st_diffs"`
		MilestoneID              int                 `json:"milestone_id"`
		State                    string              `json:"state"`
		MergeStatus              string              `json:"merge_status"`
		TargetProjectID          int                 `json:"target_project_id"`
		IID                      int                 `json:"iid"`
		Description              string              `json:"description"`
		Position                 int                 `json:"position"`
		LockedAt                 string              `json:"locked_at"`
		UpdatedByID              int                 `json:"updated_by_id"`
		MergeError               string              `json:"merge_error"`
		MergeParams              *gitlab.MergeParams `json:"merge_params"`
		MergeWhenBuildSucceeds   bool                `json:"merge_when_build_succeeds"`
		MergeUserID              int                 `json:"merge_user_id"`
		MergeCommitSHA           string              `json:"merge_commit_sha"`
		DeletedAt                string              `json:"deleted_at"`
		ApprovalsBeforeMerge     string              `json:"approvals_before_merge"`
		RebaseCommitSHA          string              `json:"rebase_commit_sha"`
		InProgressMergeCommitSHA string              `json:"in_progress_merge_commit_sha"`
		LockVersion              int                 `json:"lock_version"`
		TimeEstimate             int                 `json:"time_estimate"`
		Source                   *gitlab.Repository  `json:"source"`
		Target                   *gitlab.Repository  `json:"target"`
		LastCommit               struct {
			ID        string     `json:"id"`
			Message   string     `json:"message"`
			Timestamp *time.Time `json:"timestamp"`
			URL       string     `json:"url"`
			Author    struct {
				Name  string `json:"name"`
				Email string `json:"email"`
			} `json:"author"`
		} `json:"last_commit"`
		WorkInProgress bool                 `json:"work_in_progress"`
		URL            string               `json:"url"`
		Action         string               `json:"action"`
		OldRev         string               `json:"oldrev"`
		Assignee       gitlab.MergeAssignee `json:"assignee"`
	}

	glMergeEvent := gitlab.MergeEvent{
		ObjectAttributes: ObjectAttributes{
			TargetBranch: "foo",
		},
	}

	payload, _ := addBranch(&glMergeEvent)

	var glMergeResult glPullRequestPayload
	err := json.Unmarshal(payload, &glMergeResult)
	if err != nil {
		t.Errorf("Error during unmarshall of payload for gitlab merge request in TestAddBranchGitLabMergeRequest test")
	}

	if glMergeResult.ObjectAttributes.TargetBranch != "foo" {
		t.Errorf("Error - TargetBranch appears to have changed to %s, the Event should be unaltered", glMergeResult.ObjectAttributes.TargetBranch)
	}
	if glMergeResult.WebhookBranch != "foo" {
		t.Errorf("Error - Inccorect branch name set, expected foo, received %s", glMergeResult.WebhookBranch)
	}
}
