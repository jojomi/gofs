package gofs

import (
	"github.com/spf13/afero"
	"os"
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
	// TODO implement
	return nil
}

func (x Dir) Ensure(perm os.FileMode) error {
	return x.Create(perm)
}

func (x Dir) EnsureEmpty(perm os.FileMode) error {
	if x.Exists() {
		return x.Clear()
	}
	return x.Create(perm)
}

func (x Dir) Exists() bool {
	fi, err := x.fs.Stat(x.Path())
	return !os.IsNotExist(err) && fi.IsDir()
}

func (x Dir) NotExists() bool {
	return !x.Exists()
}

func (x Dir) Path() string {
	return x.path
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
