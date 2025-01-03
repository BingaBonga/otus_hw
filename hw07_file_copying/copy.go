package main

import (
	"bufio"
	"errors"
	"io"
	"os"

	//nolint:depguard
	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) (err error) {
	fromFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0o666)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	toFile, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if os.IsNotExist(err) {
		toFile, err = os.Create(toPath)
	}
	if err != nil {
		return err
	}
	defer toFile.Close()

	fromFileInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	if fromFileInfo.IsDir() {
		return ErrUnsupportedFile
	}

	fileSize := fromFileInfo.Size()
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 || limit+offset > fileSize {
		limit = fileSize - offset
	}

	_, err = fromFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	bufReader := bufio.NewReaderSize(fromFile, 1024*1024)
	progressBar := pb.New64(limit).Start()
	barReader := progressBar.NewProxyReader(bufReader)

	_, err = io.CopyN(toFile, barReader, limit)
	progressBar.Finish()
	return err
}
