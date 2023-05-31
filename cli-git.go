package review

import (
	"log"
	"os/exec"
	"strings"
)

// CliGit is a type that implements the NoteShower, NoteWriter interfaces
// using the git cli
type CliGit struct{}

// Show returns the notes for a given hash
func (g *CliGit) Show(ref string, hash string) (string, error) {
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
func (g *CliGit) AddNote(ref string, hash string, message string) error {
	c := exec.Command("git", "notes", "append", "-m", message, hash)
	c.Env = append(c.Env, "GIT_NOTES_REF=refs/notes/"+ref)
	o, err := c.CombinedOutput()
	log.Println(string(o))
	return err
}

// UpdateRef updates a symbolic ref to point to a given ref
func (g *CliGit) UpdateRef(ref string, target string) error {
	c := exec.Command("git", "symbolic-ref", ref, target)
	o, err := c.CombinedOutput()
	log.Println(string(o))
	return err
}

// GetRef gets a symbolic ref given a name
func (g *CliGit) GetRef(ref string) (string, error) {
	c := exec.Command("git", "symbolic-ref", ref)
	o, err := c.Output()
	target := strings.Trim(string(o), "\n")
	log.Println(target)
	return target, err
}
