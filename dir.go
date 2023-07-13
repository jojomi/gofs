package gofs

import (
	"fmt"
	"github.com/pkg/browser"
	"github.com/spf13/afero"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

// Dir is a filesystem directory.
type Dir struct {
	path string
	fs   afero.Fs
}

func DirAt(path string) Dir {
	return dirAtWithFs(path, afero.NewOsFs())
}

func dirAtWithFs(path string, fs afero.Fs) Dir {
	// replace home dir path
	path, _ = homedir.Expand(path)

	// make path absolute
	if !filepath.IsAbs(path) {
		pwd, err := os.Getwd()
		if err != nil {
			// this should not happen
			panic(err)
		}
		path = filepath.Join(pwd, path)
	}

	// remove trailing path separator if it exists
	path = strings.TrimRight(path, string(os.PathSeparator))

	return Dir{
		path: path,
		fs:   fs,
	}
}

func (x Dir) Create(perm os.FileMode) error {
	return x.fs.MkdirAll(x.path, perm)
}

func (x Dir) Ensure(perm os.FileMode) error {
	return x.Create(perm)
}

func (x Dir) MustEnsure(perm os.FileMode) Dir {
	err := x.Ensure(perm)
	if err != nil {
		panic(err)
	}
	return x
}

func (x Dir) EnsureEmpty(perm os.FileMode) error {
	if x.Exists() {
		if x.IsEmpty() {
			return nil
		}
		return x.Clear()
	}
	return x.Create(perm)
}

func (x Dir) IsEmpty() bool {
	if !x.Exists() {
		return true
	}

	files, err := afero.ReadDir(x.fs, x.Path())
	if err != nil {
		panic(err)
	}

	return len(files) == 0
}

func (x Dir) AssertEmpty() Dir {
	if !x.IsEmpty() {
		panic(fmt.Errorf("dir %s should have been empty", x))
	}
	return x
}

func (x Dir) IsReadable() bool {
	if !x.Exists() {
		return false
	}
	_, err := afero.ReadDir(x.fs, x.path)
	return err == nil
}

func (x Dir) AssertReadable() Dir {
	if !x.IsReadable() {
		panic(fmt.Errorf("dir %s should have been readable", x))
	}
	return x
}

func (x Dir) IsWritable() bool {
	if !x.Exists() {
		return false
	}
	// Testfile
	t := x.MustFileAt("__writetest_gofs")
	t.AssertNotExists()
	defer t.MustRemove()
	err := t.SetContent([]byte("a"))
	return err == nil
}

func (x Dir) AssertWritable() Dir {
	if !x.IsWritable() {
		panic(fmt.Errorf("dir %s should have been writable", x))
	}
	return x
}

func (x Dir) AssertNotEmpty() Dir {
	if x.IsEmpty() {
		panic(fmt.Errorf("dir %s should not have been empty", x))
	}
	return x
}

func (x Dir) Exists() bool {
	fi, err := x.fs.Stat(x.Path())
	return !os.IsNotExist(err) && fi.IsDir()
}

func (x Dir) ReadDir() ([]os.FileInfo, error) {
	return afero.ReadDir(x.fs, x.Path())
}

func (x Dir) AssertExists() Dir {
	if !x.Exists() {
		panic(fmt.Errorf("dir %s should have existed", x))
	}
	return x
}

func (x Dir) NotExists() bool {
	return !x.Exists()
}

func (x Dir) AssertNotExists() Dir {
	if !x.NotExists() {
		panic(fmt.Errorf("dir %s should not have existed", x))
	}
	return x
}

func (x Dir) Path() string {
	return x.path
}

func (x Dir) Parent() Dir {
	return dirAtWithFs(path.Dir(x.path), x.fs)
}

func (x Dir) RelativeTo(dir Dir) string {
	result, err := filepath.Rel(dir.Path(), x.Path())
	if err != nil {
		return ""
	}
	return result
}

func (x Dir) MustFileAt(relativePath string) File {
	file, err := x.FileAt(relativePath)
	if err != nil {
		panic(err)
	}
	return file
}

func (x Dir) FileAt(relativePath string) (File, error) {
	if filepath.IsAbs(relativePath) {
		return File{}, NewFilePathNotRelativeError(relativePath)
	}

	return fileAtWithFs(filepath.Join(x.path, relativePath), x.fs), nil
}

func (x Dir) MustDirAt(relativePath string) Dir {
	dir, err := x.DirAt(relativePath)
	if err != nil {
		panic(err)
	}
	return dir
}

func (x Dir) DirAt(relativePath string) (Dir, error) {
	if filepath.IsAbs(relativePath) {
		return Dir{fs: x.fs}, NewDirPathNotRelativeError(relativePath)
	}

	return dirAtWithFs(filepath.Join(x.path, relativePath), x.fs), nil
}

func (x Dir) Remove() error {
	if !x.Exists() {
		return nil
	}
	return x.fs.RemoveAll(x.path)
}

func (x Dir) Clear() error {
	dir := x.Path()
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func (x Dir) MustClear() Dir {
	err := x.Clear()
	if err != nil {
		panic(err)
	}
	return x
}

func (x Dir) WithTrailingPathSeparator() string {
	return x.path + string(os.PathSeparator)
}

func (x Dir) WithoutTrailingPathSeparator() string {
	return x.path
}

func (x Dir) Equals(otherDir Dir) bool {
	return x.Path() == otherDir.Path()
}

func (x Dir) String() string {
	return x.path
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func (x Dir) OpenStandard() error {
	return browser.OpenFile(x.Path())
}
