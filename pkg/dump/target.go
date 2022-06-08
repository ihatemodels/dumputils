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

package dump

type Input struct {
	DestinationDir string
	Binaries       BinPaths
}

type BinPaths struct {
	PgTools10 string
	PgTools11 string
	PgTools12 string
	PgTools13 string
	PgTools14 string

	MySQLDump string
	MongoDump string
}

func NewPaths() BinPaths {
	return BinPaths{
		PgTools10: "/usr/lib/postgresql/10/bin/",
		PgTools11: "/usr/lib/postgresql/11/bin/",
		PgTools12: "/usr/lib/postgresql/12/bin/",
		PgTools13: "/usr/lib/postgresql/13/bin/",
		PgTools14: "/usr/lib/postgresql/14/bin/",

		MySQLDump: "/usr/bin/mysqldump",
		MongoDump: "/usr/bin/mongodump",
	}
}

type Target interface {
	Dump() error
	Probe() error
}
