package stopwatchapp

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// Unpack unpacks an firmware bundle, writing files to destination directory.
func Unpack(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0700)

	for _, f := range r.File {
		err := extractFile(f, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

// MustClose closes given io.Closer or panics
func MustClose(f io.Closer) {
	if err := f.Close(); err != nil {
		panic(err)
	}
}

func extractFile(f *zip.File, dest string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer MustClose(rc)

	path := filepath.Join(dest, f.Name)

	if f.FileInfo().IsDir() {
		if err := os.MkdirAll(filepath.Dir(path), f.Mode()); err != nil {
			return err
		}

	} else {
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return err
		}

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer MustClose(f)

		_, err = io.Copy(f, rc)
		if err != nil {
			return err
		}
	}
	return nil
}
