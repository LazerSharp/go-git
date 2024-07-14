package gogit

import (
	"bufio"
	"compress/zlib"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
)

type GitObject interface {
	Type() string
	Serialize(w io.Writer)
}

type Blob struct {
	Content []byte
}

func (b *Blob) Serialize(w io.Writer) {
	w.Write(b.Content)
}

func (b *Blob) Type() string {
	return "blob"
}

func DeSerializeObject(r io.Reader) (GitObject, error) {
	br := bufio.NewReader(r)
	// read object type
	typ, err := br.ReadString(' ')
	if err != nil {
		return nil, err
	}
	typ = typ[:len(typ)-1]

	// read length
	lnstr, err := br.ReadString('\x00')
	if err != nil {
		return nil, err
	}
	ln, err := strconv.ParseInt(lnstr[:len(lnstr)-1], 10, 64)
	if err != nil {
		return nil, err
	}

	slog.Debug("DeSerializeObject: ", "type", typ, "length", ln)
	bytes := make([]byte, ln)
	n, err := br.Read(bytes)

	if err != nil {
		return nil, err
	}
	if n != int(ln) {
		return nil, fmt.Errorf("file corrupted! wrong length: %v", ln)
	}

	switch typ {
	case "blob":
		return &Blob{bytes}, nil
	}

	return nil, err

}

func ReadObject(repo *GitRepository, sha string) (GitObject, error) {

	fpath, err := RepoFile(repo, false, "objects", sha[:2], sha[2:])

	if err != nil {
		return nil, err
	}

	// if e, _ := exists(fpath); !e {
	// 	return nil, fmt.Errorf("Object file '%s' does not exists", fpath)
	// }

	f, err := os.Open(fpath)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	r, err := zlib.NewReader(f)

	if err != nil {
		return nil, err
	}
	defer r.Close()

	return DeSerializeObject(r)

}
