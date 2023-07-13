package gofs

import (
	"fmt"
	"github.com/spf13/afero"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/browser"
)

// File is a file in the filesystem.
type File struct {
	// path is the absolute filepath
	path              string
	createPermissions os.FileMode
	fs                afero.Fs
}

func FileAt(filePath string) File {
	return fileAtWithFs(filePath, afero.NewOsFs())
}

func fileWithSameFs(filePath string, f File) File {
	return fileAtWithFs(filePath, f.fs)
}

func FileWithFs(filePath string, fs afero.Fs) File {
	return fileAtWithFs(filePath, fs)
}

func fileAtWithFs(filePath string, fs afero.Fs) File {
	// replace home dir path
	filePath, err := homedir.Expand(filePath)
	if err != nil {
		// this should really not happen
		panic(err)
	}

	// make path absolute
	if !filepath.IsAbs(filePath) {
		pwd, err := os.Getwd()
		if err != nil {
			// this should not happen
			panic(err)
		}
		filePath = filepath.Join(pwd, filePath)
	}

	return File{
		path:              filePath,
		createPermissions: 0640,
		fs:                fs,
	}
}

// SetCreatePermissions allows you to define the FileMode used when creating this file (if it did not exist).
func (x File) SetCreatePermissions(perm os.FileMode) File {
	x.createPermissions = perm
	return x
}

func (x File) Exists() bool {
	fi, err := x.fs.Stat(x.Path())
	return !os.IsNotExist(err) && !fi.IsDir()
}

func (x File) AssertExists() File {
	if !x.Exists() {
		panic(fmt.Errorf("file %s should have existed", x))
	}
	return x
}

func (x File) NotExists() bool {
	return !x.Exists()
}

func (x File) AssertNotExists() File {
	if !x.NotExists() {
		panic(fmt.Errorf("file %s should not have existed", x))
	}
	return x
}

func (x File) Path() string {
	return x.path
}

func (x File) RelativeTo(to Dir) string {
	result, err := filepath.Rel(to.Path(), x.Path())
	if err != nil {
		return x.Path()
	}
	return result
}

func (x File) IsHidden() bool {
	return strings.HasPrefix(x.Filename(), ".")
}

func (x File) Equals(otherFile File) bool {
	return x.Path() == otherFile.Path()
}

func (x File) String() string {
	return x.path
}

func (x File) Filename() string {
	return path.Base(x.path)
}

func (x File) WithoutExtension() File {
	return fileAtWithFs(strings.TrimRight(x.path, path.Ext(x.path)), x.fs)
}

func (x File) OpenStandard() error {
	return browser.OpenFile(x.Path())
}

func (x File) Filesize() int64 {
	if !x.Exists() {
		return 0
	}
	if !x.IsReadable() {
		panic(fmt.Errorf("file %s was not readable to check filesize", x))
	}
	fi, err := x.fs.Stat(x.Path())
	if err != nil {
		// this should not happen
		panic(err)
	}
	return fi.Size()
}

func (x File) FilesizeHuman() string {
	const unit = 1024
	b := x.Filesize()
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func (x File) IsEmpty() bool {
	return x.Filesize() == 0
}

func (x File) AssertEmpty() File {
	if !x.IsEmpty() {
		panic(fmt.Errorf("file %s should have been empty", x))
	}
	return x
}

func (x File) AssertNotEmpty() File {
	if x.IsEmpty() {
		panic(fmt.Errorf("file %s should not have been empty", x))
	}
	return x
}

func (x File) WithFileReadOnly(logic func(f afero.File) error) error {
	f, err := x.fs.Open(x.Path())
	if err != nil {
		return err
	}

	err = logic(f)

	closeErr := f.Close()

	if err != nil {
		return err
	}
	return closeErr
}
