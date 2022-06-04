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

package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"sort"
	"sync"
)

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
)

const (
	StackTraceField = "stacktrace"
	ErrorField      = "error"
	TimeField       = "time"
)

var (
	consoleBufPool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 100))
		},
	}
)

type consoleWriter struct {
	level zerolog.Level
	out   io.Writer
}

func colorize(color int, v string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, v)
}

// Write -
// This code is copied from zerolog.ConsoleWriter
func (w *consoleWriter) Write(p []byte) (n int, err error) {
	var buf = consoleBufPool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		consoleBufPool.Put(buf)
	}()

	var evt map[string]interface{}
	d := json.NewDecoder(bytes.NewReader(p))
	d.UseNumber()
	err = d.Decode(&evt)
	if err != nil {
		return n, fmt.Errorf("cannot decode event: %s", err)
	}
	keys := make([]string, 0, len(evt))

	for field := range evt {
		switch field {
		case StackTraceField, ErrorField, TimeField:
			continue
		}
		keys = append(keys, field)
	}

	sort.Strings(keys)

	buf.WriteString(colorize(colorGreen, "time"))
	buf.WriteString(fmt.Sprintf("=%v", evt[TimeField]))
	buf.WriteByte(' ')

	for i, k := range keys {
		buf.WriteString(colorize(colorYellow, k))
		buf.WriteString(fmt.Sprintf("=%v", evt[k]))
		if i != len(keys)-1 { // Skip space for last part
			buf.WriteByte(' ')
		}
	}

	if _, ok := evt[ErrorField]; ok {
		buf.WriteByte(' ')
		buf.WriteString(colorize(colorRed, "error"))
		buf.WriteString(fmt.Sprintf("=%v", evt[ErrorField]))
	}

	if _, ok := evt[StackTraceField]; ok {
		buf.WriteByte('\n')
		buf.WriteString(colorize(colorRed, "StackTrace -----> "))
		buf.WriteByte('\n')
		buf.WriteString(fmt.Sprintf("%v", evt[StackTraceField]))
	}

	buf.WriteByte('\n')
	_, err = buf.WriteTo(w.out)
	return len(p), err
}

func (w *consoleWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level < w.level {
		return len(p), nil
	}
	return w.Write(p)
}
