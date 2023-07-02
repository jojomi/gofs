package gofs

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFile(t *testing.T) {
	a := assert.New(t)

	f := FileAt("/tmp/testfile.log")
	f.fs = afero.NewMemMapFs()
	parentDir := DirAt("/tmp")
	testContent := "go test\nanother line"

	// cleanup
	defer cleanup(f)
	cleanup(f)

	// test non existing file
	a.False(f.Exists())
	a.True(f.NotExists())
	a.True(f.IsEmpty())
	a.False(f.IsReadable())
	a.True(f.IsWritable())
	a.True(f.Dir().Equals(parentDir))
	f.MustEnsureDir(0750)
	a.True(f.HasAnyExtension())
	a.True(f.HasExtension(ExtLog))
	f.AssertExtension(ExtLog)
	a.Equal("/tmp/testfile.pdf", f.WithExtension(ExtPdf).Path())

	// test assertions
	f.AssertNotExists()
	f.AssertWritable()
	f.AssertEmpty()

	// create file
	err := f.SetContentString(testContent)
	a.Nil(err)
	a.True(f.Exists())
	a.False(f.NotExists())
	a.False(f.IsEmpty())

	f.AssertExists()
	f.AssertReadable()
	f.AssertNotEmpty()

	// read content
	byteContent := []byte(testContent)
	b, err := f.Content()
	a.Nil(err)
	a.Equal(byteContent, b)

	// read content (string)
	c, err := f.StringContent()
	a.Nil(err)
	a.Equal(testContent, c)

	extraContent := "third line"
	err = f.AppendStringln(extraContent)
	a.Nil(err)

	a.Equal([]byte(testContent+extraContent+"\n"), f.MustContent())
	a.Equal(testContent+extraContent+"\n", f.MustStringContent())

	// test hashing
	expectedMD5Hash := "9656e885da2a0510f480c5ff8a6a57b5"
	a.Equal(expectedMD5Hash, f.MustMd5Hash())
	f.AssertMd5Hash(expectedMD5Hash)
}

func cleanup(f File) {
	f.MustRemove()
}
