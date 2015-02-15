package gomatimport

import (
	"fmt"
	"os"
)

const (
	HeaderLength      = 128
	DescriptionLength = 116
)

type Mat struct {
	H    Header
	Data interface{}
}

func Import(filename string, dst interface{}) (Mat, error) {
	var out Mat

	f, err := os.Open(filename)
	if err != nil {
		return out, fmt.Errorf("Failed to open file %s, error message is : %s\n", err.Error())
	}
	defer func() {
		f.Close()
	}()

	// Read Header
	out.H = DecodeHeader(f)
	out.Data = DecodeData(f)

	return out, nil
}
