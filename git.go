package review

// Hasher is an interface for getting the git hash of a file
type Hasher interface {
	Hash(path string) (string, error)
}

// NoteShower is an interface for showing git notes for a hash
type NoteShower interface {
	Show(ref string, hash string) (string, error)
}

// NoteWriter is an interface for appending to git notes for a hash
type NoteWriter interface {
	WriteNote(ref string, hash string, message string) error
}

// FileLister is an interface for listing files in a commit
type FileLister interface {
	List(ref string) ([]string, error)
}
