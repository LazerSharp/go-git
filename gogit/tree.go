package gogit

import "io"

type Tree struct {
}

func (c *Tree) Type() string {
	return "tree"
}

func (c *Tree) Serialize(w io.Writer) error {
	return nil
}

func (c *Tree) Deserialize(w io.Writer) error {
	return nil
}
