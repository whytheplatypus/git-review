package review

import (
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Git is a type that implements the Hasher and FileLister interfaces
// using the go-git library
type Git struct{}

// Hash returns the hash of a file at a given ref
func (g *Git) Hash(path string) (string, error) {
	//read the bytes of the file
	bytes, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	// hash the contents of the file
	hash := plumbing.ComputeHash(plumbing.BlobObject, bytes)
	return hash.String(), nil
}

// List returns the files in a commit
func (g *Git) List(ref string) ([]string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal(err)
	}
	// if ref is the empty string then use HEAD
	if ref == "" {
		ref = "HEAD"
	}

	reference, err := repo.Reference(plumbing.ReferenceName(ref), true)
	if err != nil {
		return nil, err
	}
	commit, err := repo.CommitObject(reference.Hash())
	if err != nil {
		return nil, err
	}

	files, err := commit.Files()
	if err != nil {
		return nil, err
	}

	var fileNames []string
	for file, err := files.Next(); err == nil; file, err = files.Next() {
		fileNames = append(fileNames, file.Name)
	}
	return fileNames, nil
}
