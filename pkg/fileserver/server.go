package fileserver

import (
	"net/url"
)

type FileServer interface {
	Upload(File) (*url.URL, error)
	Download(string) (*File, error)
}

type File struct {
	Name string
	Type string
	Data []byte
}
