package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type Path struct {
	name string
	ref  string
}

type ByName []Path

func (a ByName) Len() int      { return len(a) }
func (a ByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool {
	splitsi := strings.Split(a[i].name[1:], "/")
	splitsj := strings.Split(a[j].name[1:], "/")

	min_index := min(len(splitsi), len(splitsj))

	for k := 0; k < min_index; k++ {
		if splitsi[k] != splitsj[k] {
			return splitsi[k] < splitsj[k]
		}
	}

	return len(splitsi) < len(splitsj)
}

func (self *Path) AsRef() string {
	return fmt.Sprintf("\"%v\": {\"$ref\": \"%v\"}", self.name, self.ref)
}

func walkDir(path string, path2name func(string) string) []Path {
	var out []Path
	entries, err := os.ReadDir(path)

	if err != nil {
		log.Fatalf("ReadDir error: %v\n", err)
	}

	for _, e := range entries {
		path := fmt.Sprintf("%v/%v", path, e.Name())

		if !e.IsDir() {
			out = append(out, Path{path2name(path), path})
		} else {
			out = append(out, walkDir(path, path2name)...)
		}
	}

	return out
}

func AsString(paths []Path) string {
	out := "{"

	limit := len(paths)

	sort.Sort(ByName(paths))

	for i, path := range paths {
		fmt.Fprintf(os.Stderr, "PATH: '%v'\n", path.name)
		out += path.AsRef()

		if i+1 != limit {
			out += ","
		}
	}

	out += "}"

	return out
}

func createPaths() []Path {
	if err := os.Chdir("./openapi/"); err != nil {
		log.Fatalf("CD error: %v\n", err)
	}

	return walkDir("./paths", func(path string) string {
		return path[7 : len(path)-5]
	})
}

func createSchemas() []Path {
	if err := os.Chdir("./openapi/components/"); err != nil {
		log.Fatalf("CD error: %v\n", err)
	}

	return walkDir("./schemas", func(path string) string {
		return path[10 : len(path)-5]
	})
}

func main() {
	functions := make(map[string]func() []Path)

	functions["paths"] = createPaths
	functions["schemas"] = createSchemas

	names := make([]string, len(functions))
	for k := range functions {
		names = append(names, k)
	}

	if 2 != len(os.Args) {
		names := make([]string, len(functions))

		for k := range functions {
			names = append(names, k)
		}

		log.Fatalf("Wrong arguments: %v {what: %v}", os.Args[0], names)
	}

	if f := functions[os.Args[1]]; nil == f {
		log.Fatalf("Wrong function: %v not in %v", os.Args[1], names)
	} else {
		fmt.Println(AsString(f()))
	}
}

