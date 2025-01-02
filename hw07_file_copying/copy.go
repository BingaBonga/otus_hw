package main

import (
	"errors"
	"io"
	"log"
	"os"
)

var (
	ErrFailedToRead          = errors.New("failed to read")
	ErrFailedToWrite         = errors.New("failed to write")
	ErrFileDoesNotExist      = errors.New("file does not exist")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) (err error) {
	defer func() {
		if r := recover(); r != nil && err != nil {
			log.Println("Recovered from panic:", r)
			err = r.(error)
		}
	}()

	fromFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0666)
	check(err)
	defer func() {
		check(fromFile.Close())
	}()

	toFile, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if os.IsNotExist(err) {
		toFile, err = os.Create(toPath)
	}

	check(err)
	defer func() {
		check(toFile.Close())
	}()

	fromFileInfo, err := os.Stat(fromPath)
	check(err)

	fileSize := fromFileInfo.Size()
	if offset > fileSize {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 {
		limit = fileSize
	}

	read := int64(0)
	buf := make([]byte, 1024)
	for read < limit {
		readAt, errRead := fromFile.ReadAt(buf, offset+read)
		if errRead != nil && errRead != io.EOF {
			return ErrFailedToRead
		}

		readAt64 := int64(readAt)
		if read+readAt64 > limit {
			readAt64 = limit - read
		}

		_, errWrite := toFile.WriteAt(buf[:readAt64], read)
		if errWrite != nil {
			return ErrFailedToWrite
		}

		read += readAt64
		if errRead == io.EOF {
			break
		}
	}

	return nil
}

func check(err error) {
	if err != nil {
		if os.IsNotExist(err) {
			panic(ErrFileDoesNotExist)
		}

		panic(ErrUnsupportedFile)
	}
}
