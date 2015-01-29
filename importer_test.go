package gomatimport

import (
	"testing"
)

func TestImport(t *testing.T) {
	_, err := Import("index.mat", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}
}
