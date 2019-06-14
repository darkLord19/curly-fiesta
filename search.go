package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	patternColor    = "\033[31m"
	filenameColor   = "\033[35m"
	lineNumberColor = "\033[34m"
	resetColor      = "\033[0m"

	EXIT_FAILURE = 1
)

var (
	showLineNum    bool
	showColoredOut bool
	fileCount      int
	stdOutWriter   io.Writer
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func printUsage() {
	val := "usage: grep [-n] [-c/--colored] [-h/--help] [pattern] [file ...]"
	fmt.Fprintf(stdOutWriter, "%s\n", val)
}

func printOut(filename string, matchedLine string, lnum string) {
	if fileCount > 1 {
		fmt.Fprintf(stdOutWriter, "%s:", filename)
	}
	if showLineNum {
		fmt.Fprintf(stdOutWriter, "%s: %s\n", lnum, matchedLine)
	} else {
		fmt.Fprintf(stdOutWriter, "%s\n", matchedLine)
	}
}

func printColoredOut(filename string, matchedLine string, lnum string) {
	filename = filenameColor + filename + resetColor
	lnum = lineNumberColor + lnum + resetColor
	printOut(filename, matchedLine, lnum)
}

func init() {
	flag.BoolVar(&showLineNum, "n", false, "Flag to specify if you want to print line numbers or not")
	flag.BoolVar(&showColoredOut, "colored", false, "Flag to specify if you want colored output or not")
	flag.BoolVar(&showColoredOut, "c", false, "Flag to specify if you want colored output or not (shorthand)")
	flag.Parse()
	stdOutWriter = bufio.NewWriter(os.Stdout)
}

func main() {

	args := flag.Args()

	if len(args) < 2 {
		printUsage()
		os.Exit(EXIT_FAILURE)
	}

	searchTerm := args[0]
	filenames := args[1:]
	fileCount = len(filenames)

	for i := range filenames {
		ln := 0
		file, err := os.Open(filenames[i])
		check(err)
		defer file.Close()

		scanner := bufio.NewScanner(file)

		//Read each line one by one of file
		for scanner.Scan() {
			line := scanner.Text()
			// Check if line contains given search string
			if strings.Contains(line, searchTerm) {
				if showColoredOut {
					printColoredOut(filenames[i], line, strconv.Itoa(ln))
				} else {
					printOut(filenames[i], line, strconv.Itoa(ln))
				}
			}
			ln++
		}
	}
}
