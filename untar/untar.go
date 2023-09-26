package untar

import (
	"archive/tar"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
)

type TarBoll struct {
	Source string
	Target string
}

func (t *TarBoll) loopAndSave(tarReader *tar.Reader) error {

	// create target dir
	if _, err := os.Stat(t.Target); err != nil {
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

func (t *TarBoll) UnTar() error {

	switch s := path.Ext(t.Source); s {

	case ".tar":
		return t.tar()
	case ".gz":
		return t.tgz()
	case ".tgz":
		return t.tgz()
	case ".bz2":
		return t.bz2()
	default:
		err := errors.New("invalid file type")
		return err
	}
}

func (t *TarBoll) tar() error {
	var tarReader *tar.Reader
	fw, err := os.Open(t.Source)
	if err != nil {
		return err
	}
	defer fw.Close()
	tarReader = tar.NewReader(fw)

	return t.loopAndSave(tarReader)
}

func (t *TarBoll) bz2() error {

	fw, err := os.Open(t.Source)
	if err != nil {
		return err
	}
	defer fw.Close()

	br := bzip2.NewReader(fw)
	tarReader := tar.NewReader(br)
	return t.loopAndSave(tarReader)
}

func (t *TarBoll) tgz() error {
	var tarReader *tar.Reader
	fw, err := os.Open(t.Source)
	if err != nil {
		return err
	}
	defer fw.Close()

	// read gzip
	gr, err := gzip.NewReader(fw)
	if err != nil {
		return err
	}
	defer gr.Close()

	tarReader = tar.NewReader(gr)

	return t.loopAndSave(tarReader)
}

func (t *TarBoll) saveFile(fw *os.File, tr *tar.Reader) error {
	defer fw.Close()
	_, err := io.Copy(fw, tr)
	if err != nil {
		return err
	}
	return nil
}
