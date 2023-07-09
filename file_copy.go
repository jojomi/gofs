package gofs

import (
	"github.com/juju/errors"
	"io"
)

func (x File) CopyTo(target File) error {
	// Open the source file
	src, err := x.fs.Open(x.path)
	if err != nil {
		return errors.Annotatef(err, "Error opening source file")
	}
	defer src.Close()

	// Create the destination file
	dest, err := x.fs.Create(target.Path())
	if err != nil {
		return errors.Annotatef(err, "Error creating destination file")
	}
	defer dest.Close()

	// Copy the contents from source to destination
	_, err = io.Copy(dest, src)
	if err != nil {
		return errors.Annotatef(err, "Error copying file")
	}
	return nil
}
