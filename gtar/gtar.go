package gtar

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type UnTar struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

func (t *UnTar) saveFile(fw *os.File, tr *tar.Reader) error {
	defer fw.Close()
	_, err := io.Copy(fw, tr)
	if err != nil {
		return err
	}
	return nil
}

func (t *UnTar) detFileType(f *os.File) string {
	buff := make([]byte, 512)

	_, err := f.Read(buff)
	if err != nil {
		panic(err)
	}
	filetype := http.DetectContentType(buff)

	f.Seek(0, io.SeekStart)
	return filetype
}

func (t *UnTar) UnTarGzip() error {
	var tarReader *tar.Reader
	fw, err := os.Open(t.Source)
	if err != nil {
		return err
	}
	defer fw.Close()

	if fileType := t.detFileType(fw); fileType == "application/x-gzip" {
		// read gzip
		gr, err := gzip.NewReader(fw)
		if err != nil {
			return err
		}
		defer gr.Close()

		tarReader = tar.NewReader(gr)
	} else {
		tarReader = tar.NewReader(fw)
	}

	// create target dir

	if _, err = os.Stat(t.Target); err != nil {
		if err = os.MkdirAll(t.Target, 0755); err != nil {
			return err
		}
	}

	// loop tar reader
	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		targetPath := filepath.Join(t.Target, header.Name)

		fileInfo := header.FileInfo()
		if fileInfo.IsDir() {
			if err = os.MkdirAll(targetPath, fileInfo.Mode()); err != nil {
				return err
			}
			continue
		}
		targetFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
		if err != nil {
			return err
		}
		err = t.saveFile(targetFile, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}
