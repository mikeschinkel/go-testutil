package testutil

import (
	"testing"

	"github.com/mikeschinkel/go-dt"
)

func LoadFile(t *testing.T, file dt.Filepath, mustLoad bool) (data []byte) {
	var err error

	data, err = file.ReadFile()
	if err != nil && mustLoad {
		t.Fatal(err)
	}

	return data
}
