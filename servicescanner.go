// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
)

// Used to find service definitions in the source code. Not the best regex but it should cover most cases.
var serviceRegex = regexp.MustCompile(`^[ \t]*[a-zA-Z]+[ \t]+([a-z][a-zA-Z0-9]+Service)[; \t]*`)

type Dependency struct {
	Source string
	Line   int
	Name   string
}

func (d Dependency) String() string {
	return fmt.Sprintf("%s:%d %s", d.Source, d.Line, d.Name)
}

// Scans the given reader for occurences of service definitions as defined by serviceRegex. Just a simple
// line-by-line scanner, not a parser so it doesn't recognize multiline comments.
func ScanForDependencies(filename string, reader *bufio.Reader) []Dependency {
	var deps []Dependency

	keepSearching := true
	lineNumber := 0

	for keepSearching {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			keepSearching = false
		} else if err != nil {
			log.Fatal(err)
		}

		lineNumber++

		if !serviceRegex.MatchString(line) {
			// Skip lines that are not service definitions.
			continue
		}

		submatch := serviceRegex.FindStringSubmatch(line)
		serviceName := submatch[1]
		capitalizedServiceName := strings.Title(serviceName)

		dep := Dependency{
			Name:   capitalizedServiceName,
			Line:   lineNumber,
			Source: filename,
		}

		deps = append(deps, dep)
	}

	return deps
}
