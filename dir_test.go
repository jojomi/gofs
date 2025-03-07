package gofs

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDir(t *testing.T) {
	a := assert.New(t)

	d := DirAt("/tmp/testdir")
	d.fs = afero.NewMemMapFs()
	parentDir := DirAt("/tmp")

	a.Equal("testdir", d.Name())

	// test non existing dir
	a.False(d.Exists())
	a.True(d.NotExists())
	a.True(d.IsEmpty())
	a.True(d.Parent().Equals(parentDir))
	d.MustEnsure(0750)

	// test assertions
	d.AssertExists()
	d.AssertReadable()
	d.AssertWritable()
	d.AssertEmpty()

	// create file
	f := d.MustFileAt("testfile")
	err := f.SetContentString("content")
	a.Nil(err)

	d.AssertNotEmpty()
}
