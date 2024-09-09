package main_test

import (
	"testing"

	main "github.com/BHillGH/devops-demo"
)

func TestMain(t *testing.T) {
	rec, err := main.GetContentFromCSV()
	if err != nil {
		t.Error(err)
	}
	if len(rec) < 7 {
		t.Errorf("Incorrect Length")
	}
}
