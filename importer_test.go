package gomatimport

import (
	"fmt"
	"testing"
)

func TestImport(t *testing.T) {
	data, err := Import("index.mat", nil)
	if err != nil {
		t.Error(err.Error())
		return
	}

	fmt.Println("Description:")
	fmt.Print(data.H.Description, "\n")
	fmt.Println("Version:", data.H.Version)
	fmt.Println("Endian:", data.H.Endian)
}
