package review

import (
	"log"
	"strconv"
	"strings"
)

// Reviews is a map of reviews keyed by an integer line number
// TODO: eventually this should be a map of objects that contain the author and date as well, these would be obtained from the commit.
type Reviews map[int][]string

// Reviewer is a type that can list reviews, write reviews, and switch reviews
type Reviewer struct {
	Hasher
	NoteShower
	NoteWriter
	ref string
}

func (r *Reviewer) Switch(ref string) {
	r.ref = ref
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
		reviews[lineNumber] = append(reviews[lineNumber], strings.Join(parts[1:], ":"))
	}
	return reviews
}

func (r *Reviewer) Write(path string, line int, message string) error {
	hash, err := r.Hash(path)
	if err != nil {
		return err
	}
	return r.WriteNote(r.ref, hash, strconv.Itoa(line)+":"+message)
}
