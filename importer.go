package gomatimport

import (
	"fmt"
	"os"
)

const (
	HeaderLength      = 128
	DescriptionLength = 116
)

type Header struct {
	Description string
	SubsysData  interface{}
	Version     interface{}
	Endian      interface{}
}

type ParsedData struct {
	H    Header
	Data interface{}
}

func NewHeaderBytes(hdbuffer []byte) Header {
	var out Header

	out.Description = string(hdbuffer[:DescriptionLength])

	return out
}

func Import(filename string, dst interface{}) (ParsedData, error) {
	var out ParsedData

	f, err := os.Open(filename)
	if err != nil {
		return out, fmt.Errorf("Failed to open file %s, error message is : %s\n", err.Error())
	}
	defer func() {
		f.Close()
	}()

	// Read Header
	hdbuffer := make([]byte, HeaderLength)
	f.Read(hdbuffer)
	out.H = NewHeaderBytes(hdbuffer)

	return out, nil
}
