// Copyright 2018 Palantir Technologies, Inc.
// Copyright 2020 G-Research Limited
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"context"
	"encoding/json"

	"github.com/google/go-github/v53/github"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/palantir/policy-bot/policy/common"
	"github.com/pkg/errors"

	"github.com/G-Research/tfe-plan-bot/pull"
)

type PullRequestReview struct {
	Base
}

func (h *PullRequestReview) Handles() []string { return []string{"pull_request_review"} }

// Handle pull_request_review
// https://developer.github.com/v3/activity/events/types/#pullrequestreviewevent
func (h *PullRequestReview) Handle(ctx context.Context, eventType, deliveryID string, payload []byte) error {
	var event github.PullRequestReviewEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return errors.Wrap(err, "failed to parse pull request review event payload")
	}

	// Ignore events triggered by policy-bot (e.g. for dismissing stale reviews)
	if event.GetSender().GetLogin() == h.AppName+"[bot]" {
		return nil
	}

	installationID := githubapp.GetInstallationIDFromEvent(&event)
	ctx, _ = h.PreparePRContext(ctx, installationID, event.GetPullRequest())

	return h.Evaluate(ctx, installationID, common.TriggerReview, pull.Locator{
		Owner:  event.GetRepo().GetOwner().GetLogin(),
		Repo:   event.GetRepo().GetName(),
		Number: event.GetPullRequest().GetNumber(),
		Value:  event.GetPullRequest(),
	})
}
