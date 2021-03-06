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

package utils

import (
	"compress/gzip"
	"fmt"
	"github.com/rotisserie/eris"
	"io"
	"os"
	"path/filepath"
)

// Contains reports if given element of type T is found in slice of type []T
func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func Gzip(source string) error {
	reader, err := os.Open(source)

	if err != nil {
		return eris.Wrapf(err, "can not read %s file", source)
	}

	filename := filepath.Base(source)
	target := fmt.Sprintf("%s.gz", source)
	writer, err := os.Create(target)

	if err != nil {
		return eris.Wrapf(err, "can not create %s target file", target)
	}

	defer writer.Close()
	archiver := gzip.NewWriter(writer)
	archiver.Name = filename
	defer archiver.Close()

	_, err = io.Copy(archiver, reader)

	if err != nil {
		return eris.Wrapf(err, "can not write data while archiving to %s destination file", target)
	}

	return nil
}
