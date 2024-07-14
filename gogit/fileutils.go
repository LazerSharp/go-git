package gogit

import (
	"io"
	"io/fs"
	"os"
)

func isDir(path string) bool {
	exts, f := exists(path)
	if !exts {
		return false
	}
	return f.Mode().IsDir()
}

func isDirEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}

func exists(path string) (bool, fs.FileInfo) {
	f, err := os.Stat(path)
	if err != nil {
		return false, nil
	}
	return true, f
}

func Must[T any](obj T, err error) T {
	Check(err)
	return obj
}

func Check(err error) {
	if err != nil {
		panic(err)
		// fmt.Println("fatal: ", err)
		// os.Exit(1)
	}
}
