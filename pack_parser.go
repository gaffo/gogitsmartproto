package gogitsmartproto

import (
	"io"
	"errors"
	"fmt"
	"encoding/binary"
	"compress/zlib"
	"bytes"
	"bufio"
)

type Pack struct {
	CItems  int
	Version int
}

type packObjectType uint8

const (
	_             packObjectType = iota
	OBJ_COMMIT
	OBJ_TREE
	OBJ_BLOB
	OBJ_TAG
	_
	OBJ_OFS_DELTA
	OBJ_REF_DELTA
)

func (this packObjectType) String() string {
	switch this {
	case OBJ_COMMIT:
		return "COMMIT"
	case OBJ_TREE:
		return "TREE"
	case OBJ_BLOB:
		return "BLOB"
	case OBJ_TAG:
		return "TAG"
	case OBJ_OFS_DELTA:
		return "OFS_DELTA"
	case OBJ_REF_DELTA:
		return "REF_DELTA"
	default:
		return "Unknown"
	}
}

func ParsePack(rdr io.ReadSeeker) (pack Pack, err error) {
	err = nil
	four := make([]byte, 4, 4)
	rdr.Read(four)
	header := string(four)
	if header != "PACK" {
		err = errors.New("Invalid Header")
		return
	}

	rdr.Read(four)
	pack.Version = int(binary.BigEndian.Uint32(four))

	rdr.Read(four)
	pack.CItems = int(binary.BigEndian.Uint32(four))

	for i := 0; i < pack.CItems; i++ {
		_bytes := make([]byte, 1)
		rdr.Read(_bytes)
		_byte := _bytes[0]

		objectTypeBits := ((_byte >> 4) & 7)
		MSB := (_byte & 128)
		chunk := uint(_byte) & 15
		objectSize := int(chunk)
		var shift uint = 4
		fmt.Println("i", chunk, shift, objectSize)
		for MSB > 0 {
			fmt.Println(MSB, ".")
			_bytes := make([]byte, 1)
			rdr.Read(_bytes)
			_byte := _bytes[0]

			MSB = (_byte & 128)

			chunk := uint(_byte) & 127
			newIncr := int(chunk << shift)
			objectSize += newIncr
			shift += 7
			fmt.Println("r", _bytes, chunk, newIncr, shift, objectSize)
		}

		before, err2 := rdr.Seek(0, io.SeekCurrent)
		objectType := packObjectType(objectTypeBits)

		fmt.Println(objectTypeBits, objectSize, objectType)

		zr, err2 := zlib.NewReader(bufio.NewReader(rdr))
		if err2 != nil {
			err = err2
			return
		}

		buf := new(bytes.Buffer)
		bytesRead, err2 := io.CopyN(buf, zr, int64(objectSize))
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}
		after, err2 := rdr.Seek(0, io.SeekCurrent)

		fmt.Println(bytesRead, objectSize, before, after)


		fmt.Println("======================")
		fmt.Println(buf.String())
		fmt.Println("---------------------")
	}

	return
}
