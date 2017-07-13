package gogitsmartproto

import (
	"net/http"
	"fmt"
)

type Ref struct {
	Checksum string
	Ref      string
}

type Refs struct {
	Refs []Ref
}

type Client struct {
	base string
}

func NewClient(base string) (*Client, error){
	return &Client{
		base: base,
	}, nil
}

func (this *Client) Refs() (Refs, error) {
	resp, err := http.Get(fmt.Sprintf("%s.git/info/refs?service=git-upload-pack", this.base))
	if err != nil {
		return Refs{}, err
	}

	return ParseRefs(resp.Body)
}