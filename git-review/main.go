package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	review "github.com/whytheplatypus/git-review"
	cligit "github.com/whytheplatypus/git-review/cli-git"
	gogit "github.com/whytheplatypus/git-review/go-git"
)

/**
* lists reviews: git review
* switches review: git review switch <review name>
* lists reivews for a directory: git review <dir>
* lists reviews for a file: git review <file>
* shows reviews for a line: git review -l <line> <file>
* adds review for a line: git review add -m "message" -l <line> <file>
* prunes notes: git review prune
* opens default editor to add a review for a line: git review <file> <line> -e
* opens a specified editor to add a review for a line: git review <file> <line> -e <editor>
 */

func init() {
	log.SetFlags(log.Lshortfile)
	log.SetOutput(io.Discard)
}

type command func(args []string, reviewer review.Reviewer) error

func Switch(args []string, reviewer review.Reviewer) error {
	f := flag.NewFlagSet("switch", flag.ExitOnError)
	f.Parse(args)
	ref := f.Arg(0)

	if ref == "" {
		refs, err := gogit.ListRefs("refs/notes/review/")
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
	cligit.UpdateRef("REVIEW_HEAD", "refs/notes/"+r)
	return nil
}

func Review(args []string, reviewer review.Reviewer) error {
	f := flag.NewFlagSet("review", flag.ExitOnError)
	line := f.Int("l", -1, "Line number to show review for")
	f.Parse(args)

	file := f.Arg(0)
	if file == "" {
		return errors.New("file must be specified")
	}

	reviews := reviewer.List(file)

	if *line > -1 {
		fmt.Println(reviews[*line])
		return nil
	}

	fmt.Println(reviews)

	return nil
}

func Add(args []string, reviewer review.Reviewer) error {
	f := flag.NewFlagSet("add", flag.ExitOnError)
	message := f.String("m", "", "Message to add to review")
	line := f.Int("l", -1, "Line number to add review for")
	f.Parse(args)

	file := f.Arg(0)
	if file == "" {
		return errors.New("file must be specified")
	}

	if *message == "" {
		return errors.New("message must be specified")
	}

	return reviewer.Add(file, *line, *message)
}

var commands = map[string]command{
	"switch": Switch,
	"add":    Add,
	// I think I need a flag to either show the results for a file as is or from history
	"list": Review, // default command
}

func main() {
	var verbose bool
	var tracked bool
	flag.BoolVar(&verbose, "v", false, "Show verbose logging")
	flag.BoolVar(&tracked, "t", false, "Use tracked refs for files")
	flag.Parse()

	if verbose {
		log.SetOutput(os.Stderr)
	}

	var reviewer = review.Reviewer{
		Hash:    gogit.Hash,
		Show:    cligit.Show,
		AddNote: cligit.AddNote,
		GetRef:  cligit.GetRef,
	}

	if tracked {
		log.Println("tracking")
		reviewer = review.Reviewer{
			Hash:    cligit.TrackedHash,
			Show:    cligit.Show,
			AddNote: cligit.AddNote,
			GetRef:  cligit.GetRef,
		}
	}

	args := flag.Args()[1:]
	command, ok := commands[flag.Arg(0)]
	if !ok {
		log.Fatalf("Unknown command %s", flag.Arg(0))
	}

	if err := reviewer.Init(); err != nil {
		log.Fatal(err)
	}

	if err := command(args, reviewer); err != nil {
		log.Fatal(err)
	}
}
