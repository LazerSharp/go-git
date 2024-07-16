package gogit

import "io"

type Blob struct {
	Content []byte
}

func NewBlob(r io.Reader) *Blob {
	b := &Blob{}
	if r != nil {
		b.DeSerialize(r)
	}
	return b
}

func NewEmptyBlob() *Blob {
	return NewBlob(nil)
}

func (b *Blob) DeSerialize(r io.Reader) (err error) {
	b.Content, err = io.ReadAll(r)
	return err
}

func (b *Blob) Serialize(w io.Writer) error {
	_, err := w.Write(b.Content)
	return err
}

func (b *Blob) Type() string {
	return "blob"
}

func (b *Blob) Len() int {
	return len(b.Content)
}
