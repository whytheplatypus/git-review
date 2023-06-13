package cligit

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// Show returns the notes for a given hash
func Show(ref string, hash string) (string, error) {
	c := exec.Command("git", "notes", "show", hash)
	c.Env = append(c.Env, "GIT_NOTES_REF=refs/notes/"+ref)
	o, err := c.Output()
	if err != nil {
		return string(o), err
	}
	log.Println(string(o))
	return string(o), nil
}

// WriteNote appends a note to a given hash
func AddNote(ref string, hash string, message string) error {
	c := exec.Command("git", "notes", "append", "-m", message, hash)
	c.Env = append(c.Env, "GIT_NOTES_REF=refs/notes/"+ref)
	o, err := c.CombinedOutput()
	log.Println(string(o))
	return err
}

func TrackedHash(path string) (string, error) {
	c := exec.Command("git", "rev-parse", fmt.Sprintf("HEAD:%s", path))
	o, err := c.CombinedOutput()
	if err != nil {
		return string(o), err
	}
	log.Println(string(o))
	return string(o), nil
}

// UpdateRef updates a symbolic ref to point to a given ref
func UpdateRef(ref string, target string) error {
	c := exec.Command("git", "symbolic-ref", ref, target)
	o, err := c.CombinedOutput()
	log.Println(string(o))
	return err
}

// GetRef gets a symbolic ref given a name
func GetRef(ref string) (string, error) {
	c := exec.Command("git", "symbolic-ref", ref)
	o, err := c.Output()
	target := strings.Trim(string(o), "\n")
	log.Println(target)
	return target, err
}
