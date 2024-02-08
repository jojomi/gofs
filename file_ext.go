package gofs

import (
	"fmt"
	"regexp"
	"strings"
)

func (x File) HasAnyExtension() bool {
	filename := x.Filename()
	return strings.Contains(filename[1:], ".")
}

func (x File) HasExtension(fileExtension FileExtension) bool {
	return strings.HasSuffix(x.path, fileExtension.WithDot())
}

func (x File) AssertExtension(fileExtension FileExtension) File {
	if !x.HasExtension(fileExtension) {
		panic(fmt.Errorf("file %s should have had file extension %s", x, fileExtension))
	}
	return x
}

func (x File) WithExtension(fileExtension FileExtension) File {
	if !x.HasAnyExtension() {
		return fileWithSameFs(x.path+fileExtension.WithDot(), x)
	}

	r := regexp.MustCompile(`(\.tar)?\.[^.]+$`)
	newFilename := r.ReplaceAllString(x.Filename(), fileExtension.WithDot())

	return x.Dir().MustFileAt(newFilename)
}

func (x File) Extension() *FileExtension {
	if !x.HasAnyExtension() {
		return nil
	}

	r := regexp.MustCompile(`(\.tar)?\.[^.]+$`)
	extension := r.FindString(x.Filename())

	result := FileExtensionFrom(extension)
	return &result
}
