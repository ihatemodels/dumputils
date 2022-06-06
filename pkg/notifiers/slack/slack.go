/*
 * Copyright (c) 2022. Dumputils Authors
 *
 * https://opensource.org/licenses/MIT
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ihatemodels/dumputils/pkg/notifiers"
	"github.com/rotisserie/eris"
	"io/ioutil"
	"net/http"
	"time"
)

const slackApiPostMessageURL = "https://slack.com/api/chat.postMessage"

// Notifier implements a notifiers.Notifier for Slack notifications.
type Notifier struct {
	client         *http.Client
	token          string
	channel        string
	defaultPayload request
}

// New returns slack.Notifier
func New(token string, channel string) (*Notifier, error) {
	return &Notifier{
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
		token:   token,
		channel: channel,
		defaultPayload: request{
			Channel: channel,
			Text:    "DumpUtils Notification",
		},
	}, nil
}

func (n *Notifier) Notify(in notifiers.Input) error {
	payload := n.defaultPayload

	header := succsefullHeader

	if in.State == notifiers.Failed {
		header = failedHeader
	}

	payloadFields := []fields{
		{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Type:*\n%s", in.InstanceType),
		},
		{
			Type: "mrkdwn",
			Text: fmt.Sprintf("*Created at:*\n%s", time.Now().Format("2006-01-02 15-01-05")),
		},
	}

	payload.Blocks = append(payload.Blocks,
		blocks{
			Type: "header",
			Text: &text{
				Type: "plain_text",
				Text: header,
			},
		}, blocks{
			Type: "section",
			Text: &text{
				Type: "mrkdwn",
				Text: fmt.Sprintf("Instance %s took %v to backup", in.InstanceName, in.Duration),
			},
		}, blocks{
			Type:   "section",
			Fields: &payloadFields,
		},
	)

	j, err := json.Marshal(payload)

	if err != nil {
		return eris.Wrap(err, "notifiers/slack: payload failed to serialise to json")
	}

	fmt.Println(string(j))
	req, err := http.NewRequest("POST", slackApiPostMessageURL, bytes.NewBuffer(j))

	if err != nil {
		return eris.Wrap(err, "notifiers/slack: post request failed to construct")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", n.token))

	resp, err := n.client.Do(req)

	if err != nil {
		return eris.Wrapf(err, "notifiers/slack: post request to %s failed", slackApiPostMessageURL)
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return eris.Wrapf(err, "notifiers/slack: reading response body failed")
	}

	var slackResponse response

	err = json.Unmarshal(respBody, &slackResponse)

	if err != nil {
		return eris.Wrap(err, "notifiers/slack: failed to serialize response")
	}

	if !slackResponse.Ok {
		return eris.New(fmt.Sprintf("notifiers/slack: notification failed to send. slack api response: %s", string(respBody)))
	}

	return nil
}
