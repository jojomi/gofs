package gofs

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/spf13/afero"
	"io"
	"os"
)

func (x File) IsReadable() bool {
	f, err := x.fs.Open(x.Path())
	if err == nil {
		_ = f.Close()
	}
	return err == nil
}

func (x File) AssertReadable() File {
	if !x.IsReadable() {
		panic(fmt.Errorf("file %s should have been readable", x))
	}
	return x
}

func (x File) IsWritable() bool {
	filePath := x.Path()
	existed := x.Exists()
	f, err := os.OpenFile(filePath, os.O_RDWR+os.O_CREATE, x.createPermissions)
	if err == nil {
		_ = f.Close()
		if !existed {
			err = os.Remove(filePath)
			if err != nil {
				panic(err)
			}
		}
	}
	return err == nil
}

func (x File) AssertWritable() File {
	if !x.IsWritable() {
		panic(fmt.Errorf("file %s should have been writable", x))
	}
	return x
}

func (x File) Content() ([]byte, error) {
	return afero.ReadFile(x.fs, x.Path())
}

func (x File) MustContent() []byte {
	content, err := x.Content()
	if err != nil {
		panic(errors.Annotatef(err, "could not read content of %s", x))
	}
	return content
}

func (x File) StringContent() (string, error) {
	content, err := afero.ReadFile(x.fs, x.Path())
	return string(content), err
}

func (x File) MustStringContent() string {
	content, err := x.StringContent()
	if err != nil {
		panic(errors.Annotatef(err, "could not read content of %s", x))
	}
	return string(content)
}

func (x File) Append(newContent []byte) error {
	f, err := x.fs.OpenFile(x.Path(), os.O_RDWR, x.createPermissions)
	if err != nil {
		return err
	}

	// seek to end of file
	_, err = f.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	_, err = f.Write(newContent)
	return err
}

func (x File) AppendString(newContent string) error {
	return x.Append([]byte(newContent))
}

func (x File) AppendStringln(newContent string) error {
	return x.AppendString(newContent + "\n")
}

func (x File) SetContent(newContent []byte) error {
	return afero.WriteFile(x.fs, x.Path(), newContent, x.createPermissions)
}

func (x File) SetContentString(newContent string) error {
	return x.SetContent([]byte(newContent))
}

func (x File) Remove() error {
	if x.NotExists() {
		return nil
	}
	return x.fs.Remove(x.Path())
}

func (x File) MustRemove() File {
	err := x.Remove()
	if err != nil {
		panic(err)
	}
	return x
}
