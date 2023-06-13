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
	ref     string
}

func (r *Reviewer) Init() error {
	// if there is no current review at REVIEW_HEAD warn that one must be created with `git review switch`
	// can possibly extract this
	ref, err := r.GetRef("REVIEW_HEAD")
	if err != nil {
		log.Fatal(err)
	}

	// these are both the same error case, good place for refactoring
	if ref == "" {
		return errors.New("no current review. Create one with `git review switch <review name>`")
	}

	// hide this
	log.Println(r.Switch(ref))
	return nil
}

func (r *Reviewer) Switch(ref string) string {

	ref = strings.TrimPrefix(ref, "refs/notes/review/")

	r.ref = "review/" + ref
	return r.ref
}

func (r *Reviewer) List(path string) Reviews {
	reviews := make(Reviews)
	hash, err := r.Hash(path)
	if err != nil {
		return reviews
	}
	note, err := r.Show(r.ref, hash)
	if err != nil {
		return reviews
	}
	noteLines := strings.Split(note, "\n")
	for _, line := range noteLines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ":")
		lineNumber, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Println(err)
			continue
		}
		message, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			log.Println(err)
			continue
		}
		reviews[lineNumber] = append(reviews[lineNumber], string(message))
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

	//TODO: return an error as a warning if the hash is not tracked in git
	// This requires an error type that can be checked later

	//base64 encode the message so that it can contain formatting characters
	encodedMessage := base64.StdEncoding.EncodeToString([]byte(message))
	return r.AddNote(r.ref, hash, fmt.Sprintf("%d:%s", line, encodedMessage))
}
