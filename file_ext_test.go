package gofs

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFile_AssertExtension(t *testing.T) {
	type fields struct {
		path              string
		createPermissions os.FileMode
		fs                afero.Fs
	}
	type args struct {
		fileExtension FileExtension
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		panics bool
	}{
		{
			name: "valid extension",
			fields: fields{
				path:              "image.jpg",
				createPermissions: 0,
				fs:                nil,
			},
			args: args{
				fileExtension: ExtJpg,
			},
			panics: false,
		},
		{
			name: "invalid extension",
			fields: fields{
				path:              "video.mp4",
				createPermissions: 0,
				fs:                nil,
			},
			args: args{
				fileExtension: ExtPng,
			},
			panics: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := File{
				path:              tt.fields.path,
				createPermissions: tt.fields.createPermissions,
				fs:                tt.fields.fs,
			}
			if tt.panics {
				assert.Panicsf(t, func() {
					x.AssertExtension(tt.args.fileExtension)
				}, "AssertExtension panics")
			} else {
				assert.NotPanicsf(t, func() {
					x.AssertExtension(tt.args.fileExtension)
				}, "AssertExtension does not panic")
			}
		})
	}
}

func TestFile_Extension(t *testing.T) {
	type fields struct {
		path              string
		createPermissions os.FileMode
		fs                afero.Fs
	}

	extTarGz := FileExtensionFrom(".tar.gz")
	extMp3 := FileExtensionFrom("mp3")

	tests := []struct {
		name   string
		fields fields
		want   *FileExtension
	}{
		{
			name: "file with standard extension",
			fields: fields{
				path:              "my.log",
				createPermissions: 0,
				fs:                nil,
			},
			want: &ExtLog,
		},
		{
			name: "file with number in extension",
			fields: fields{
				path:              "music.mp3",
				createPermissions: 0,
				fs:                nil,
			},
			want: &extMp3,
		},
		{
			name: "file with double extension",
			fields: fields{
				path:              "my.tar.gz",
				createPermissions: 0,
				fs:                nil,
			},
			want: &extTarGz,
		},
		{
			name: "file without extension",
			fields: fields{
				path:              "my-binary",
				createPermissions: 0,
				fs:                nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := File{
				path:              tt.fields.path,
				createPermissions: tt.fields.createPermissions,
				fs:                tt.fields.fs,
			}
			assert.Equalf(t, tt.want, x.Extension(), "Extension()")
		})
	}
}

func TestFile_HasAnyExtension(t *testing.T) {
	type fields struct {
		path              string
		createPermissions os.FileMode
		fs                afero.Fs
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "basic file",
			fields: fields{
				path:              "/var/test/video.mp4",
				createPermissions: 0,
				fs:                nil,
			},
			want: true,
		},
		{
			name: "binary file",
			fields: fields{
				path:              "my-binary",
				createPermissions: 0,
				fs:                nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := File{
				path:              tt.fields.path,
				createPermissions: tt.fields.createPermissions,
				fs:                tt.fields.fs,
			}
			assert.Equalf(t, tt.want, x.HasAnyExtension(), "HasAnyExtension()")
		})
	}
}

func TestFile_HasExtension(t *testing.T) {
	type fields struct {
		path              string
		createPermissions os.FileMode
		fs                afero.Fs
	}
	type args struct {
		fileExtension FileExtension
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "matching extension",
			fields: fields{
				path:              "invoice.pdf",
				createPermissions: 0,
				fs:                nil,
			},
			args: args{
				fileExtension: ExtPdf,
			},
			want: true,
		},
		{
			name: "not matching extension",
			fields: fields{
				path:              "invoice.pdf",
				createPermissions: 0,
				fs:                nil,
			},
			args: args{
				fileExtension: ExtJpg,
			},
			want: false,
		},
		{
			name: "no extension, no match",
			fields: fields{
				path:              "binary-file",
				createPermissions: 0,
				fs:                nil,
			},
			args: args{
				fileExtension: ExtJpg,
			},
			want: false,
		},
		{
			name: "double extension match",
			fields: fields{
				path:              "archive.tar.gz",
				createPermissions: 0,
				fs:                nil,
			},
			args: args{
				fileExtension: FileExtensionFrom("tar.gz"),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := File{
				path:              tt.fields.path,
				createPermissions: tt.fields.createPermissions,
				fs:                tt.fields.fs,
			}
			assert.Equalf(t, tt.want, x.HasExtension(tt.args.fileExtension), "HasExtension(%v)", tt.args.fileExtension)
		})
	}
}

func TestFile_WithExtension(t *testing.T) {
	type fields struct {
		path              string
		createPermissions os.FileMode
		fs                afero.Fs
	}
	type args struct {
		fileExtension FileExtension
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "switch extension",
			fields: fields{
				path:              "whatever.jpg",
				createPermissions: 0,
				fs:                nil,
			},
			args: args{
				fileExtension: ExtPng,
			},
			want: "whatever.png",
		},
		{
			name: "set same extension again",
			fields: fields{
				path:              "whatever.jpg",
				createPermissions: 0,
				fs:                nil,
			},
			args: args{
				fileExtension: ExtJpg,
			},
			want: "whatever.jpg",
		},
		{
			name: "replace .tar.gz",
			fields: fields{
				path:              "whatever.tar.gz",
				createPermissions: 0,
				fs:                nil,
			},
			args: args{
				fileExtension: ExtZip,
			},
			want: "whatever.zip",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := File{
				path:              tt.fields.path,
				createPermissions: tt.fields.createPermissions,
				fs:                tt.fields.fs,
			}
			assert.Equalf(t, tt.want, x.WithExtension(tt.args.fileExtension).Filename(), "WithExtension(%v)", tt.args.fileExtension)
		})
	}
}
