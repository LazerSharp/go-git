package gogit

import "io"

type Blob struct {
	Content []byte
}

func NewBlob(r io.Reader) (*Blob, error) {
	b := &Blob{}
	if r != nil {
		if err := b.DeSerialize(r); err != nil {
			return nil, err
		}
	}
	return b, nil
}

func NewEmptyBlob() (*Blob, error) {
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
