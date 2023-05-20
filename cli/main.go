package main

import (
	"flag"
	"log"
	"os"

	review "github.com/whytheplatypus/git-review"
)

/**
* lists reviews: git review
* switches review: git review switch <review name>
* starts a new review: git review switch -c <review name>
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
}

func main() {
	// Print the arguments
	log.Println(os.Args)
	flag.Parse()
	log.Println(flag.Args())
	if len(flag.Args()) > 0 {
		args := flag.Args()
		f := flag.NewFlagSet(args[0], flag.ExitOnError)
		f.Parse(args[1:])
		log.Println(f.Args())
	}

	reviewer.Switch("review")

	// List reviews for the current directory

	fs := review.Fs{}
	files, err := fs.List(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		log.Println(file)
		reviews := reviewer.List(file)

		for _, review := range reviews {
			log.Println(review)
		}
	}

}

type command interface {
	Parse() error
	Execute() error
}
