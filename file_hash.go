package gofs

import (
	"crypto/md5"
	"fmt"
	"github.com/spf13/afero"
	"hash"
	"io"
)

func (x File) Md5Hash() (string, error) {
	var md5hash hash.Hash
	err := x.WithFileReadOnly(func(f afero.File) error {
		md5hash = md5.New()
		_, err := io.Copy(md5hash, f)
		return err
	})

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", md5hash.Sum(nil)), nil
}

func (x File) MustMd5Hash() string {
	md5Hash, err := x.Md5Hash()
	if err != nil {
		panic(err)
	}
	return md5Hash
}

func (x File) AssertMd5Hash(hash string) File {
	if x.MustMd5Hash() != hash {
		panic(fmt.Errorf("file %s should have had md5 hash %s", x, hash))
	}
	return x
}
