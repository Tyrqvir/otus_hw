package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSeekFile              = errors.New("can't seek")
	size                     int64
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	source, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	if offset != 0 {
		_, err := source.Seek(offset, io.SeekStart)
		if err != nil {
			return ErrSeekFile
		}
	}

	defer source.Close()

	sourceInfo, err := source.Stat()
	if err != nil {
		return err
	}

	size = sourceInfo.Size()

	if offset > size {
		return ErrOffsetExceedsFileSize
	}
	if size == 0 {
		return ErrUnsupportedFile
	}

	destination, err := os.Create(toPath)
	if err != nil {
		return err
	}

	defer destination.Close()

	limit = func() int64 {
		if limit == 0 || limit > size {
			return size - offset
		}
		return limit
	}()

	bar := pb.Default.Start64(limit)
	barReader := bar.NewProxyReader(source)

	_, err = io.CopyN(destination, barReader, limit)
	if err != nil {
		return err
	}

	bar.Finish()

	return nil
}
