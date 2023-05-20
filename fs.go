package review

import "os"

// Fs is a type that implements the FileLister interface
// using the os package
type Fs struct{}

// List returns the files in a commit
func (f *Fs) List(ref string) ([]string, error) {
	fileInfos, err := os.ReadDir(ref)
	if err != nil {
		return nil, err
	}
	filePaths := []string{}
	for _, fileInfo := range fileInfos {
		// skip directories
		if fileInfo.IsDir() {
			continue
		}
		filePaths = append(filePaths, fileInfo.Name())
	}
	return filePaths, nil
}
