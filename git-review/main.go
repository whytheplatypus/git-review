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
			if fmt.Sprintf("review/%s", ref) == reviewer.Ref() {
				fmt.Printf("* ")
			}
			fmt.Println(ref)
		}
		return nil
	}
	//TODO: check ref structure with path

	// update the symbolic ref "REVIEW_HEAD" to point to the specified ref
	cligit.UpdateRef("REVIEW_HEAD", "refs/notes/review/"+ref)
	return nil
}

func Review(args []string, reviewer review.Reviewer) error {
	f := flag.NewFlagSet("review", flag.ExitOnError)
	line := f.Int("l", -1, "Line number to show review for")
	f.Parse(args)

	file := f.Arg(0)
	if file == "" {
		// list directory
		return errors.New("file must be specified")
	}

	reviews := reviewer.List(file)

	if *line > -1 {
		fmt.Printf("%d - %s", *line, reviews[*line])
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

func Prune(args []string, reviewer review.Reviewer) error {
	f := flag.NewFlagSet("prune", flag.ExitOnError)
	f.Parse(args)

	//FIXME: aweful back and forth naming
	return reviewer.Clean()
}

func Noop(msg string) func([]string, review.Reviewer) error {
	return func(args []string, reviewer review.Reviewer) error {
		fmt.Printf("%s is not yet implemented.\n", msg)
		return nil
	}
}

var commands = map[string]command{
	"switch": Switch,
	"add":    Add,
	"list":   Review, // default command
	"prune":  Prune,
	"lsp":    Noop("lsp"),
	"server": Noop("server"),
}

func main() {
	var verbose bool
	var tracked bool
	flag.BoolVar(&verbose, "v", false, "Show verbose logging")
	flag.BoolVar(&tracked, "t", true, "compute refs from current file contents")
	flag.Parse()

	if verbose {
		log.SetOutput(os.Stderr)
	}

	var reviewer = review.Reviewer{
		Hash:    gogit.Hash,
		Show:    cligit.Show,
		AddNote: cligit.AddNote,
		GetRef:  cligit.GetRef,
		Prune:   cligit.Prune,
	}

	if tracked {
		log.Println("using the hash from git history")
		reviewer.Hash = cligit.TrackedHash
	}

	args := flag.Args()[1:]
	command, ok := commands[flag.Arg(0)]
	if !ok {
		log.Fatalf("Unknown command %s", flag.Arg(0))
	}

	if err := command(args, reviewer); err != nil {
		log.Fatal(err)
	}
}
