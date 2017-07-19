package gogitsmartproto_test

import (
	"testing"
	"os"
	"github.com/gaffo/gogitsmartproto"
	"reflect"
)

func failIfError(t* testing.T, err error) {
	if err != nil {
		t.Fatal("Unable to open file", err)
	}
}

func Test_ParseWithResults(t *testing.T) {
	f, err := os.Open("rails_refs.txt")
	failIfError(t, err)
	results, err := gogitsmartproto.ParseRefs(f)
	failIfError(t, err)
	expected := 31073
	assertLength(t, expected, results.Refs)

	ref := results.Refs[0]
	assertEquals(t, "3802de4a769092a4b6477e9b5ec0636938c5a957", ref.Checksum)
	assertEquals(t, "refs/__temp__/3802de4a769092a4b6477e9b5ec0636938c5a957", ref.Ref)
}

func assertLength(t *testing.T, expected int, col interface{}) {
	objValue := reflect.ValueOf(col)
	l := objValue.Len()
	if l != expected {
		t.Fatalf("Expected %d refs but got %d", expected, l)
	}
}

func assertEquals(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Fatalf("Expected [%s] to be equal to [%s]", expected, actual)
	}
}

func assertIEquals(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Fatalf("Expected [%d] to be equal to [%d]", expected, actual)
	}
}