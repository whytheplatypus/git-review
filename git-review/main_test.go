package main

import (
	"testing"

	review "github.com/whytheplatypus/git-review"
)

func TestReview(t *testing.T) {
	// want to be able to test output or results
	// want to be able to affect the directory it oporates in
	// that could be the reviewer object
	Switch([]{"test-%d"}, nil)
	Add([]{"-m", "test"}, nil)
	Review([]{"test-%d"}, nil)


	type args struct {
		args     []string
		reviewer review.Reviewer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Review(tt.args.args, tt.args.reviewer); (err != nil) != tt.wantErr {
				t.Errorf("Review() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
