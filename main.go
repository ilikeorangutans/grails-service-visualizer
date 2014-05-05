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
	"code.google.com/p/gographviz"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var IGNORED_DIRS []string = []string{".git", "target"}

func main() {

	dir := os.Args[1]
	absolutePath, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}

	outputName := filepath.Base(absolutePath)
	outputFilename := fmt.Sprintf("%s.dot", outputName)

	log.Printf("Scanning %s... ", absolutePath)
	files := buildListOfFiles(absolutePath)
	log.Printf("found %d files", len(files))

	graph := gographviz.NewGraph()
	graph.SetDir(true)
	graph.SetName("dependencyGraph")
	graph.SetStrict(true)

	graph.AddAttr("dependencyGraph", "splines", "\"ortho\"")
	graph.AddAttr("dependencyGraph", "ranksep", "\"2.0\"")
	graph.AddAttr("dependencyGraph", "concentrate", "true") // TODO: Not sure if that is actually valid for digraphs

	for _, f := range files {
		name := filepath.Base(f)

		isController := strings.Contains(name, "Controller")
		isService := strings.Contains(name, "Service")

		if !isController && !isService {
			continue
		}
		nameWithoutExtension := name[:strings.Index(name, ".groovy")]

		var params map[string]string = make(map[string]string)
		if isService {
			params["shape"] = "\"hexagon\""
			params["tooltip"] = fmt.Sprintf("\"From File %s\"", f)
		}

		if !graph.IsNode(nameWithoutExtension) {
			graph.AddNode("", nameWithoutExtension, params)
		}

		src, err := os.Open(f)

		if err != nil {
			log.Fatal(err)
		}
		reader := bufio.NewReader(src)

		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			if strings.Contains(line, "import") {
				continue
			}

			trimmed := strings.Trim(line, "\n\r\t")

			serviceDefinition := EndsWith(trimmed, "Service")

			if serviceDefinition {

				split := strings.Split(trimmed, " ")
				name := strings.Title(strings.Trim(split[1], " "))
				if len(name) == 0 {
					continue
				}
				var params map[string]string = make(map[string]string)
				if strings.Contains(name, "Service") {
					params["shape"] = "\"hexagon\""
					params["tooltip"] = fmt.Sprintf("\"From Source %s\"", trimmed)
				}

				if !graph.IsNode(name) {
					graph.AddNode("", name, params)
				}

				graph.AddEdge(nameWithoutExtension, "", name, "", true, nil)
			}
		}

	}

	log.Printf("Writing output to %s...", outputFilename)
	err = ioutil.WriteFile(outputFilename, []byte(graph.String()), 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func EndsWith(haystack, needle string) bool {
	if len(haystack) == 0 {
		return false
	}
	lastIndex := strings.LastIndex(haystack, needle)

	if lastIndex < 0 {
		return false
	}

	expectedIndex := len(haystack) - len(needle)
	return lastIndex == expectedIndex

}

func buildListOfFiles(dir string) []string {

	startDir, err := os.Open(dir)
	if err != nil {
		log.Fatal(err)
	}
	defer startDir.Close()

	var result []string

	fis, err := startDir.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}

	for _, fi := range fis {

		localName := filepath.Base(fi.Name())

		if fi.IsDir() {

			ignoreDir := false
			for _, ignore := range IGNORED_DIRS {
				if ignore == localName {
					ignoreDir = true
					break
				}
			}

			if ignoreDir {
				continue
			}

			s := filepath.Join(dir, fi.Name())
			if err != nil {
				log.Fatal(err)
			}

			result = append(result, buildListOfFiles(s)...)
		} else {
			if filepath.Ext(localName) == ".groovy" && !strings.Contains(localName, "Spec") && !strings.Contains(localName, "Test") {
				result = append(result, filepath.Join(dir, fi.Name()))
			}
		}

	}

	return result
}
