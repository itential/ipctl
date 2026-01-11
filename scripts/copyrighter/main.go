// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const header = `// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

`

var checkOnly bool

func main() {
	flag.BoolVar(&checkOnly, "check", false, "Only check files")
	flag.Parse()

	filepath.Walk(".", processFile)
}

func processFile(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() && info.Name() == "vendor" {
		return filepath.SkipDir
	}

	if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
		fmt.Printf("checking %s\n", path)
		err = checkAndFixFile(path)
		if err != nil {
			fmt.Printf("Error processing file %s: %v\n", path, err)
		}
	}

	return nil
}

func checkAndFixFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if checkOnly {
		if !strings.HasPrefix(string(content), header) {
			log.Fatalf("missing or invalid source header: %s\n", path)
		}
	} else {
		reader := bufio.NewReader(bytes.NewReader(content))
		var restOfFile bytes.Buffer

		inCopyrightBlock := false
		foundPackage := false

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err.Error() == "EOF" {
					restOfFile.WriteString(line)
					break
				}
				return err
			}
			if strings.HasPrefix(strings.TrimSpace(line), "package ") {
				foundPackage = true
				restOfFile.WriteString(line)
				break
			}
			if strings.HasPrefix(strings.TrimSpace(line), "// Copyright") {
				inCopyrightBlock = true
			}
			if inCopyrightBlock {
				if strings.TrimSpace(line) == "" {
					inCopyrightBlock = false
				}
			} else {
				restOfFile.WriteString(line)
			}
		}

		// Read the rest of the file in so we can write back
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err.Error() == "EOF" {
					restOfFile.WriteString(line)
					break
				}
				return err
			}
			restOfFile.WriteString(line)
		}

		if foundPackage {
			newContent := append([]byte(header), restOfFile.Bytes()...)
			err = os.WriteFile(path, newContent, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
