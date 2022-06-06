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

const (
	failedHeader     = ":relieved: DumpUtils -> backup failed to create"
	succsefullHeader = ":relieved: DumpUtils -> backup created successfully"
)

type request struct {
	Channel string   `json:"channel"`
	Text    string   `json:"text"`
	Blocks  []blocks `json:"blocks"`
}
type text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
type fields struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
type blocks struct {
	Type   string    `json:"type"`
	Text   *text     `json:"text,omitempty"`
	Fields *[]fields `json:"fields,omitempty"`
}

type response struct {
	Ok               bool     `json:"ok"`
	Error            string   `json:"error"`
	Errors           []string `json:"errors"`
	Warning          string   `json:"warning"`
	ResponseMetadata struct {
		Messages []string `json:"messages"`
		Warnings []string `json:"warnings"`
	} `json:"response_metadata"`
}
