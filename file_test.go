package gofs

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"text/template"
)

func TestFile(t *testing.T) {
	a := assert.New(t)
	fs := afero.NewMemMapFs()

	f := FileAt("/tmp/testfile.log")
	f.fs = fs
	parentDir := DirAt("/tmp")
	testContent := "go test\nanother line"

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
	c, err := f.ContentString()
	a.Nil(err)
	a.Equal(testContent, c)

	extraContent := "third line"
	err = f.AppendStringln(extraContent)
	a.Nil(err)

	a.Equal([]byte(testContent+extraContent+"\n"), f.MustContent())
	a.Equal(testContent+extraContent+"\n", f.MustContentString())

	// test io
	fCopy := FileAt("/tmp/copy")
	fCopy.fs = fs
	err = f.CopyTo(fCopy)
	a.Nil(err)
	a.True(fCopy.Exists())
	a.Equal(f.MustContentString(), fCopy.MustContentString())

	// test hashing
	expectedMD5Hash := "9656e885da2a0510f480c5ff8a6a57b5"
	a.Equal(expectedMD5Hash, f.MustMd5Hash())
	f.AssertMd5Hash(expectedMD5Hash)

	// test template
	outFile := FileAt("/tmp/template.out")
	outFile.fs = fs
	err = f.SetContentString("")
	a.Nil(err)
	err = f.AssertEmpty().AppendString("{{ .Name }} was here")
	a.Nil(err)
	err = f.MustRenderer().WithData(map[string]string{
		"Name": "John",
	}).RenderToFile(outFile)
	a.Nil(err)
	c, err = outFile.ContentString()
	a.True(outFile.Exists())
	a.False(outFile.IsEmpty())
	a.Nil(err)
	a.Equal("John was here", c)

	// template with custom funcs
	err = f.MustClear().AppendString("{{ .Name | toLower | toUpper }} was here")
	a.Nil(err)
	var buf strings.Builder
	err = f.MustRenderer().WithData(map[string]string{
		"Name": "Tom",
	}).AddFuncs(template.FuncMap{
		"toUpper": strings.ToUpper,
	}).AddFuncs(template.FuncMap{
		"toLower": strings.ToLower,
	}).RenderTo(&buf)
	a.Nil(err)
	a.Equal("TOM was here", buf.String())
}

func TestFileAtDir(t *testing.T) {
	a := assert.New(t)
	fs := afero.NewMemMapFs()

	f := FileAtDir(DirAt("/tmp/"), "testfile.log")
	f.fs = fs

	a.Equal("/tmp/testfile.log", f.Path())
}

func TestWithFilename(t *testing.T) {
	a := assert.New(t)
	fs := afero.NewMemMapFs()

	f := FileAt("/tmp/testfile.log")
	f.fs = fs

	a.Equal("/tmp/out.pdf", f.WithFilename("out.pdf").Path())
}

func TestInDir(t *testing.T) {
	a := assert.New(t)
	fs := afero.NewMemMapFs()

	f := FileAt("/tmp/testfile.log")
	f.fs = fs

	a.Equal("/run/testfile.log", f.InDir(DirAt("/run/")).Path())
}

func TestWithoutExtension(t *testing.T) {
	a := assert.New(t)
	fs := afero.NewMemMapFs()

	f := FileAt("/tmp/testfile.exe")
	f.fs = fs

	a.Equal("testfile", f.WithoutExtension().Filename())
}
