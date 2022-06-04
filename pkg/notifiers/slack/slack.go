/*
 * Copyright (c) 2022.  Dumputils Authors
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
	"github.com/ihatemodels/dumputils/pkg/notifiers"
	"net/http"
	"time"
)

// Notifier implements a notifiers.Notifier for Slack notifications.
type Notifier struct {
	client         *http.Client
	token          string
	channel        string
	defaultPayload request
}

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

func (n *Notifier) Notify(state notifiers.NotificationState) error {
	payload := n.defaultPayload

	payload.Blocks = append(payload.Blocks, block{
		Type: "header",
	})
	return nil
}
