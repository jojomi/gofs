package gofs

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileExtensionFrom(t *testing.T) {
	type args struct {
		extension string
	}
	tests := []struct {
		name string
		args args
		want FileExtension
	}{
		{
			name: "jpg",
			args: args{
				extension: "jpg",
			},
			want: FileExtension{
				extension: "jpg",
			},
		},
		{
			name: ".jpg",
			args: args{
				extension: ".jpg",
			},
			want: FileExtension{
				extension: "jpg",
			},
		},
		{
			name: ".tar.gz",
			args: args{
				extension: ".tar.gz",
			},
			want: FileExtension{
				extension: "tar.gz",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, FileExtensionFrom(tt.args.extension), "FileExtensionFrom(%v)", tt.args.extension)
		})
	}
}

func TestFileExtension_String(t *testing.T) {
	type fields struct {
		extension string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "basic",
			fields: fields{
				extension: "png",
			},
			want: ".png",
		},
		{
			name: "tar.gz",
			fields: fields{
				extension: "tar.gz",
			},
			want: ".tar.gz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := FileExtension{
				extension: tt.fields.extension,
			}
			assert.Equalf(t, tt.want, x.String(), "String()")
		})
	}
}

func TestFileExtension_WithDot(t *testing.T) {
	type fields struct {
		extension string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "txt",
			fields: fields{
				extension: "txt",
			},
			want: ".txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := FileExtension{
				extension: tt.fields.extension,
			}
			assert.Equalf(t, tt.want, x.WithDot(), "WithDot()")
		})
	}
}

func TestFileExtension_WithoutDot(t *testing.T) {
	type fields struct {
		extension string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "tex",
			fields: fields{
				extension: "tex",
			},
			want: "tex",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := FileExtension{
				extension: tt.fields.extension,
			}
			assert.Equalf(t, tt.want, x.WithoutDot(), "WithoutDot()")
		})
	}
}
