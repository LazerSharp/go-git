package gogit

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
)

type EntryType string

const (
	BlobEntry   EntryType = "10"
	TreeEntry   EntryType = "4"
	CommitEntry EntryType = "16"
)

func (b EntryType) Val() string {
	return string(b)
}

func (b EntryType) String() string {
	return map[EntryType]string{
		BlobEntry:   "blob",
		TreeEntry:   "tree",
		CommitEntry: "commit"}[b]
}

type Entry struct {
	Mode string
	Type EntryType
	Hash string
	Path string
}

func writeBytes(w *bufio.Writer, byts ...[]byte) (err error) {
	for _, bs := range byts {
		_, err = w.Write(bs)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Entry) Serialize(w io.Writer) error {
	bufw := bufio.NewWriter(w)

	hash, err := hex.DecodeString(e.Hash)
	if err != nil {
		return err
	}

	err = writeBytes(
		bufw,
		[]byte(e.Type.Val()),
		[]byte(e.Mode),
		[]byte(" "),
		[]byte(e.Path),
		[]byte{'\x00'},
		hash,
	)
	if err != nil {
		return err
	}

	if err = bufw.Flush(); err != nil {
		return err
	}
	return nil
}

func readEntry(r *bufio.Reader) (*Entry, error) {

	// read mode
	mode, err := r.ReadString(' ')
	if err != nil {
		return nil, err
	}
	mode = mode[:len(mode)-1]
	// if len(mode) == 5 {
	// 	mode = fmt.Sprintf("0%s", mode)
	// }

	// read path
	path, err := r.ReadString('\x00')
	if err != nil {
		return nil, err
	}
	path = path[:len(path)-1]
	//fmt.Println("path", path)

	// read hash
	shaBytes := make([]byte, 20)
	_, err = r.Read(shaBytes)
	if err != nil {
		return nil, err
	}
	sha := hex.EncodeToString(shaBytes)
	//fmt.Println("sha", sha)

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

func parseMode(data string) (typ EntryType, mode string, err error) {
	l := len(data)
	return EntryType(data[:l-4]), data[l-4:], nil
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

func (t *Tree) Serialize(w io.Writer) error {
	for _, e := range t.Entries {
		err := e.Serialize(w)
		if err != nil {
			return err
		}
	}
	return nil
}
