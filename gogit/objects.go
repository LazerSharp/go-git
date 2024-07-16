package gogit

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
)

type GitObject interface {
	Len() int
	Type() string
	Serialize(w io.Writer) error
	DeSerialize(r io.Reader) error
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
	byts := make([]byte, ln)
	n, err := br.Read(byts)

	if err != nil {
		return nil, err
	}
	if n != int(ln) {
		return nil, fmt.Errorf("file corrupted! wrong length: %v", ln)
	}
	reader := bytes.NewReader(byts)
	switch typ {
	case "blob":
		return NewBlob(reader), nil
	case "commit":
		return NewBlob(reader), nil
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

func Sha1(content []byte) string {
	h := sha1.New()
	h.Write(content)
	return hex.EncodeToString(h.Sum(nil))
}

func WriteObject(obj GitObject, repo *GitRepository) (sha string, err error) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s %d\x00", obj.Type(), obj.Len())
	obj.Serialize(&buf)

	sha = Sha1(buf.Bytes())
	if repo == nil {
		return sha, nil
	}

	// write to file
	fpth, err := RepoFile(repo, true, sha[:2], sha[2:])

	if e, _ := exists(fpth); e {
		return sha, nil
	}

	if err != nil {
		return "", err
	}
	f, err := os.Create(fpth)
	if err != nil {
		return "", err
	}
	w := zlib.NewWriter(f)
	defer w.Close()
	w.Write(buf.Bytes())

	return sha, nil
}
