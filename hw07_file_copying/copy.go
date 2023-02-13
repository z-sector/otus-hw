package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNegativeOffset        = errors.New("offset cannot be negative")
	ErrNegativeLimit         = errors.New("limit cannot be negative")
	ErrFileIsDirectory       = errors.New("file is a directory")
	ErrFromPathIsEmpty       = errors.New("from path should be don't empty")
	ErrToPathIsEmpty         = errors.New("to path should be don't empty")
)

func Copy(fromPath, toPath string, offset, limit int64) (err error) {
	if err := validateCopyParams(fromPath, toPath, offset, limit); err != nil {
		return err
	}

	source, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		if err := source.Close(); err != nil {
			log.Println(err)
		}
	}()

	size, err := validateSource(source, offset)
	if err != nil {
		return err
	}
	if limit == 0 || limit > size {
		limit = size
	}

	if _, err := source.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	dest, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() {
		errCl := dest.Close()
		if errCl != nil && err == nil {
			err = errCl
		}
	}()

	return copyCont(source, dest, limit)
}

func validateCopyParams(fromPath, toPath string, offset, limit int64) error {
	switch {
	case fromPath == "":
		return ErrFromPathIsEmpty
	case toPath == "":
		return ErrToPathIsEmpty
	case offset < 0:
		return ErrNegativeOffset
	case limit < 0:
		return ErrNegativeLimit
	}
	return nil
}

func validateSource(f *os.File, offset int64) (int64, error) {
	sourceInfo, err := f.Stat()
	if err != nil {
		return 0, err
	}

	size := sourceInfo.Size()
	switch {
	case sourceInfo.IsDir():
		return 0, ErrFileIsDirectory
	case size < offset:
		return 0, ErrOffsetExceedsFileSize
	}
	if !sourceInfo.Mode().IsRegular() {
		return 0, ErrUnsupportedFile
	}
	return size - offset, nil
}

func copyCont(source io.Reader, dest io.Writer, limit int64) error {
	const bufferSize = 10 * 1024
	buf := make([]byte, bufferSize)
	var total int64
	var prevProgress int

	for {
		exit := false
		bytesCount, err := source.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				exit = true
			} else {
				return err
			}
		}

		total += int64(bytesCount)
		if total >= limit {
			bytesCount -= int(total - limit)
			exit = true
			total = limit
		}

		if bytesCount > 0 {
			_, err = dest.Write(buf[:bytesCount])
			if err != nil {
				return err
			}
		}

		progress := int(float32(total) / float32(limit) * 100)
		if progress > prevProgress {
			fmt.Printf("%d%%...", progress)
			prevProgress = progress
		}

		if exit {
			break
		}
	}

	fmt.Println()
	return nil
}
