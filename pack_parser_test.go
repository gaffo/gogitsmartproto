package gogitsmartproto_test

import (
	"testing"
	"os"
	"github.com/gaffo/gogitsmartproto"
)


func Test_PackParse(t *testing.T) {
	f, err := os.Open("pack.dat")
	failIfError(t, err)
	pack, err := gogitsmartproto.ParsePack(f)
	failIfError(t, err)
	assertIEquals(t, 2, pack.Version)
	assertIEquals(t, 24, pack.CItems)
}
