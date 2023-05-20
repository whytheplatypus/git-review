package review

import (
	"encoding/base64"
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
	hash, err := r.Hash(path)
	if err != nil {
		return err
	}
	//base64 encode the message so that it can contain formatting characters
	encodedMessage := base64.StdEncoding.EncodeToString([]byte(message))
	return r.AddNote(r.ref, hash, strconv.Itoa(line)+":"+encodedMessage)
}
