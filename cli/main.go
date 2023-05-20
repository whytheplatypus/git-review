package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"strconv"

	review "github.com/whytheplatypus/git-review"
)

/**
* lists reviews: git review
* switches review: git review switch <review name>
* starts a new review: git review switch -c <review name>
* lists reivews for a directory: git review <dir>
* lists reviews for a file: git review <file>
* shows reviews for a line: git review <file> <line>
* adds review for a line: git review -m "message" <file> <line>
* prunes notes: git review prune
* opens default editor to add a review for a line: git review <file> <line> -e
* opens a specified editor to add a review for a line: git review <file> <line> -e <editor>
 */

var cliGit = review.CliGit{}
var goGit = review.Git{}
var reviewer = review.Reviewer{
	Hasher:     &goGit,
	NoteShower: &cliGit,
	NoteWriter: &cliGit,
}

func init() {
	log.SetFlags(log.Lshortfile)
}

type command interface {
	Parse() error
	Execute() error
}

func Review(args []string) error {
	f := flag.NewFlagSet("review", flag.ExitOnError)
	var message string
	f.StringVar(&message, "m", "", "message")
	f.Parse(args)
	file := f.Arg(0)
	line := f.Arg(1)
	log.Println(file, line, message)
	if message != "" {
		log.Println(message)

		// Check that both file and line are present
		if file == "" || line == "" {
			return errors.New("file and line must be specified")
		}

		// Check that line is a number
		lineNumber, err := strconv.Atoi(line)
		if err != nil {
			return err
		}

		// Add the note
		return reviewer.Add(file, lineNumber, message)
	}

	reviews := reviewer.List(file)

	if line != "" {
		lineNumber, err := strconv.Atoi(line)
		if err != nil {
			return err
		}
		log.Println(reviews[lineNumber])
		return nil
	}

	for _, review := range reviews {
		log.Println(review)
	}

	return nil
}

func main() {
	// Print the arguments

	reviewer.Switch("review")

	log.Println(Review(os.Args[1:]))
}
