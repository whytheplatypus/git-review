package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	review "github.com/whytheplatypus/git-review"
)

/**
* lists reviews: git review
* switches review: git review switch <review name>
* lists reivews for a directory: git review <dir>
* lists reviews for a file: git review <file>
* shows reviews for a line: git review <file> <line>
* adds review for a line: git review <file> <line> -m "message"
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
	log.SetOutput(io.Discard)
}

type command interface {
	Parse() error
	Execute() error
}

func Switch(args []string) error {
	f := flag.NewFlagSet("switch", flag.ExitOnError)
	f.Parse(args)
	ref := f.Arg(0)

	if ref == "" {
		refs, err := goGit.ListRefs("refs/notes/review/")
		if err != nil {
			return err
		}
		for _, ref := range refs {
			fmt.Println(ref)
		}
		return nil
	}

	r := reviewer.Switch(ref)

	// update the symbolic ref "REVIEW_HEAD" to point to the specified ref
	cliGit.UpdateRef("REVIEW_HEAD", "refs/notes/"+r)
	return nil
}

func Review(args []string) error {
	f := flag.NewFlagSet("review", flag.ExitOnError)
	var verbose bool
	f.BoolVar(&verbose, "v", false, "Show verbose logging")
	f.Parse(args)

	if verbose {
		log.SetOutput(os.Stderr)
	}

	file := f.Arg(0)
	line := f.Arg(1)

	if len(f.Args()) > 2 {
		mf := flag.NewFlagSet("message", flag.ExitOnError)
		var message string
		mf.StringVar(&message, "m", "", "Message to add to the review")
		mf.Parse(f.Args()[2:])

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
	}

	reviews := reviewer.List(file)

	if line != "" {
		lineNumber, err := strconv.Atoi(line)
		if err != nil {
			return err
		}
		fmt.Println(reviews[lineNumber])
		return nil
	}

	fmt.Println(reviews)

	return nil
}

func main() {
	// if there is no current review at REVIEW_HEAD warn that one must be created with `git review switch`
	ref, err := cliGit.GetRef("REVIEW_HEAD")
	if err != nil {
		log.Fatal(err)
	}

	if ref == "" {
		log.Println("No current review. Create one with `git review switch <review name>`")
		return
	}

	if len(os.Args) < 2 {
		fmt.Println(ref)
		return
	}

	if os.Args[1] == "switch" {
		log.Println(Switch(os.Args[2:]))
		return
	}

	log.Println(ref)

	reviewer.Switch(ref)

	log.Println(Review(os.Args[1:]))
}
