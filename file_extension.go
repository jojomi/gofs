package gofs

import "strings"

type FileExtension struct {
	extension string
}

func FileExtensionFrom(extension string) FileExtension {
	extension = strings.TrimPrefix(extension, ".")

	return FileExtension{
		extension: extension,
	}
}

func (x FileExtension) WithDot() string {
	return "." + x.extension
}

func (x FileExtension) WithoutDot() string {
	return x.extension
}

func (x FileExtension) String() string {
	return x.WithDot()
}

var (
	ExtPdf = FileExtensionFrom("pdf")
	ExtJpg = FileExtensionFrom("jpg")
	ExtPng = FileExtensionFrom("png")
	ExtLog = FileExtensionFrom("log")
	ExtZip = FileExtensionFrom("zip")
)
