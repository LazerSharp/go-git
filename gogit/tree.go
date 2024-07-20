package gogit

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
)

type Entry struct {
	Mode string
	Type string
	Hash string
	Path string
}

type Tree struct {
	Entries []Entry
}

func NewTree(r io.Reader) (*Tree, error) {
	tree := &Tree{
		Entries: make([]Entry, 0),
	}
	err := tree.DeSerialize(r)
	if err != nil {
		return nil, err
	}
	return tree, nil
}

func (t *Tree) Type() string {
	return "tree"
}

func (e *Entry) Serialize(w io.Writer) error {
	w.Write([]byte(e.Type))
	w.Write([]byte(e.Mode))
	w.Write([]byte(" "))

	return nil
}

func (t *Tree) Serialize(w io.Writer) error {
	for _, e := range t.Entries {
		err := e.Serialize(w)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseMode(data string) (typ string, mode string, err error) {
	t := data[:2]
	switch t {
	case "10":
		typ = "blob"
	case "04":
		typ = "tree"
	case "12":
		typ = "blob" // A symlink. Blob contents is link target.
	case "16":
		typ = "commit" // A submodule
	}

	return typ, data[2:], nil
}

func readEntry(r *bufio.Reader) (*Entry, error) {

	// read mode
	mode, err := r.ReadString(' ')
	if err != nil {
		return nil, err
	}
	mode = mode[:len(mode)-1]
	if len(mode) == 5 {
		mode = fmt.Sprintf("0%s", mode)
	}

	// read path
	path, err := r.ReadString('\x00')
	if err != nil {
		return nil, err
	}
	path = path[:len(path)-1]
	fmt.Println("path", path)

	// read hash
	shaBytes := make([]byte, 20)
	_, err = r.Read(shaBytes)
	if err != nil {
		return nil, err
	}
	sha := hex.EncodeToString(shaBytes)
	fmt.Println("sha", sha)

	entry := &Entry{
		Path: path,
		Hash: sha,
	}
	fmt.Printf("mode= [%s] len = %d\n", mode, len(mode))
	entry.Type, entry.Mode, err = parseMode(mode)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (t *Tree) DeSerialize(r io.Reader) error {

	reader := bufio.NewReader(r)

	for {
		entry, err := readEntry(reader)
		if err != nil && err != io.EOF {
			return err
		}
		if entry == nil {
			break
		}
		t.Entries = append(t.Entries, *entry)
	}

	return nil

}
