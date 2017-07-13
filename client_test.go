package gogitsmartproto_test

import (
	"testing"
	"github.com/gaffo/gogitsmartproto"
)

func Test_RefsRails(t *testing.T) {
	client, err := gogitsmartproto.NewClient("https://github.com/rails/rails")
	failIfError(t, err)
	resp, err := client.Refs()
	failIfError(t, err)
	if len(resp.Refs) < 10 {
		t.Fatal("Didn't have enough refs")
	}
}