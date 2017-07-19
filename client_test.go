package gogitsmartproto_test

import (
	"testing"
	"github.com/gaffo/gogitsmartproto"
)

/*
func Test_RefsRails(t *testing.T) {
	client, err := gogitsmartproto.NewClient("https://github.com/rails/rails")
	failIfError(t, err)
	resp, err := client.Refs()
	failIfError(t, err)
	if len(resp.Refs) < 10 {
		t.Fatal("Didn't have enough refs")
	}
}
*/

func Test_Scan(t *testing.T) {
	client, err := gogitsmartproto.NewClient("https://github.com/gaffo/parses_travel_emails")
	failIfError(t, err)
	//resp, err := client.Refs()
	//failIfError(t, err)
	//ref := resp.Refs[0]
	res, err := client.Packs("3802de4a769092a4b6477e9b5ec0636938c5a957")
	failIfError(t, err)
	t.Fatal(res)
}
