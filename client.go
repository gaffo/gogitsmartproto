package gogitsmartproto

type Ref struct {
	Checksum string
	Ref      string
}

type Refs struct {
	Refs []Ref
}

type Client struct {
}

func (this *Client) Refs(repo string) (Refs, error) {
	return Refs{}, nil
}
