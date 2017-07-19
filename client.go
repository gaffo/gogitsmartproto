package gogitsmartproto

import (
	"net/http"
	"fmt"
	"strings"
	"bytes"
	"compress/gzip"
	"os"
	"io"
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

func NewClient(base string) (*Client, error) {
	return &Client{
		base: base,
	}, nil
}

func rts(resp *http.Response) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	return buf.String()
}

func rtsgz(resp *http.Response) string {
	gzr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return ""
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(gzr)
	return buf.String()
}

func (this *Client) uploadPackUrl() string {
	return fmt.Sprintf("%s.git/info/refs?service=git-upload-pack", this.base)
}

func (this *Client) walkPackUrl() string {
	return fmt.Sprintf("%s/git-upload-pack", this.base)
}

func (this *Client) Refs() (Refs, error) {
	resp, err := http.Get(this.uploadPackUrl())
	if err != nil {
		return Refs{}, err
	}

	return ParseRefs(resp.Body)
}

func (this *Client) Packs(want string) (string, error) {
	req := "0032want 20a6d91a4ab05a95635fba68872ee5b38c68e16e\n00000009done\n"
	url := this.walkPackUrl()
	hReq, err := http.NewRequest("POST", url, strings.NewReader(req))
	if err != nil {
		return "", err
	}
	hReq.Header.Add("User-Agent", "JGit/4.8.0.201706111038-r")

	client := &http.Client{}
	resp, err := client.Do(hReq)

	f, err := os.Create("tmp.dat")
	if err != nil {
		return "", err
	}

	io.Copy(f, resp.Body)

	return "", nil

	//return rts(resp), nil
}
