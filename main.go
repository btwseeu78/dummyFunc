// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// [START functions_cloudevent_pubsub]

// Package helloworld provides a set of Cloud Functions samples.
package cloudfunction

import (
	"context"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/event"
	"log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/alecthomas/kingpin/v2"
)

func init() {
	functions.CloudEvent("HelloPubSub", helloPubSub)
}

// MessagePublishedData contains the full Pub/Sub message
// See the documentation for more details:
// https://cloud.google.com/eventarc/docs/cloudevents#pubsub
type MessagePublishedData struct {
	Message PubSubMessage
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

var (
	dynUrl    = kingpin.Flag("dynamic-url", "Dynatrace Url").Envar("DYNA_URL").String()
	projectId = kingpin.Flag("Project Id", "GCP Project ID").Envar("GKE_PROJECT_ID").String()
	apiToken  = kingpin.Flag("dynamic-token", "Dynatrace Token").Envar("DT_API_TOKEN").String()
)

// helloPubSub consumes a CloudEvent message and extracts the Pub/Sub message.
func helloPubSub(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	kingpin.Parse()
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}
	name := string(msg.Message.Data) // Automatically decoded from base64.
	if name == "" {
		name = "World"
	}
	log.Printf("Hello, %s!", name)
	log.Printf("dynaUrl: %s\n,projectId: %s\n,apiToken: %s\n", *dynUrl, *projectId, *apiToken)
	return nil
}

// [END functions_cloudevent_pubsub]
