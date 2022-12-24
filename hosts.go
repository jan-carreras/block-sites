package main

import (
	"io/ioutil"
)

type hostFile struct {
	hostPath string
}

func NewHostFile(hostPath string) hostFile {
	return hostFile{hostPath: hostPath}
}

func (h *hostFile) Read() ([]byte, error) {
	return ioutil.ReadFile(h.hostPath)
}

func (h *hostFile) Write(data []byte) error {
	return ioutil.WriteFile(h.hostPath, data, 0644)
}

type Cmd struct {
	Pause  bool
	Resume bool
}
