package gogitsmartproto

import (
	"io"
	"errors"
	"fmt"
	"encoding/binary"
	"compress/zlib"
	"bytes"
	"encoding/hex"
)

type byteReader struct {
	r io.ReadSeeker
}

func (this *byteReader) Read(p []byte) (n int, err error) {
	return this.r.Read(p)
}

func (this *byteReader) ReadByte() (byte, error) {
	buf := make([]byte, 1)
	_, err := this.r.Read(buf)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func (this *byteReader) Seek(offset int64, whence int) (int64, error) {
	return this.r.Seek(offset, whence)
}

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

func ParsePack(r io.ReadSeeker) (pack Pack, err error) {
	rdr := &byteReader{r: r}
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
		for MSB > 0 {
			_bytes := make([]byte, 1)
			rdr.Read(_bytes)
			_byte := _bytes[0]

			MSB = (_byte & 128)

			chunk := uint(_byte) & 127
			newIncr := int(chunk << shift)
			objectSize += newIncr
			shift += 7
		}

		objectType := packObjectType(objectTypeBits)

		fmt.Println(objectType)

		switch objectType {
		case OBJ_REF_DELTA:
			baseObjName := make([]byte, 20)
			rdr.Read(baseObjName)
			fmt.Println("baseObjName", string(baseObjName))
		}

		zr, err2 := zlib.NewReader(rdr)
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

		if bytesRead != int64(objectSize) {
			err = errors.New("Object size doesn't match decompressed size")
			return
		}

		fmt.Println("======================")
		switch objectType {
		case OBJ_BLOB:
			//fmt.Println(buf.String())
		case OBJ_TREE:
			b := buf.Bytes()
			start := 0
			for i := 0; i < len(b); i++ {
				v := b[i]
				if v != 0 {
					continue
				}

				pre := string(b[start:i])
				start = i + 21
				sha := hex.EncodeToString(b[i+1:start])
				fmt.Println(pre, sha)
			}
		default:
			//fmt.Println(buf.String())
		}

		fmt.Println("---------------------")
	}

	return
}

func nulary(buf []byte) []string {
	ret := make([]string, 0, 1024)
	start := 0
	for i, v := range buf {
		if v == 0 {
			str := string(buf[start:i])
			ret = append(ret, str)
			start = i + 1
		}
	}
	return ret
}

func nulidx(ary []byte) int {
	for i, v := range ary {
		if v == 0 {
			return i
		}
	}
	return -1
}