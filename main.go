package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/doraemonkeys/doraemon"
	"github.com/doraemonkeys/sedr/version"
)

const usage = `Usage: sedr <regex> <original> <substitution> <file>

Example:
sedr '(\d{4})-(\d{2})-(\d{2})' '$0' '$1-$3-$2' file.txt

Options:
  <regex>        Regular expression pattern
  <original>     Original text to replace (default: $0)
  <substitution> Replacement text
  <file>         File to process

Note: 
- $0 represents the entire matched string
- $+ can be used to represent a literal $`

type Config struct {
	Regex        string
	Original     string
	Substitution string
	Filename     string
}

func main() {
	config, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		fmt.Println(usage)
		fmt.Println()
		fmt.Println("Version:", version.Version)
		fmt.Println("Build Hash:", version.BuildHash)
		fmt.Println("Build Time:", version.BuildTime)
		os.Exit(1)
	}

	err = processFile(config)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("File processed successfully")
}

func parseArgs(args []string) (Config, error) {
	if len(args) < 3 {
		return Config{}, errors.New("insufficient arguments")
	}

	config := Config{
		Regex:        args[0],
		Original:     "$0",
		Substitution: args[len(args)-2],
		Filename:     args[len(args)-1],
	}

	if len(args) > 3 {
		config.Original = args[1]
	}

	return config, nil
}

func processFile(config Config) error {
	f, err := os.Open(config.Filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	stat, err := f.Stat()
	if err != nil {
		fmt.Println("Error getting file stat:", err)
		os.Exit(1)
	}
	fileMode := stat.Mode()
	content := make([]byte, stat.Size())
	_, err = f.Read(content)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	err = f.Close()
	if err != nil {
		fmt.Println("Error closing file:", err)
		os.Exit(1)
	}

	re, err := regexp.Compile(config.Regex)
	if err != nil {
		return fmt.Errorf("error compiling regex: %v", err)
	}

	captureGroupIndex, err := parseCaptureGroupIndex(config.Original)
	if err != nil {
		return fmt.Errorf("error parsing capture group index: %v", err)
	}

	matches := re.FindSubmatchIndex(content)
	if matches == nil {
		return errors.New("no matches found")
	}

	if len(matches)/2-1 < captureGroupIndex {
		return fmt.Errorf("capture group index out of range, max index is %d", len(matches)/2-1)
	}

	substitution := handleSubstitution(config.Substitution, re.FindSubmatch(content))
	start := matches[2*captureGroupIndex]
	end := matches[2*captureGroupIndex+1]
	err = doraemon.WriteFile(config.Filename, fileMode, content[:start], []byte(substitution), content[end:])
	if err != nil {
		return fmt.Errorf("error replacing content: %v", err)
	}
	return nil
}

func handleSubstitution(substitution string, matches [][]byte) string {
	for i := 1; i < len(matches); i++ {
		placeholder := fmt.Sprintf("$%d", i)
		substitution = strings.ReplaceAll(substitution, placeholder, string(matches[i]))
	}
	return strings.ReplaceAll(substitution, "$+", "$")
}

func parseCaptureGroupIndex(original string) (int, error) {
	if original == "$0" {
		return 0, nil
	}
	n, err := strconv.Atoi(original[1:])
	if err != nil {
		return 0, err
	}
	return n, nil
}
