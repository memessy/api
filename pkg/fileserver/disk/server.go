package disk

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"memessy-api/pkg/fileserver"
	"net/url"
	"path"
	"strings"
)

type FileServer struct {
	BaseUrl string
	Dir     string
}

func (server *FileServer) Upload(file fileserver.File) (*url.URL, error) {
	uniqueName := server.makeUnique(file.Name)
	savedPath := server.getPath(uniqueName)
	err := ioutil.WriteFile(savedPath, file.Data, 0644)
	if err != nil {
		return nil, err
	}
	name := path.Base(savedPath)
	u := server.getUrl(name)
	return &u, nil
}

func (server *FileServer) Download(filename string) (*fileserver.File, error) {
	filePath := server.getPath(filename)
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return &fileserver.File{Name: filename, Data: bytes}, nil
}

func (server *FileServer) getPath(name string) string {
	return path.Join(server.Dir, name)
}
func (server *FileServer) makeUnique(name string) string {
	ext := path.Ext(name)
	name = strings.TrimSuffix(name, ext)
	return name + "_" + randString(6) + ext
}
func (server *FileServer) getUrl(name string) url.URL {
	rawUrl := fmt.Sprintf("%s%s", server.BaseUrl, name)
	u, err := url.Parse(rawUrl)
	if err != nil {
		return url.URL{}
	}
	return *u
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
