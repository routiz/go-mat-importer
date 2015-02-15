package gomatimport

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

const (
	_              = iota
	TypeInt8       = iota
	TypeUint8      = iota
	TypeInt16      = iota
	TypeUint16     = iota
	TypeInt32      = iota
	TypeUint32     = iota
	TypeSingle     = iota
	TypeReserved1  = iota
	TypeDouble     = iota
	TypeReserved2  = iota
	TypeReserved3  = iota
	TypeInt64      = iota
	TypeUint64     = iota
	TypeMatrix     = iota
	TypeCompressed = iota
	TypeUtf8       = iota
	TypeUtf16      = iota
	TypeUtf32      = iota
)

type DataElement struct {
	Type          uint
	NumberOfBytes uint
	RawBytes      []byte
}

func SplitData(r io.Reader) []DataElement {
	delements := make([]DataElement, 0, 10)
	for {
		var delm DataElement

		// Read first 4 bytes and check if this element is small.
		tgbuffer1 := [4]byte{}
		readcnt, err := r.Read(tgbuffer1[:])
		if err != nil && err != io.EOF {
			panic(err)
		}
		if readcnt == 0 || err == io.EOF {
			break
		}
		// Check if first 2 bytes of the tag are zeroes.
		typ := binary.LittleEndian.Uint16(tgbuffer1[2:])
		// Number of bytes
		var nob uint
		if typ != 0 {
			// If first 2 bytes are not zeros, this is a small data element.
			nob = uint(binary.LittleEndian.Uint16(tgbuffer1[2:]))
		} else {
			// If first 2 bytes are zeros, this is a normal data element
			// which is not a small data element.

			// First 4 bytes are type.
			typ = binary.LittleEndian.Uint16(tgbuffer1[:])

			tgbuffer2 := [4]byte{}
			readcnt, err = r.Read(tgbuffer2[:])
			if err != nil && err != io.EOF {
				panic(err)
			}
			if readcnt == 0 || err == io.EOF {
				break
			}
			nob = uint(binary.LittleEndian.Uint32(tgbuffer2[:]))
		}
		dataBuffer := make([]byte, nob)
		readcnt, err = r.Read(dataBuffer) // maybe removed

		delm.Type = uint(typ)
		delm.NumberOfBytes = uint(nob)
		delm.RawBytes = dataBuffer

		delements = append(delements, delm)

		fmt.Println("type = ", typ)
		fmt.Println("number of bytes = ", nob)
		fmt.Println("read count = ", readcnt)
		// do something
	}

	return delements
}

func DecodeDataElement(delm DataElement) interface{} {
	fmt.Printf("Raw bytes size : %d\n", len(delm.RawBytes))
	reader := bytes.NewReader(delm.RawBytes)
	if reader == nil {
		return nil
	}

	switch delm.Type {

	case TypeCompressed:
		fmt.Printf("Data element is a compressed data\n")
		gzreader, err := zlib.NewReader(reader)
		if err != nil {
			log.Printf("ERROR - %s\n", err.Error())
			return nil
		}
		defer func() {
			gzreader.Close()
		}()

		rbuffer := make([]byte, 8*1024)
		wbuffer := make([]byte, 0)
		for {
			readcnt, err := gzreader.Read(rbuffer)
			if readcnt == 0 {
				break
			}
			if err != nil {
				log.Printf("ERROR - %s\n", err.Error())
				return nil
			}
			wbuffer = append(wbuffer, rbuffer[:readcnt]...)
		}
		fmt.Println("Header of decoded bytes :", wbuffer[:8])
		DecodeData(bytes.NewBuffer(wbuffer))
		break
	}
	return nil
}

func DecodeData(f io.Reader) interface{} {
	dataElements := SplitData(f)
	for _, delm := range dataElements {
		DecodeDataElement(delm)
	}
	var out Mat
	return out
}
