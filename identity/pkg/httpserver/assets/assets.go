package assets

import (
	"io/fs"
	"path/filepath"
)

type HttpStaticFS struct {
	fs fs.FS
}

func NewHttpStaticFS(fs fs.FS) *HttpStaticFS {
	return &HttpStaticFS{
		fs: fs,
	}
}

// Implements the interface fs.FS
func (w HttpStaticFS) Open(name string) (fs.File, error) {
	file, err := w.fs.Open(name)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		// We should never list the content of a folder, index.html is mandatory
		idxPath := filepath.Join(name, "index.html")

		_, err := w.fs.Open(idxPath)
		if err != nil {
			closeErr := file.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return file, nil
}
