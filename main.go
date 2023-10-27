package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var re = regexp.MustCompile(`[^a-zA-z]t\(\s*'([\s\S]+?)',?\s*\)`)

func main() {
	if len(os.Args[1:]) == 0 {
		log.Printf("path is not specifed")
		os.Exit(1)
	}
	m := map[string]string{}
	for _, p := range os.Args[1:] {
		log.Printf("walk path %s", p)
		err := filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			log.Printf("read file: %s", path)
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			matches := re.FindAllStringSubmatch(string(b), -1)
			for _, match := range matches {
				_, ok := m[match[1]]
				if !ok {
					m[match[1]] = match[1]
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("fail to walk path: %s %v", p, err)
		}
	}

	log.Println("completed!")
	fmt.Println()

	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	err := e.Encode(m)
	if err != nil {
		log.Printf("fail to encode json: %v", err)
	}
}
