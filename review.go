package review

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Reviews is a map of reviews keyed by an integer line number
// TODO: eventually this should be a map of objects that contain the author and date as well, these would be obtained from the commit.
type Reviews map[int][]string

// Reviewer is a type that can list reviews, write reviews, and switch reviews
type Reviewer struct {
	Hash    func(path string) (string, error)
	Show    func(ref string, hash string) (string, error)
	GetRef  func(ref string) (string, error)
	AddNote func(ref string, hash string, message string) error
	Prune   func(ref string) error
}

func (r *Reviewer) Ref() string {
	ref, err := r.GetRef("REVIEW_HEAD")
	if err != nil {
		log.Fatal(err)
	}

	// these are both the same error case, good place for refactoring
	if ref == "" {
		err := errors.New("no current review. Create one with `git review switch <review name>`")
		log.Fatal(err)
	}
	ref = strings.TrimPrefix(ref, "refs/notes/")
	log.Printf("current review is %s\n", ref)
	return ref
}

func (r *Reviewer) List(path string) Reviews {
	reviews := make(Reviews)
	hash, err := r.Hash(path)
	if err != nil {
		log.Printf("[ERROR]: %s\n", err)
		return reviews
	}
	log.Printf("listing reviews for %s at %s\n", hash, path)
	note, err := r.Show(r.Ref(), hash)
	if err != nil {
		log.Printf("[ERROR]: %s\n", err)
		return reviews
	}
	//TODO: this is a good place to refactor
	noteLines := strings.Split(note, "\n")
	for _, line := range noteLines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ":")
		lineNumber, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Printf("[ERROR]: %s\n", err)
			continue
		}
		m, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			log.Printf("[ERROR]: %s\n", err)
			continue
		}
		message := string(m)
		log.Printf("reading note %s at %d\n", message, lineNumber)
		reviews[lineNumber] = append(reviews[lineNumber], message)
	}
	return reviews
}

func (r *Reviewer) Add(path string, line int, message string) error {
	// this is the point of transformation,
	// we can have a reviewer with a hash function of the current file,
	// and a hash tracked at the HEAD
	hash, err := r.Hash(path)
	if err != nil {
		return err
	}
	log.Printf("adding note to hash %s at %s\n", hash, path)

	//TODO: return an error as a warning if the hash is not tracked in git
	// This requires an error type that can be checked later

	//base64 encode the message so that it can contain formatting characters
	encodedMessage := base64.StdEncoding.EncodeToString([]byte(message))
	return r.AddNote(r.Ref(), hash, fmt.Sprintf("%d:%s", line, encodedMessage))
}

func (r *Reviewer) Clean() error {
	return r.Prune(r.Ref())
}
